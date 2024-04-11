package resource

import (
	"encoding/json"
	"fmt"
	"sync"
)

type Painting struct {
	Char    byte
	Data    []string
	Runes   [][]rune
	Width   int
	Height  int
	Borders [][2]int
	Valid   bool
}

func (p Painting) DebugPrint() {
	fmt.Println("========================")
	fmt.Printf("Char: %c %d\n", p.Char, p.Char)
	if !p.Valid {
		fmt.Println("---[EMPTY]---")
		return
	}
	fmt.Printf("Height: %d\n", p.Height)
	fmt.Printf("Width : %d\n", p.Width)
	for i := range p.Data {
		fmt.Printf("%3d|%s|%3d\n", p.Borders[i][0], p.Data[i], p.Borders[i][1])
	}
}

func (p *Painting) Build() {
	maxWidth := 0
	p.Borders = make([][2]int, len(p.Data))
	p.Runes = make([][]rune, len(p.Data))
	for i := range p.Data {
		w := 0
		var border [2]int
		rs := []rune(p.Data[i])
		for j := len(rs) - 1; j >= 0; j-- {
			if rs[j] != ' ' {
				w = j + 1
				border[1] = w
				break
			}
		}
		if w != 0 {
			for j := 0; j < w; j++ {
				if rs[j] != ' ' {
					border[0] = j
					break
				}
			}
		}
		if w > maxWidth {
			maxWidth = w
		}
		p.Runes[i] = rs
		p.Borders[i] = border
	}
	if maxWidth > 0 {
		p.Valid = true
	} else {
		p.Data = nil
		p.Runes = nil
		return
	}
	p.Height = len(p.Data)
	p.Width = maxWidth
	for i := range p.Data {
		p.Runes[i] = p.Runes[i][:p.Width]
		p.Data[i] = string(p.Runes[i])
	}
}

var (
	pdata     map[string]PM
	pdataLock sync.Mutex
)

type PM struct {
	M    map[byte]Painting
	H    int
	AvgW int
}

func (p PM) Dump() {
	for _, pt := range p.M {
		pt.DebugPrint()
	}
}

func RegisterFromJson(name, j string, h int) {
	var ps []Painting
	if e := json.Unmarshal([]byte(j), &ps); e != nil {
		panic(fmt.Errorf("RegisterFromJson fail, %v", e))
	}
	var pm PM
	pm.M = make(map[byte]Painting)
	pm.H = h
	for i := range ps {
		if !ps[i].Valid {
			continue
		}
		pm.M[ps[i].Char] = ps[i]
		pm.AvgW += ps[i].Width
	}
	pm.AvgW /= len(pm.M)
	pdataLock.Lock()
	defer pdataLock.Unlock()
	if pdata == nil {
		pdata = make(map[string]PM)
	}
	pdata[name] = pm
}

func Get(name string) (PM, bool) {
	v, ok := pdata[name]
	return v, ok
}
