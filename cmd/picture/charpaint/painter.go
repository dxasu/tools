// author: dxasu@foxmail.com
// date: 2021-08-11 16:23:58
// description: draw text on terminal using ANSI escape sequence
// version: 1.0.0
package charpaint

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/dxasu/tools/cmd/picture/charpaint/resource"
)

var (
	font = "ANSI_Shadow"
)

const (
	COLOR_RESET       = "\033[0;0m"
	COLOR_BLACK       = "\033[0;30m"
	COLOR_RED         = "\033[0;31m"
	COLOR_GREEN       = "\033[0;32m"
	COLOR_ORANGE      = "\033[0;33m"
	COLOR_BLUE        = "\033[0;34m"
	COLOR_PURPLE      = "\033[0;35m"
	COLOR_CYAN        = "\033[0;36m"
	COLOR_LIGHTGRAY   = "\033[0;37m"
	COLOR_DARKGRAY    = "\033[1;30m"
	COLOR_LIGHTRED    = "\033[1;31m"
	COLOR_LIGHTGREEN  = "\033[1;32m"
	COLOR_YELLOW      = "\033[1;33m"
	COLOR_LIGHTBLUE   = "\033[1;34m"
	COLOR_LIGHTPURPLE = "\033[1;35m"
	COLOR_LIGHTCYAN   = "\033[1;36m"
	COLOR_WHITE       = "\033[1;37m"
)

var (
	rainbowlist = []string{COLOR_RESET, COLOR_RED, COLOR_GREEN, COLOR_ORANGE, COLOR_BLUE, COLOR_PURPLE, COLOR_CYAN}
)

func SetFont(name string) {
	if _, ok := resource.Get(name); !ok {
		panic(fmt.Errorf("invalid font[%s]", name))
	}
	font = name
}

func NewPainter(name ...string) Painter {
	if len(name) == 0 {
		m, _ := resource.Get(font)
		return Painter{PM: m}
	} else {
		m, ok := resource.Get(name[0])
		if !ok {
			panic(fmt.Errorf("invalid font[%s]", name))
		}
		return Painter{PM: m}
	}
}

type Painter struct {
	resource.PM
	blank []string
	space []string
}

func (p Painter) RainbowStagger(s string) []string {
	l := 1 // stagger level
	colors := make([]int, len(s))
	bs := make([]strings.Builder, p.H+l)
	mark := rand.Intn(1)
	for i := range s {
		colors[i] = rand.Intn(len(rainbowlist))
		for i > 0 && colors[i] == colors[i-1] {
			colors[i] = rand.Intn(len(rainbowlist))
		}
		pt := p.M[s[i]]
		var data []string
		var lineWidth int
		if s[i] == ' ' {
			data = p.Space()
			lineWidth = p.AvgW
		} else if pt.Valid {
			data = pt.Data
			lineWidth = pt.Width
		} else {
			data = p.Blank()
			lineWidth = p.AvgW
		}
		level := (i + mark) % 2
		if level == 1 {
			bs[0].WriteString(rainbowlist[colors[i]] + Repeating(' ', lineWidth) + " ")
		}
		for j := 0; j < p.H; j++ {
			bs[level+j].WriteString(rainbowlist[colors[i]] + data[j] + " ")
		}
		if level == 0 {
			bs[len(bs)-1].WriteString(rainbowlist[colors[i]] + Repeating(' ', lineWidth) + " ")
		}
	}
	ss := make([]string, p.H+l)
	for i := 0; i < p.H+l; i++ {
		ss[i] = bs[i].String() + COLOR_RESET
	}
	return ss
}

func (p Painter) ColorLoop(s string, colors []string) []string {
	bs := make([]strings.Builder, p.H)
	for i := range s {
		pt := p.M[s[i]]
		var data []string
		if s[i] == ' ' {
			data = p.Space()
		} else if pt.Valid {
			data = pt.Data
		} else {
			data = p.Blank()
		}
		for j := 0; j < p.H; j++ {
			bs[j].WriteString(colors[i%len(colors)] + data[j] + " ")
		}
	}
	ss := make([]string, p.H)
	for i := 0; i < p.H; i++ {
		ss[i] = bs[i].String() + COLOR_RESET
	}
	return ss
}

func (p Painter) Rainbow(s string) []string {
	colors := make([]int, len(s))
	bs := make([]strings.Builder, p.H)
	for i := range s {
		colors[i] = rand.Intn(len(rainbowlist))
		for i > 0 && colors[i] == colors[i-1] {
			colors[i] = rand.Intn(len(rainbowlist))
		}
		pt := p.M[s[i]]
		var data []string
		if s[i] == ' ' {
			data = p.Space()
		} else if pt.Valid {
			data = pt.Data
		} else {
			data = p.Blank()
		}
		for j := 0; j < p.H; j++ {
			bs[j].WriteString(rainbowlist[colors[i]] + data[j] + " ")
		}
	}
	ss := make([]string, p.H)
	for i := 0; i < p.H; i++ {
		ss[i] = bs[i].String() + COLOR_RESET
	}
	return ss
}

func (p Painter) Color(s, color string) []string {
	bs := make([]strings.Builder, p.H)
	for i := range s {
		pt := p.M[s[i]]
		var data []string
		if s[i] == ' ' {
			data = p.Space()
		} else if pt.Valid {
			data = pt.Data
		} else {
			data = p.Blank()
		}
		for j := 0; j < p.H; j++ {
			bs[j].WriteString(data[j] + " ")
		}
	}
	ss := make([]string, p.H)
	for i := 0; i < p.H; i++ {
		ss[i] = color + bs[i].String() + COLOR_RESET
	}
	return ss
}

func (p Painter) String(s string) []string {
	bs := make([]strings.Builder, p.H)
	for i := range s {
		pt := p.M[s[i]]
		var data []string
		if s[i] == ' ' {
			data = p.Space()
		} else if pt.Valid {
			data = pt.Data
		} else {
			data = p.Blank()
		}
		for j := 0; j < p.H; j++ {
			bs[j].WriteString(data[j] + " ")
		}
	}
	ss := make([]string, p.H)
	for i := 0; i < p.H; i++ {
		ss[i] = bs[i].String()
	}
	return ss
}

func (p *Painter) Space() []string {
	if p.space != nil {
		return p.space
	}
	s := Repeating(' ', p.AvgW)
	ss := make([]string, p.H)
	for i := range ss {
		ss[i] = s
	}
	p.space = ss
	return ss
}

func (p *Painter) Blank() []string {
	if p.blank != nil {
		return p.blank
	}
	s := Repeating('█', p.AvgW)
	ss := make([]string, p.H)
	for i := range ss {
		ss[i] = s
	}
	p.blank = ss
	return ss
}

func Color(s, color string) []string {
	return NewPainter(font).Color(s, color)
}

func String(s string) []string {
	return NewPainter(font).String(s)
}

func Rainbow(s string) []string {
	return NewPainter(font).Rainbow(s)
}

func RainbowStagger(s string) []string {
	return NewPainter(font).RainbowStagger(s)
}

func ColorLoop(s string, colors []string) []string {
	return NewPainter(font).ColorLoop(s, colors)
}

func Repeating(c rune, n int) string {
	s := make([]rune, n)
	for i := range s {
		s[i] = c
	}
	return string(s)
}

func (p Painter) Join(joint string, in ...[]string) []string {
	if len(in) == 0 {
		return nil
	}
	h := p.H
	cs := p.String(joint)
	bs := make([]bytes.Buffer, h)
	for i := range in {
		for j := 0; j < h; j++ {
			bs[j].WriteString(in[i][j])
			if i != len(in)-1 {
				bs[j].WriteString(cs[j])
			}
		}
	}
	ss := make([]string, h)
	for i := range ss {
		ss[i] = bs[i].String()
	}
	return ss
}

func Join(joint string, in ...[]string) []string {
	return NewPainter(font).Join(joint, in...)
}

func Print(ss ...[]string) {
	s := Join(" ", ss...)
	for i := range s {
		fmt.Println(s[i])
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
