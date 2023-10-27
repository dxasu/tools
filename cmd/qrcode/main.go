package main

import (
	"image"
	"image/color"
	"io"
	"net/http"
	"os"
	"strings"

	"bay.core/lancet/rain"
	"github.com/atotto/clipboard"
	"github.com/dxasu/qrcode"
	_ "github.com/dxasu/tools/lancet/version"
	"github.com/makiuchi-d/gozxing"
	uncode "github.com/makiuchi-d/gozxing/qrcode"
	"github.com/spf13/cast"
)

func main() {
	if len(os.Args) < 2 || rain.NeedHelp() {
		println(`please input http://xxxxx
-c or "-cs 256" with copy data from clipboard
-u xxx.png | http://xxxx.png with unpack from qrcode`)
		return
	}
	content := os.Args[1]
	length := 256
	if content == "-c" || content == "-cs" && len(os.Args) >= 3 {
		var err error
		content, err = clipboard.ReadAll()
		rain.ExitIf(err)

		if content == "-cs" {
			length = cast.ToInt(os.Args[2])
		}
	} else if content == "-u" && len(os.Args) >= 3 {
		file := os.Args[2]
		println(decodeFile(file))
		return
	}

	err := qrcode.WriteColorFile(content, qrcode.Medium, length, color.White, color.Black, "qr.png")
	rain.ExitIf(err)

}

func decodeFile(qrCodePath string) string {
	var r io.Reader
	if strings.HasPrefix(qrCodePath, "http") {
		resp, err := http.Get(qrCodePath)
		rain.ExitIf(err)
		defer resp.Body.Close()
		r = resp.Body
	} else {
		file, err := os.Open(qrCodePath)
		rain.ExitIf(err)
		defer file.Close()
		r = file
	}
	content := decodeReader(r)
	return content
}

func decodeReader(file io.Reader) string {
	img, _, err := image.Decode(file)
	rain.ExitIf(err)
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	rain.ExitIf(err)
	// decode image
	qrReader := uncode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)
	rain.ExitIf(err)
	return result.String()
}
