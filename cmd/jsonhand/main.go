package main

import (
	"bytes"
	"encoding/json"
	"os"
	"strconv"
	"strings"

	"bay.core/lancet/rain"
	"github.com/atotto/clipboard"
	_ "github.com/dxasu/tools/lancet/version"
)

func main() {
	if len(os.Args) < 2 {
		println(`args empty, example: jsonhand -[uzeg] {"today":"go"}
-u switch json between string and json string
-z switch json between single line and multi
-g switch json to go struct
`)
		return
	}
	var (
		data string
		err  error
	)
	if len(os.Args) < 3 {
		data, err = clipboard.ReadAll()
		rain.ExitIf(err)
	} else {
		data = os.Args[2]
	}

	if strings.Contains(os.Args[1], "u") {
		if data[0] == '"' {
			data, err = strconv.Unquote(data)
			rain.ExitIf(err)
		} else {
			data = strconv.Quote(data)
		}
	}

	if strings.Contains(os.Args[1], "z") {
		if strings.Contains(data, "\n") {
			dst := &bytes.Buffer{}
			err = json.Indent(dst, []byte(data), "", "")
			rain.ExitIf(err)
			data = dst.String()
			data = strings.ReplaceAll(data, "\n", "")
		} else {
			dst := &bytes.Buffer{}
			err = json.Indent(dst, []byte(data), "", "    ")
			rain.ExitIf(err)
			data = dst.String()
		}
	}

	if strings.Contains(os.Args[1], "g") {
		dst := &bytes.Buffer{}
		err = json.Indent(dst, []byte(data), "", "")
		rain.ExitIf(err)
		data = dst.String()
	}

	print(data)
}
