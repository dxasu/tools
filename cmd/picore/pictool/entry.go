package pictool

import (
	"bytes"
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/dxasu/tools/cmd/picore/charpaint"
)

type Option struct {
	Opt   string
	Param string
}

func (o Option) Int() int {
	level := 0
	if o.Param != "" {
		level, _ = strconv.Atoi(o.Param)
	}
	return level
}

type PicStruct struct {
	Command string
	Option  Option
	Source  string
	Target  string
}

type picFunc = func(p PicStruct, m image.Image) error

var commandList map[string]picFunc

func writeImg(target, source string, rgba *image.RGBA) {
	f, _ := os.Create(target)
	defer f.Close()
	encode(source, f, rgba)
}

func init() {
	commandList = map[string]picFunc{}
	// 颜色翻转
	commandList["ysfz"] = func(p PicStruct, m image.Image) error {
		newRgba := fzImage(m)
		writeImg(p.Target, p.Source, newRgba)
		return nil
	}

	// 灰度
	commandList["hd"] = func(p PicStruct, m image.Image) error {
		newRgba := hdImage(m)
		writeImg(p.Target, p.Source, newRgba)
		return nil
	}

	// 反色
	commandList["sf"] = func(p PicStruct, m image.Image) error {
		rectWidth := 200
		if p.Option.Param != "" {
			rectWidth, _ = strconv.Atoi(p.Option.Param)
		}
		newRgba := rectImage(m, rectWidth)
		writeImg(p.Target, p.Source, newRgba)
		return nil
	}

	// 图片转字符
	commandList["piczfh"] = func(p PicStruct, m image.Image) error {
		datas := ascllimage(m, p.Option.Int())
		for _, v := range datas {
			fmt.Println(v)
		}
		return nil
	}

	// 字符转字符画
	commandList["zfh"] = func(p PicStruct, _ image.Image) error {
		charList := make([][]string, 0, 8)
		zifu := strings.Fields(p.Source)
		for _, v := range zifu {
			charList = append(charList, charpaint.String(v))
		}
		charpaint.Print(charList...)
		return nil
	}
	// 字符转字符画 彩色
	commandList["zfh2"] = func(p PicStruct, _ image.Image) error {
		charList := make([][]string, 0, 8)
		zifu := strings.Fields(p.Source)
		for _, v := range zifu {
			charList = append(charList, charpaint.Rainbow(v))
		}
		charpaint.Print(charList...)
		return nil
	}
}

func HandlePic(p PicStruct) error {
	picFn, ok := commandList[p.Command]
	if !ok {
		return fmt.Errorf("command:%s not exsit", p.Command)
	}

	ff, err := ioutil.ReadFile(p.Source)
	if err != nil {
		picFn(p, nil)
		return nil
	}
	bbb := bytes.NewBuffer(ff)
	m, _, _ := image.Decode(bbb)
	return picFn(p, m)
}
