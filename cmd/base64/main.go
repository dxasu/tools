package main

import (
	"encoding/base64"
	"errors"
	"net/url"
	"os"
	"strings"

	"bay.core/lancet/rain"
	"github.com/atotto/clipboard"
	_ "github.com/dxasu/tools/lancet/version"
)

func main() {
	if rain.NeedHelp() {
		println(`Usage:
	base64 -[udecn] [xxx]
Flags:
	-u base64url mode
	-d base64 decode
	-e base64 encode
	-D url QueryUnescape decode
	-E url QueryEscape encode
	-c copy to clipboard
	-n print nothing but error
`)
		return
	}

	var (
		data string
		err  error
	)
	if len(os.Args) == 1 || len(os.Args) == 2 && os.Args[1][0] == '-' {
		data, err = clipboard.ReadAll()
		rain.ExitIf(err)
	} else {
		if len(os.Args) == 2 {
			data = os.Args[1]
		} else {
			data = os.Args[2]
		}
	}

	var cmd string
	if len(os.Args) > 1 {
		cmd = os.Args[1]
		if cmd[0] != '-' {
			cmd = "-d"
		}
	} else {
		cmd = "-d"
	}

	show := true
	urlMode := false
	copy := false
	for _, v := range cmd {
		switch v {
		case '-':
		case 'u':
			urlMode = true
		case 'U':
			data = url.QueryEscape(data)
		case 'E':
			data, err = url.QueryUnescape(data)
			rain.ExitIf(err)
		case 'd':
			bytes, err := base64.StdEncoding.DecodeString(data)
			rain.ExitIf(err)
			data = string(bytes)
		case 'e':
			bytes := base64.StdEncoding.EncodeToString([]byte(data))
			data = string(bytes)
		case 'c':
			copy = true
		case 'n':
			show = false
		}
	}

	if len(data) == 0 {
		rain.ExitIf(errors.New("empty data"))
	}

	if urlMode {
		data = strings.NewReplacer("+", "-", "/", "_", "=", "").Replace(data)
	}

	if show {
		print(data)
	}
	if copy {
		clipboard.WriteAll(data)
	}
}
