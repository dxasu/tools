package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"reflect"
	"strings"
	"time"

	"github.com/dxasu/coreemail"
	"github.com/dxasu/pure/rain"
	"github.com/dxasu/pure/stdin"
	"github.com/dxasu/pure/storage"
	"github.com/h2non/filetype"
)

func main() {
	d := &storage.Data{}
	err := d.Init("hub_temp_mail.json")
	if err != nil {
		fmt.Println(err)
		return
	}

HELP:
	if len(os.Args) == 2 && (os.Args[1] == "help" || os.Args[1] == "-h") {
		fmt.Println("Usage: <command> [args]")
		fmt.Println("Commands:")
		fmt.Println("  get - Get the value of all keys")
		fmt.Println("  set <key> <value> - Set the value for the specified key")
		fmt.Println("Infos:")
		fmt.Println("QQ: https://wx.mail.qq.com/list/readtemplate?name=app_intro.html#/agreement/authorizationCode")
		return
	}
	if len(os.Args) >= 2 && os.Args[1] == "set" {
		if len(os.Args) < 4 || len(os.Args)%2 != 0 {
			fmt.Println("Usage: set <key> <value>")
			return
		}
		for i := 2; i < len(os.Args); i += 2 {
			key := os.Args[i]
			value := os.Args[i+1]
			err = d.Set(key, value)
			if err != nil {
				fmt.Println(key, value, err)
			}
		}

		err = d.Save()
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	m := &coreemail.Mail{
		SenderAddr:   "sender@163.com",
		SenderName:   "senderName",
		ReceiverAddr: []string{"receiver@163.com"},
		Subject:      "subject",
		Text:         "test",
		FilePaths:    []string{}, // data.txt
		Host:         "smtp.163.com",
		Port:         25,
		Username:     "username@163.com",
		Password:     "",
	}

	for _, k := range reflect.VisibleFields(reflect.TypeOf(*m)) {
		v := reflect.ValueOf(m).Elem().FieldByName(k.Name)
		switch v.Type().Kind() {
		case reflect.String:
			value := d.GetString(k.Name)
			if value != "" {
				v.SetString(value)
			}
		case reflect.Slice:
			intValue := d.GetString(k.Name)
			if intValue != "" {
				v.Set(reflect.ValueOf([]string{intValue}))
				// fmt.Printf("Key:%s Value:%v \n", k.Name, v.Interface())
			}
		case reflect.Int, reflect.Int64:
			intValue := d.GetInt(k.Name)
			if intValue != 0 {
				v.SetInt(int64(intValue))
			}
		}
	}

	if len(os.Args) >= 2 && os.Args[1] == "get" {
		fmt.Println("\033[32mMail file path:\033[0m")
		fmt.Println()
		fmt.Println(d.GetFilePath())
		d.PrintJSON()
		mail, _ := json.MarshalIndent(m, "", "  ")
		fmt.Println("\n\033[32mMail settings:\033[0m")
		fmt.Println(string(mail))
		return
	}

	if m.Password == "" {
		os.Args = append([]string{os.Args[0]}, "help")
		goto HELP
	}

	data := stdin.GetInput(true)
	if len(data) > 0 {
		fileExt := "txt"
		if len(data) < 2<<12 { // 4K
			stat, err := os.Stat(string(data))
			if err != nil || stat.IsDir() {
				m.Text = strings.TrimSuffix(data, "\n")
				fileExt = ""
			} else {
				datas, err := os.ReadFile(data)
				rain.ExitIf(err)
				fileExt = path.Ext(data)
				fileExt = fileExt[1:]
				m.Text = data
				m.Subject = path.Base(data)
				data = string(datas)
			}
		}
		if fileExt != "" {
			m.FilePaths = append(m.FilePaths, GetFile(fileExt, []byte(data)))
		}
	} else {
		path, err := saveImageFromClipboard()
		rain.ExitIf(err)
		m.FilePaths = append(m.FilePaths, path)
	}
	err = m.Validate()
	rain.ExitIf(err)

	err = m.Send()
	rain.ExitIf(err)
}

// saveImageFromClipboard 从剪贴板保存图片
func saveImageFromClipboard() (string, error) {
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("clipboard_image_%s.png", timestamp)
	filePath := storage.NewDataPath(storage.DirTemp, "gmail_temp").GetFilePath(filename)

	// 使用AppleScript从剪贴板获取图片并保存
	script := fmt.Sprintf(`
		try
			set imageData to (the clipboard as «class PNGf»)
			set fileRef to open for access POSIX file "%s" with write permission
			write imageData to fileRef
			close access fileRef
			return "success"
		on error errMsg
			try
				close access fileRef
			end try
			return "error: " & errMsg
		end try
	`, filePath)

	cmd := exec.Command("osascript", "-e", script)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("保存图片失败: %v", err)
	}

	result := strings.TrimSpace(string(output))
	if strings.HasPrefix(result, "error:") {
		return "", fmt.Errorf("AppleScript错误: %s", result)
	}

	return filePath, nil
}

func GetFile(fileExt string, data []byte) string {
	t, err := filetype.Get(data)
	rain.ExitIf(err)
	if t.Extension != "unknown" {
		fileExt = t.Extension
	}
	path := storage.NewDataPath(storage.DirTemp, "gmail_temp").GetFilePath(time.Now().Format("060102_150405") + "." + fileExt)
	err = os.WriteFile(path, data, os.ModePerm)
	rain.ExitIf(err)
	return path
}
