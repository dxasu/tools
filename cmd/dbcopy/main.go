package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/dxasu/pure/rain"
	_ "github.com/go-sql-driver/mysql"
)

const (
	// MaxCount is the maximum number of rows to copy from each table.
	MaxCount = 100 // 设置为0表示不限制
	// MaxConcurrent is the maximum number of concurrent database operations.
	MaxConcurrent = 10 // 并发数限制
	// MaxBatchSize is the maximum number of rows to insert in a single batch.
	MaxBatchSize = 1000 // 每次批量插入的最大行数
	// MaxOpenConns is the maximum number of open connections to the database.
	MaxOpenConns = 100 // 最大打开连接数
)

func main() {
	if rain.NeedHelp() {
		fmt.Println("Usage: dbcopy")
		fmt.Println("Environment Variables:")
		fmt.Println("  SOURCE_DSN: The Data Source Name for the source MySQL database")
		fmt.Println("  TARGET_DSN(default): root:123456@tcp(127.0.0.1:3306)?charset=utf8mb4&parseTime=True&loc=Local")
		fmt.Println("This tool copies all non-system databases and their tables from a source MySQL database to a target MySQL database.")
		fmt.Println("It connects to the source database, retrieves all non-system databases, and copies each database and its tables to the target database.")
		return
	}

	sourceDSN := os.Getenv("SOURCE_DSN")
	if sourceDSN == "" {
		log.Fatal("SOURCE_DSN environment variable is not set")
	}
	// Default target DSN if not set
	targetDSN := os.Getenv("TARGET_DSN")
	if targetDSN == "" {
		targetDSN = "root:123456@tcp(127.0.0.1:3306)?charset=utf8mb4&parseTime=True&loc=Local"
	}
	sourceDB, err := connectDB(sourceDSN)
	rain.ExitIf(err)
	targetDB, err := connectDB(targetDSN)
	rain.ExitIf(err)
	defer sourceDB.Close()
	defer targetDB.Close()
	sourceDB.SetMaxOpenConns(MaxOpenConns)
	targetDB.SetMaxOpenConns(MaxOpenConns)
	databases, err := getNonSystemDatabases(sourceDB)
	rain.ExitIf(err)
	var defaulCount int = MaxCount
	var wg sync.WaitGroup
	sem := make(chan struct{}, MaxConcurrent)
	for _, dbName := range databases {
		wg.Add(1)
		go func(dbName string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			err = createDatabaseIfNotExists(targetDB, dbName)
			if err != nil {
				log.Printf("Failed to create database %s: %v", dbName, err)
				return
			}
			tables, err := getTablesInDatabase(sourceDB, dbName)
			if err != nil {
				log.Printf("Failed to get tables in %s: %v", dbName, err)
				return
			}
			for _, tableName := range tables {
				err = copyTable(sourceDB, targetDB, dbName, dbName, tableName, defaulCount)
				if err != nil {
					log.Printf("Failed to copy table %s.%s: %v", dbName, tableName, err)
				}
			}
		}(dbName)
	}
	wg.Wait()
}

func connectDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func getNonSystemDatabases(db *sql.DB) ([]string, error) {
	rows, err := db.Query(`
        SELECT schema_name 
        FROM information_schema.SCHEMATA 
        WHERE schema_name NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var databases []string
	for rows.Next() {
		var dbName string
		err = rows.Scan(&dbName)
		if err != nil {
			return nil, err
		}
		databases = append(databases, dbName)
	}
	return databases, nil
}

func createDatabaseIfNotExists(db *sql.DB, dbName string) error {
	_, err := db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", dbName))
	return err
}

func getTablesInDatabase(db *sql.DB, dbName string) ([]string, error) {
	rows, err := db.Query(fmt.Sprintf("SHOW TABLES IN `%s`", dbName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		err = rows.Scan(&tableName)
		if err != nil {
			return nil, err
		}
		tables = append(tables, tableName)
	}
	return tables, nil
}

func copyTable(sourceDB, targetDB *sql.DB, sourceDBName, targetDBName, tableName string, maxCount int) error {
	// 1. 复制表结构
	var createTableSQL string
	err := sourceDB.QueryRow(fmt.Sprintf("SHOW CREATE TABLE `%s`.`%s`", sourceDBName, tableName)).Scan(&tableName, &createTableSQL)
	if err != nil {
		return err
	}
	_, err = targetDB.Exec(fmt.Sprintf("USE `%s`", targetDBName))
	if err != nil {
		return err
	}
	_, err = targetDB.Exec(createTableSQL)
	if err != nil {
		return err
	}

	// 2. 复制数据（分批处理）
	offset := 0
	batchSize := MaxCount
	for {
		if maxCount > 0 {
			if offset >= maxCount {
				break
			}
			if offset+batchSize > maxCount {
				batchSize = maxCount - offset
			}
		}
		query := fmt.Sprintf("SELECT * FROM `%s`.`%s` LIMIT %d OFFSET %d", sourceDBName, tableName, batchSize, offset)
		rows, err := sourceDB.Query(query)
		if err != nil {
			return err
		}
		defer rows.Close()

		if !rows.Next() {
			break
		}

		columns, err := rows.Columns()
		if err != nil {
			return err
		}

		// 准备批量插入语句
		placeholders := strings.Repeat("?,", len(columns))
		placeholders = placeholders[:len(placeholders)-1]
		insertSQL := fmt.Sprintf("INSERT INTO `%s`.`%s` (%s) VALUES (%s)",
			targetDBName, tableName, strings.Join(columns, ","), placeholders)

		stmt, err := targetDB.Prepare(insertSQL)
		if err != nil {
			return err
		}
		defer stmt.Close()

		// 处理当前批次数据
		batchCount := 0
		for {
			values := make([]interface{}, len(columns))
			scanArgs := make([]interface{}, len(values))
			for i := range values {
				scanArgs[i] = &values[i]
			}
			err = rows.Scan(scanArgs...)
			if err != nil {
				return err
			}

			_, err = stmt.Exec(values...)
			if err != nil {
				return err
			}
			batchCount++

			if batchCount >= batchSize || !rows.Next() {
				break
			}
		}
		offset += batchCount
	}
	return nil
}
