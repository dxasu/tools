package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"bay.core/lancet/rain"
	"github.com/atotto/clipboard"
	_ "github.com/dxasu/tools/lancet/version"
	"github.com/spf13/cast"
)

// var r = strings.NewReplacer("\r", "", "\n", "")

func main() {
	if rain.NeedHelp() {
		println(`Usage:
	datetime -[cCnduUf]
Flags:
	-f formatStr need, 2006-01-02 15:04:05
	-u parse from unixtime
	-U parse to unixtime
	-d now time
	-c copy from clipboard
	-C copy to clipboard
	-D duration like "-2h45m". unit as "ns", "us", "ms", "s", "m", "h"
	-n print nothing but error
`)
		return
	}

	var (
		cmd    string
		data   string
		format string = "2006-01-02 15:04:05"
		err    error
	)

	if len(os.Args) == 1 {
		println(time.Now().Unix())
		println(time.Now().Format(format))
		return
	}

	param := os.Args[1]
	if param[0] != '-' {
		cmd = "-d"
		data = param
	} else {
		cmd = param
		if cmd == "-f" {
			cmd = "-fd"
		}
		if len(os.Args) > 2 {
			if strings.ContainsRune(cmd, 'f') {
				format = os.Args[2]
				if len(os.Args) > 3 {
					data = os.Args[3]
				} else {
					data = time.Now().Format(format)
				}
			} else {
				data = os.Args[2]
			}
		} else {
			data = time.Now().Format(format)
		}
	}

	if len(data) == 0 {
		rain.ExitIf(errors.New("data is empty"))
	}

	t := &timeFly{Data: []byte(data)}
	show := true

	for _, v := range cmd {
		switch v {
		case '-', 'f':
		case 'u':
			t.ParseFromUnixTime(format)
		case 'U':
			t.ParseToUnixTime(format) // unixtime in data
		case 'c':
			data, err = clipboard.ReadAll()
			rain.ExitIf(err)
			t.Data = []byte(data)
		case 'C':
			t.ToClipboard()
		case 'd':
			t.ParseTime(format)
		case 'D':
			t.Data = []byte(cast.ToString(t.ParseDuration().Seconds()))
			goto RESULT
		case 'n':
			show = false
		default:
			rain.ExitIf(fmt.Errorf("invalid param: %c", v))
		}
	}

RESULT:
	if show {
		println(string(t.Data))
	}
}

type timeFly struct {
	t    time.Time
	Data []byte
}

func (t *timeFly) ToClipboard() {
	clipboard.WriteAll(string(t.Data))
}

func (t *timeFly) ParseFromUnixTime(format string) {
	tSec, err := cast.ToInt64E(string(t.Data))
	rain.ExitIf(err)
	t.t = time.Unix(tSec, 0)
	t.Data = []byte(t.t.Format(format))
}

func (t *timeFly) ParseToUnixTime(format string) {
	t.FillTime(format)
	t.Data = []byte(cast.ToString(t.t.Unix()))
}

func (t *timeFly) FillTime(format string) {
	if t.t.IsZero() {
		t.ParseTime(format)
	}
}

func (t *timeFly) ParseTime(format string) {
	tTime, err := cast.ToTimeE(t.Data)
	if err == nil {
		t.t = tTime
	} else {
		tTime, err = cast.StringToDate(string(t.Data))
		rain.ExitIf(err)
		t.t = tTime
	}
	t.Data = []byte(t.t.Format(format))
}

// such as "300ms", "-1.5h" or "2h45m".
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
func (t *timeFly) ParseDuration() time.Duration {
	dur, err := time.ParseDuration(string(t.Data))
	rain.ExitIf(err)
	return dur
}
