package main

import (
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
	datetime -[fcCnduU]
Flags:
	-f formatStr need, default use -c
	-c copy from clipboard
	-C copy to clipboard
	-u parse unixtime
	-U parse to unixtime
	-n print nothing but error
	-d now time
	-D duration like "-2h45m". unit as "ns", "us" (or "µs"), "ms", "s", "m", "h".
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
		data = param
	} else {
		cmd = param
		if len(os.Args) > 2 {
			if strings.ContainsRune(cmd, 'f') {
				format = os.Args[2]
			} else {
				data = os.Args[2]
			}
		}

	}

	t := &timeFly{Data: []byte(data)}
	show := true

	if cmd == "" {
		t.ParseTime(format)
	} else {
		for _, v := range cmd {
			switch v {
			case '-':
			case 'u':
				t.FillTime(format)
				t.ParseUnixTime(format)
			case 'U':
				t.FillTime(format)
				t.ParseToUnixTime(format)
				goto RESULT
			case 'c', 'f':
				data, err = clipboard.ReadAll()
				rain.ExitIf(err)
				t.Data = []byte(data)
				t.ParseTime(format)
			case 'C':
				t.ToClipboard()
			case 'd':
				t.Data = []byte(time.Now().Format(format))
				t.ParseTime(format)
			case 'D':
				t.Data = []byte(t.ParseDuration().String())
				goto RESULT
			case 'n':
				show = false
			default:
				rain.ExitIf(fmt.Errorf("invalid param: %c", v))
			}
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

func (t *timeFly) ParseUnixTime(format string) {
	tSec, err := cast.ToInt64E(string(t.Data))
	rain.ExitIf(err)
	t.t = time.Unix(tSec, 0)
	t.Data = []byte(t.t.Format(format))
}

func (t *timeFly) ParseToUnixTime(format string) {
	t.Data = []byte(cast.ToString(t.t.Unix()))
}

func (t *timeFly) FillTime(format string) {
	if len(t.Data) == 0 {
		t.t = time.Now()
		t.Data = []byte(t.t.Format(format))
	} else if t.t.IsZero() {
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
