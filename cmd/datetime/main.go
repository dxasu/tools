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
	datetime -[zfcCtTuUnh]
Flags:
	-z use utc+0. not location
	-f formatStr need, "2006-01-02 15:04:05"
	-a parse by auto fit
	-u parse from unixtime(revert -U)
	-c copy from clipboard(revert -C)
	-t duration from "-2h45m6s6ms6us8ns"(revert -T)
	-n print nothing but error
	-h help
`)
		return
	}

	var (
		cmd        string
		data       string
		format     string = "2006-01-02 15:04:05"
		err        error
		timeRegion *time.Location = time.Local
	)

	if len(os.Args) == 1 {
		println(time.Now().Unix())
		println(time.Now().Format(format))
		return
	}

	param := os.Args[1]
	if param[0] != '-' {
		cmd = "-a"
		data = strings.Join(os.Args[1:], " ")
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

	t := &timeFly{Data: []byte(data), region: timeRegion}
	show := true

	for _, v := range cmd {
		switch v {
		case '-', 'f':
		case 'a':
			t.AutoParse(format)
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
		case 't':
			t.ParseDuration()
		case 'T':
			t.ParseToUnit()
		case 'z':
			t.region = time.UTC
		case 'n':
			show = false
		default:
			rain.ExitIf(fmt.Errorf("invalid param: %c", v))
		}
	}

	// RESULT:
	if show {
		println(string(t.Data))
	}
}

type timeFly struct {
	t      time.Time
	region *time.Location
	Data   []byte
}

func (t *timeFly) ToClipboard() {
	clipboard.WriteAll(string(t.Data))
}

func (t *timeFly) ParseFromUnixTime(format string) {
	tSec, err := cast.ToInt64E(string(t.Data))
	rain.ExitIf(err)
	t.t = time.Unix(tSec, 0).In(t.region)
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
	var err error
	t.t, err = time.ParseInLocation(format, string(t.Data), t.region)
	if err != nil {
		t.t, err = cast.ToTimeE(t.Data)
		if err != nil {
			t.t, err = cast.StringToDate(string(t.Data))
			rain.ExitIf(err)
		}
	}
	t.Data = []byte(t.t.Format(format))
}

func (t *timeFly) AutoParse(format string) {
	_, err := cast.ToInt64E(string(t.Data))
	if err == nil {
		if len(t.Data) == len(cast.ToString(time.Now().Unix())) {
			t.ParseFromUnixTime(format)
		} else {
			t.ParseToUnit()
		}
	} else {
		var err error
		t.t, err = time.ParseInLocation(format, string(t.Data), t.region)
		if err != nil {
			t.t, err = cast.ToTimeE(t.Data)
			if err != nil {
				t.t, err = cast.StringToDate(string(t.Data))
			}
		}
		if err == nil {
			t.Data = []byte(cast.ToString(t.t.Unix()))
			return
		}
		t.ParseDuration()
	}
}

// such as "300ms", "-1.5h" or "2h45m".
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
func (t *timeFly) ParseDuration() {
	dur, err := time.ParseDuration(string(t.Data))
	rain.ExitIf(err)
	t.Data = []byte(cast.ToString(dur.Seconds()))
}

func (t *timeFly) ParseToUnit() {
	t.Data = []byte(cast.ToString(cast.ToDuration(string(t.Data)) * time.Second))
}
