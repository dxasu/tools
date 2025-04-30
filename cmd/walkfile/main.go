package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
)

// 计算文件的 SHA256 哈希值
func hashFile(filePath string, fileHashes map[string]*fileHash) (*fileHash, error) {
	// 如果文件哈希值已经计算过，直接返回缓存中的值
	if hash, exists := fileHashes[filePath]; exists {
		return hash, nil
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	hash := sha256.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return nil, err
	}

	fileHashStr := hash.Sum(nil)
	// 缓存文件哈希值
	fileHashes[filePath] = &fileHash{
		fcnt: 1,
		dir:  0,
		hash: fileHashStr,
	}
	return fileHashes[filePath], nil
}

func readDirNames(dirname string) ([]string, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	names, err := f.Readdirnames(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	slices.Sort(names)
	return names, nil
}

type WalkFunc func(path string, info fs.FileInfo, err error) error

func SubWalk(root string, fn WalkFunc) error {
	info, err := os.Lstat(root)
	if err != nil {
		err = fn(root, nil, err)
	} else {
		err = walk(root, info, fn)
	}
	return err
}

func walk(path string, info fs.FileInfo, walkFn WalkFunc) error {
	if !info.IsDir() {
		return walkFn(path, info, nil)
	}
	names, err := readDirNames(path)
	err1 := walkFn(path, info, err)
	if err != nil || err1 != nil {
		return err1
	}

	for _, name := range names {
		filename := filepath.Join(path, name)
		fileInfo, err := os.Lstat(filename)
		if err := walkFn(filename, fileInfo, err); err != nil {
			return err
		}
	}
	return nil
}

// 计算文件夹的哈希值
func hashDirectory(dirPath string, fileHashes map[string]*fileHash) (*fileHash, error) {
	if hash, exists := fileHashes[dirPath]; exists {
		return hash, nil
	}

	var hashes []string
	var allCnt int64 = 0

	err := SubWalk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Size() == 0 {
			return nil
		}
		if !info.IsDir() {
			fileHash, err := hashFile(path, fileHashes)
			if err != nil {
				return err
			}
			allCnt += 1
			hashes = append(hashes, fmt.Sprintf("%x", fileHash))
		} else if path != dirPath { // 排除当前文件夹本身
			fileHashObj, err := hashDirectory(path, fileHashes)
			if err != nil {
				return err
			}
			if fileHashObj != nil {
				allCnt += fileHashObj.fcnt
				hashes = append(hashes, fmt.Sprintf("%x", fileHashObj))
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	if len(hashes) == 0 {
		return nil, nil
	}
	sort.Strings(hashes)
	hashString := strings.Join(hashes, ",")
	finalHash := sha256.New()
	finalHash.Write([]byte(hashString))
	finalHashValue := finalHash.Sum(nil)
	fileHashes[dirPath] = &fileHash{
		fcnt: allCnt,
		dir:  1,
		hash: finalHashValue,
	}

	return fileHashes[dirPath], nil
}

type fhash struct {
	hash string
	path []string
}

type fileHash struct {
	fcnt int64
	dir  int64
	hash []byte
}

func calcRoot(rootDir string) ([]fhash, error) {
	fileHashes := make(map[string]*fileHash)
	fileHs, err := hashDirectory(rootDir, fileHashes)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Root hashInfo fcnt: %d hash: %x path: %s\n", fileHs.fcnt, fileHs.hash, rootDir)

	valueToKeys := make(map[string][]string)
	for key, value := range fileHashes {
		valueStr := fmt.Sprintf("%d_%06d_%x", value.dir, value.fcnt, value.hash)
		valueToKeys[valueStr] = append(valueToKeys[valueStr], key)
	}

	ret := make([]fhash, 0)
	for value, keys := range valueToKeys {
		if len(keys) > 1 {
			ret = append(ret, fhash{hash: value, path: keys})
		}
	}

	sort.Slice(ret, func(i, j int) bool {
		return ret[i].hash > ret[j].hash
	})
	return ret, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: walkfile <rootDir>")
		return
	}
	infos, err := calcRoot(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	for _, info := range infos {
		fmt.Printf("Hash: %s\n", info.hash)
		for _, path := range info.path {
			fmt.Printf("  %s\n", path)
		}
	}
}
