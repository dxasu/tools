// imagetool project main.go
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/dxasu/tools/cmd/picture/pictool"
	_ "github.com/dxasu/tools/lancet/version"
)

//帮助提示信息
var usage = `Usage: picture COMMAND [OPTION] source_image [OUTPUT]
COMMAND: ysfz(颜色翻转) hd(图片灰度)
Cover old if OUTPUT is empty.
  Options is flow:
	-d		   文件夹 picture hd ./old [./new]
    -n         a number
    -g         图片灰度，可以传入图片缩放的宽度 如: picture hd 1.jpg ./new
    -c         缩放文本
    -t         转成文本
`

//该工具支持将图片色彩反转，图片灰化，图片转为字符画。
func main() {
	args := os.Args //获取用户输入的所有参数
	if len(args) > 6 || len(args) < 3 {
		fmt.Println(usage)
		return
	}

	p := pictool.PicStruct{
		Command: args[1],
		Option:  pictool.Option{},
	}
	if strings.Contains(args[2], "-") {
		p.Option = pictool.Option{Opt: args[2], Param: args[3]}
		p.Source = args[4]
		p.Target = args[5]
	} else {
		p.Source = args[2]
		p.Target = args[3]
	}

	err := pictool.HandlePic(p)
	errdo.errdo.ExitIf(err)
}
