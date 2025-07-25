package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/dxasu/pure/rain"
	_ "github.com/dxasu/pure/version"
)

func main() {
	if len(os.Args) == 2 && os.Args[1] == "-h" {
		cmd := os.Args[0]
		fmt.Println(cmd, `# print clipboard content`)
		fmt.Println(cmd, `[-p] xxx or echo xxx|`+cmd+` [-p] # copy xxx to clipboard, -p: print to stdout`)
		fmt.Println(cmd, `[-p] < file  # copy file to clipboard, -p: print to stdout`)
		return
	}
	stat, err := os.Stdin.Stat()
	rain.ExitIf(err)

	if (stat.Mode()&os.ModeCharDevice) != 0 && len(os.Args) == 1 {
		content, err := clipboard.ReadAll()
		rain.ExitIf(err)
		fmt.Println(content)
		return
	}

	var data string

	bPrint := len(os.Args) > 1 && os.Args[1] == "-p"
	var idx = 1
	if bPrint {
		idx = 2
	}

	if len(os.Args) > idx {
		data = strings.Join(os.Args[idx:], " ")
	} else {
		out, err := io.ReadAll(os.Stdin)
		rain.ExitIf(err)
		if len(out) == 0 {
			rain.ExitIf(errors.New("no data to copy"))
		}
		data = string(out[:len(out)-1])
	}

	err = clipboard.WriteAll(data)
	rain.ExitIf(err)
	if bPrint {
		fmt.Println(data)
	}
}
