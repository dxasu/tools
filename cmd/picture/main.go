// imagetool project main.go
package main

import (
	"fmt"
	"os"
	"strings"

	"bay.core/lancet/rain"

	"github.com/dxasu/tools/cmd/picture/pictool"
	_ "github.com/dxasu/tools/lancet/version"
)

//帮助提示信息
var usage = `Usage: picture COMMAND [OPTION] source_image [OUTPUT]
COMMAND: ysfz(颜色翻转) hd(图片灰度) piczfh(图片转字符画) zfh/zfh2(字符画) 
Example: picture ysfz source_image.jpg
Options as follows:
  -l         level, default is 0, range is 0-100
  -h         --help, show this help message
`

//该工具支持将图片色彩反转，图片灰化，图片转为字符画。
func main() {
	args := os.Args //获取用户输入的所有参数
	if len(args) > 6 || len(args) < 3 || rain.NeedHelp() {
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
		p.Target = "new_" + args[4]
	} else {
		p.Source = args[2]
		p.Target = "new_" + args[2]
	}

	err := pictool.HandlePic(p)
	rain.ExitIf(err)
}
