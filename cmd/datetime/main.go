package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/dxasu/pure/rain"
	_ "github.com/dxasu/pure/version"
	"github.com/spf13/cast"
)

// var r = strings.NewReplacer("\r", "", "\n", "")
// PS D:\arena\tools\cmd\datetime> copy-item .\datetime.exe D:\gowrok\bin\datetime.exe
const timeFormat = "2006-01-02 15:04:05"

func main() {
	if rain.NeedHelp() {
		fmt.Println(`Usage:
	datetime -[zfaAcCtTuUnh]
Flags:
	-z use utc+0. not location
	-f formatStr need, "2006-01-02 15:04:05"
	-a parse by auto fit
	-A time calculate
	-u parse from unixtime(revert -U)
	-c copy from clipboard(revert -C)
	-t duration from "-9y8d2h45m6s6ms6us8ns"(revert -T)
	-n print nothing but error
	-h help
`)
		return
	}
	var (
		cmd        string
		data       string
		format     string = timeFormat
		err        error
		timeRegion *time.Location = time.Local
	)

	if len(os.Args) == 1 {
		fmt.Println(time.Now().Unix())
		fmt.Println(time.Now().Format(format))
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
		} else if strings.ContainsRune(cmd, 'z') {
			timeRegion = time.UTC
		}

		if len(os.Args) <= 2 {
			if strings.ContainsRune(cmd, 'A') {
				rain.ExitIf(errors.New("miss param when use -A"))
			}
			data = time.Now().In(timeRegion).Format(format)
		} else if strings.ContainsAny(cmd, "fA") {
			format = os.Args[2]
			if len(os.Args) > 3 {
				data = os.Args[3]
			} else if strings.ContainsRune(cmd, 'A') {
				data = time.Now().In(timeRegion).Format(timeFormat)
			} else {
				data = time.Now().In(timeRegion).Format(format)
			}
		} else {
			data = os.Args[2]
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
		case 'A':
			t.CalculateTime(format)
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
		fmt.Println(string(t.Data))
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
	t.Data = []byte(t.t.In(t.region).Format(format))
}

func (t *timeFly) CalculateTime(tData string) {
	t.FillTime(timeFormat)
	sub := tData[0] == '-'
	if sub {
		tData = tData[1:]
	}
	secNum, err := cast.ToInt64E(string(tData))
	var dur time.Duration
	if err != nil {
		var (
			t2  time.Time
			err error
		)
		t2, err = time.ParseInLocation(timeFormat, tData, t.region)
		if err != nil {
			t2, err = cast.ToTimeE(tData)
			if err != nil {
				t2, err = cast.StringToDate(string(tData))
			}
		}
		if err == nil {
			t.Data = []byte(DurationToString(t.t.Sub(t2)))
			return
		}
		dur, err = ParseDuration(string(tData))
		rain.ExitIf(err)
	} else {
		dur = time.Duration(secNum)
	}
	if sub {
		t.t = t.t.Add(-dur)
	} else {
		t.t = t.t.Add(dur)
	}
	t.Data = []byte(t.t.In(t.region).Format(timeFormat))
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

func (t *timeFly) ParseDuration() {
	dur, err := ParseDuration(string(t.Data))
	rain.ExitIf(err)
	t.Data = []byte(cast.ToString(dur.Seconds()))
}

func (t *timeFly) ParseToUnit() {
	t.Data = []byte(DurationToString(cast.ToDuration(string(t.Data)) * time.Second))
}

var errLeadingInt = errors.New("time: bad [0-9]*") // never printed

// leadingInt consumes the leading [0-9]* from s.
func leadingInt(s string) (x uint64, rem string, err error) {
	i := 0
	for ; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			break
		}
		if x > 1<<63/10 {
			// overflow
			return 0, "", errLeadingInt
		}
		x = x*10 + uint64(c) - '0'
		if x > 1<<63 {
			// overflow
			return 0, "", errLeadingInt
		}
	}
	return x, s[i:], nil
}

// leadingFraction consumes the leading [0-9]* from s.
// It is used only for fractions, so does not return an error on overflow,
// it just stops accumulating precision.
func leadingFraction(s string) (x uint64, scale float64, rem string) {
	i := 0
	scale = 1
	overflow := false
	for ; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			break
		}
		if overflow {
			continue
		}
		if x > (1<<63-1)/10 {
			// It's possible for overflow to give a positive number, so take care.
			overflow = true
			continue
		}
		y := x*10 + uint64(c) - '0'
		if y > 1<<63 {
			overflow = true
			continue
		}
		x = y
		scale *= 10
	}
	return x, scale, s[i:]
}

var unitMap = map[string]uint64{
	"ns": uint64(time.Nanosecond),
	"us": uint64(time.Microsecond),
	"µs": uint64(time.Microsecond), // U+00B5 = micro symbol
	"μs": uint64(time.Microsecond), // U+03BC = Greek letter mu
	"ms": uint64(time.Millisecond),
	"s":  uint64(time.Second),
	"m":  uint64(time.Minute),
	"h":  uint64(time.Hour),
	"d":  uint64(time.Hour * 24),
	"y":  uint64(time.Hour * 24 * 365),
}

// ParseDuration parses a duration string.
// A duration string is a possibly signed sequence of
// decimal numbers, each with optional fraction and a unit suffix,
// such as "300ms", "-1.5h" or "2h45m".
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
func ParseDuration(s string) (time.Duration, error) {
	// [-+]?([0-9]*(\.[0-9]*)?[a-z]+)+
	orig := s
	var d uint64
	neg := false

	// Consume [-+]?
	if s != "" {
		c := s[0]
		if c == '-' || c == '+' {
			neg = c == '-'
			s = s[1:]
		}
	}
	// Special case: if all that is left is "0", this is zero.
	if s == "0" {
		return 0, nil
	}
	if s == "" {
		return 0, errors.New("time: invalid duration " + orig)
	}
	for s != "" {
		var (
			v, f  uint64      // integers before, after decimal point
			scale float64 = 1 // value = v + f/scale
		)

		var err error

		// The next character must be [0-9.]
		if !(s[0] == '.' || '0' <= s[0] && s[0] <= '9') {
			return 0, errors.New("time: invalid duration " + orig)
		}
		// Consume [0-9]*
		pl := len(s)
		v, s, err = leadingInt(s)
		if err != nil {
			return 0, errors.New("time: invalid duration " + orig)
		}
		pre := pl != len(s) // whether we consumed anything before a period

		// Consume (\.[0-9]*)?
		post := false
		if s != "" && s[0] == '.' {
			s = s[1:]
			pl := len(s)
			f, scale, s = leadingFraction(s)
			post = pl != len(s)
		}
		if !pre && !post {
			// no digits (e.g. ".s" or "-.s")
			return 0, errors.New("time: invalid duration " + (orig))
		}

		// Consume unit.
		i := 0
		for ; i < len(s); i++ {
			c := s[i]
			if c == '.' || '0' <= c && c <= '9' {
				break
			}
		}
		if i == 0 {
			return 0, errors.New("time: missing unit in duration " + (orig))
		}
		u := s[:i]
		s = s[i:]
		unit, ok := unitMap[u]
		if !ok {
			return 0, errors.New("time: unknown unit " + (u) + " in duration " + (orig))
		}
		if v > 1<<63/unit {
			// overflow
			return 0, errors.New("time: invalid duration " + (orig))
		}
		v *= unit
		if f > 0 {
			// float64 is needed to be nanosecond accurate for fractions of hours.
			// v >= 0 && (f*unit/scale) <= 3.6e+12 (ns/h, h is the largest unit)
			v += uint64(float64(f) * (float64(unit) / scale))
			if v > 1<<63 {
				// overflow
				return 0, errors.New("time: invalid duration " + (orig))
			}
		}
		d += v
		if d > 1<<63 {
			return 0, errors.New("time: invalid duration " + (orig))
		}
	}
	if neg {
		return -time.Duration(d), nil
	}
	if d > 1<<63-1 {
		return 0, errors.New("time: invalid duration " + (orig))
	}
	return time.Duration(d), nil
}

func DurationToString(d time.Duration) string {
	// Largest time is 2540400h10m10.000000000s
	var buf [32]byte
	w := len(buf)

	u := uint64(d)
	neg := d < 0
	if neg {
		u = -u
	}

	if u < uint64(time.Second) {
		// Special case: if duration is smaller than a second,
		// use smaller units, like 1.2ms
		var prec int
		w--
		buf[w] = 's'
		w--
		switch {
		case u == 0:
			return "0s"
		case u < uint64(time.Microsecond):
			// print nanoseconds
			prec = 0
			buf[w] = 'n'
		case u < uint64(time.Millisecond):
			// print microseconds
			prec = 3
			// U+00B5 'µ' micro sign == 0xC2 0xB5
			w-- // Need room for two bytes.
			copy(buf[w:], "µ")
		default:
			// print milliseconds
			prec = 6
			buf[w] = 'm'
		}
		w, u = fmtFrac(buf[:w], u, prec)
		w = fmtInt(buf[:w], u)
	} else {
		if u%60 > 0 {
			w--
			buf[w] = 's'
			w, u = fmtFrac(buf[:w], u, 9)
			// u is now integer seconds
			w = fmtInt(buf[:w], u%60)
		} else {
			u /= uint64(time.Second)
		}
		u /= 60

		// u is now integer minutes
		if u > 0 {
			if u%60 > 0 {
				w--
				buf[w] = 'm'
				w = fmtInt(buf[:w], u%60)
			}

			u /= 60
			// u is now integer hours
			// Stop at hours because days can be different lengths.
			if u > 0 && u%24 > 0 {
				w--
				buf[w] = 'h'
				w = fmtInt(buf[:w], u%24)
			}

			u /= 24
			if u > 0 && u%365 > 0 {
				w--
				buf[w] = 'd'
				w = fmtInt(buf[:w], u%365)
			}

			u /= 365
			if u > 0 {
				w--
				buf[w] = 'y'
				w = fmtInt(buf[:w], u)
			}
		}
	}

	if neg {
		w--
		buf[w] = '-'
	}

	return string(buf[w:])
}

// fmtFrac formats the fraction of v/10**prec (e.g., ".12345") into the
// tail of buf, omitting trailing zeros. It omits the decimal
// point too when the fraction is 0. It returns the index where the
// output bytes begin and the value v/10**prec.
func fmtFrac(buf []byte, v uint64, prec int) (nw int, nv uint64) {
	// Omit trailing zeros up to and including decimal point.
	w := len(buf)
	print := false
	for i := 0; i < prec; i++ {
		digit := v % 10
		print = print || digit != 0
		if print {
			w--
			buf[w] = byte(digit) + '0'
		}
		v /= 10
	}
	if print {
		w--
		buf[w] = '.'
	}
	return w, v
}

// fmtInt formats v into the tail of buf.
// It returns the index where the output begins.
func fmtInt(buf []byte, v uint64) int {
	w := len(buf)
	if v == 0 {
		w--
		buf[w] = '0'
	} else {
		for v > 0 {
			w--
			buf[w] = byte(v%10) + '0'
			v /= 10
		}
	}
	return w
}
