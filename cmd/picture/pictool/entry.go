package pictool

import (
	"bytes"
	"fmt"
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

type picFunc = func(p PicStruct, m image.Image) error

var commandList map[string]picFunc

func init() {
	commandList = map[string]picFunc{}
	commandList["ysfz"] = func(p PicStruct, m image.Image) error {
		newRgba := fzImage(m)
		f, _ := os.Create(p.Target)
		defer f.Close()
		encode(p.Source, f, newRgba)
		return nil
	}

	commandList["hd"] = func(p PicStruct, m image.Image) error {
		newGray := hdImage(m)
		f, _ := os.Create(p.Target)
		defer f.Close()
		encode(p.Source, f, newGray)
		return nil
	}

	commandList["sf"] = func(p PicStruct, m image.Image) error {
		rectWidth := 200
		if p.Option.Param != "" {
			rectWidth, _ = strconv.Atoi(p.Option.Param)
		}
		newRgba := rectImage(m, rectWidth)
		f, _ := os.Create(p.Target)
		defer f.Close()
		encode(p.Source, f, newRgba)
		return nil
	}

	commandList["zc"] = func(p PicStruct, m image.Image) error {
		ascllimage(m, p.Target)
		return nil
	}
}

func HandlePic(p PicStruct) error {
	picFn, ok := commandList[p.Command]
	if !ok {
		return fmt.Errorf("command:%s not exsit", p.Command)
	}

	ff, _ := ioutil.ReadFile(p.Source)
	bbb := bytes.NewBuffer(ff)
	m, _, _ := image.Decode(bbb)
	return picFn(p, m)
}
