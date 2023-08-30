package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"bay.core/lancet/rain"
	"github.com/atotto/clipboard"
	"github.com/dxasu/tools/cmd/jsonhand/j2struct"
	_ "github.com/dxasu/tools/lancet/version"
	"github.com/spf13/viper"
)

// var r = strings.NewReplacer("\r", "", "\n", "")

func main() {
	if rain.NeedHelp() {
		println(`args empty, example: jsonhand -[ju | fz | qgGcn] [xxx | yaml]
-j to json by yaml, yml, toml, ini, env. Data source must from clipboard
-u unquote string to json
-f format json
-z zip json
-q quote json to string
-g spawn go struct and sub (-G with single struct)
-c copy to clipboard
-n print nothing but error
-d as -uf default without params. 
`)
		return
	}

	var cmd string
	if len(os.Args) > 1 {
		cmd = os.Args[1]
		if cmd[0] != '-' {
			rain.ExitIf(errors.New("invalid param"))
		}
	} else {
		cmd = "-d"
	}
	
	var (
		data string
		err  error
	)
	if len(os.Args) < 3 || strings.ContainsRune(cmd, 'j') {
		data, err = clipboard.ReadAll()
		rain.ExitIf(err)
	} else {
		data = os.Args[2]
	}

	if len(data) == 0 {
		rain.ExitIf(errors.New("empty data"))
	}

	j := &jsonFly{[]byte(data)}
	show := true
	for _, v := range cmd {
		switch v {
		case '-':
		case 'f':
			j.Format()
		case 'z':
			j.Zip()
		case 'q':
			j.Quote()
		case 'u':
			j.UnQuote()
		case 'g':
			j.ToStruct(true)
		case 'G':
			j.ToStruct(false)
		case 'c':
			j.ToClipboard()
		case 'd':
			j.Default()
		case 'j':
			if len(os.Args) < 3 {
				rain.ExitIf(errors.New("-j need param like yaml, yml, toml, ini, env"))
			}
			j.ParseToJson(os.Args[2])
		case 'n':
			show = false
		default:
			rain.ExitIf(fmt.Errorf("invalid param: %c", v))
		}
	}
	if show {
		println(string(j.Data))
	}
}

type jsonFly struct {
	Data []byte
}

func (j *jsonFly) Zip() {
	dst := &bytes.Buffer{}
	err := json.Compact(dst, j.Data)
	rain.ExitIf(err)
	j.Data = dst.Bytes()
}

func (j *jsonFly) UnQuote() {
	data, err := strconv.Unquote(string(j.Data))
	rain.ExitIf(err)
	j.Data = []byte(data)
}

func (j *jsonFly) Quote() {
	j.Data = []byte(strconv.Quote(string(j.Data)))
}

func (j *jsonFly) Format() {
	dst := &bytes.Buffer{}
	err := json.Indent(dst, j.Data, "", strings.Repeat(" ", 4))
	rain.ExitIf(err)
	j.Data = dst.Bytes()
}

func (j *jsonFly) ToStruct(subStruct bool) {
	convertFloats := true
	data, err := j2struct.Generate(bytes.NewBuffer(j.Data), "Core", []string{"json"}, subStruct, convertFloats)
	rain.ExitIf(err)
	j.Data = data
}

func (j *jsonFly) ToClipboard() {
	clipboard.WriteAll(string(j.Data))
}

// "yaml", "yml", "json", "toml", "hcl", "tfvars", "ini", "properties", "props", "prop", "dotenv", "env"
func (j *jsonFly) ParseToJson(t string) {
	viper.SetConfigType(t)
	viper.ReadConfig(bytes.NewBuffer(j.Data))
	var err error
	j.Data, err = json.Marshal(viper.AllSettings())
	rain.ExitIf(err)
}

func (j *jsonFly) Default() {
	data, err1 := strconv.Unquote(string(j.Data))
	if err1 == nil {
		j.Data = []byte(data)
	}
	dst := &bytes.Buffer{}
	err2 := json.Indent(dst, j.Data, "", strings.Repeat(" ", 4))
	if err2 == nil {
		j.Data = dst.Bytes()
	}
	if err1 != nil && err2 != nil {
		j.Data = []byte(fmt.Errorf("parse with err1:%w \nerr2:%w", err1, err2).Error())
	}
}
