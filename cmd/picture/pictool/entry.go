package pictool

import (
	"bytes"
	"image"
	"io/ioutil"
	"os"
	"strconv"
)

type Option struct {
	Opt   string
	Param string
}

type PicStruct struct {
	Command string
	Option  Option
	Source  string
	Target  string
}

func HandlePic(p PicStruct) error {
	ff, _ := ioutil.ReadFile(p.Source)
	bbb := bytes.NewBuffer(ff)
	m, _, _ := image.Decode(bbb)
	if p.Command == "-r" {
		newRgba := fzImage(m)
		f, _ := os.Create(p.Target)
		defer f.Close()
		encode(p.Source, f, newRgba)
	} else if p.Command == "-g" {
		newGray := hdImage(m)
		f, _ := os.Create(p.Target)
		defer f.Close()
		encode(p.Source, f, newGray)
	} else if p.Command == "-c" {
		rectWidth := 200
		if p.Option.Param != "" {
			rectWidth, _ = strconv.Atoi(p.Option.Param)
		}
		newRgba := rectImage(m, rectWidth)
		f, _ := os.Create(p.Target)
		defer f.Close()
		encode(p.Source, f, newRgba)
	} else {
		ascllimage(m, p.Target)
	}
	return nil
}
