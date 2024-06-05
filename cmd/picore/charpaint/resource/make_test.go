package resource

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"testing"
)

func GenerateGoFile(name string, ps []*Painting, height int) {
	jsonVarName := "_" + name
	datafmt := `
package resource
const %s=%s

func init() {
	RegisterFromJson("%s", %s, %d)
}
`
	jb, _ := json.Marshal(ps)
	content := fmt.Sprintf(datafmt, jsonVarName, "`"+string(jb)+"`", name, jsonVarName, height)
	b, e := format.Source([]byte(content))
	if e != nil {
		panic(fmt.Sprintf("Generate gofmt fail, %v", e))
	}

	f, e := os.OpenFile(name+".go", os.O_CREATE|os.O_WRONLY, 0755)
	if e != nil {
		panic(fmt.Sprintf("Generate %s.go open fail, %v", name, e))
	}
	defer f.Close()
	f.Write(b)
}

func TestMake(t *testing.T) {
	charf, e := os.Open("ascii.txt")
	if e != nil {
		t.Fatal(e)
	}
	defer charf.Close()
	charb, e := ioutil.ReadAll(charf)
	if e != nil {
		t.Fatal(e)
	}
	bs := bytes.Split(charb, []byte("\n"))
	chars := make([]byte, len(bs))
	for i := range bs {
		chars[i] = bs[i][0]
	}
	if e = parsePaintingFile("ANSI_Shadow", 6, chars); e != nil {
		t.Fatal(e)
	}
}

func parsePaintingFile(name string, height int, chars []byte) error {
	ps := make([]*Painting, len(chars))
	for i := range ps {
		ps[i] = &Painting{
			Char: chars[i],
		}
	}
	filename := "data-" + name + ".txt"
	f, e := os.Open(filename)
	if e != nil {
		return e
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	for i := range ps {
		ps[i].Data = make([]string, height)
		for j := 0; j < height; j++ {
			line, _, _ := reader.ReadLine()
			ps[i].Data[j] = string(line)
		}
		ps[i].Build()
		reader.ReadLine()
	}
	GenerateGoFile(name, ps, height)
	return nil
}
