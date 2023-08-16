package main

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/dxasu/qrcode"
	uncode "github.com/makiuchi-d/gozxing/qrcode"

	"github.com/makiuchi-d/gozxing"
	"github.com/spf13/cast"
)

func main() {
	var content string
	if len(os.Args) >= 2 {
		content = os.Args[1]
	} else {
		panic(`args empty, input http://xxxxx
		 -c or "-cs 256" with copy data from clipboard
		 -u xxx.png | http://xxxx.png with unpack from qrcode`)
	}

	length := 256
	if content == "-c" || content == "-cs" && len(os.Args) >= 3 {
		var err error
		content, err = clipboard.ReadAll()
		panicIf(err)

		if content == "-cs" {
			length = cast.ToInt(os.Args[2])
		}
	} else if content == "-u" && len(os.Args) >= 3 {
		file := os.Args[2]
		println(decodeFile(file))
		return
	}

	err := qrcode.WriteColorFile(content, qrcode.Medium, length, color.White, color.Black, "qr.png")
	panicIf(err)

}

func panicIf(err error) {
	if err != nil {
		fmt.Printf("%+v", err)
		os.Exit(0)
	}
}

func decodeFile(qrCodePath string) string {
	var r io.Reader
	if strings.HasPrefix(qrCodePath, "http") {
		resp, err := http.Get(qrCodePath)
		panicIf(err)
		defer resp.Body.Close()
		r = resp.Body
	} else {
		file, err := os.Open(qrCodePath)
		panicIf(err)
		defer file.Close()
		r = file
	}
	content := decodeReader(r)
	return content
}

func decodeReader(file io.Reader) string {
	img, _, err := image.Decode(file)
	panicIf(err)
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	panicIf(err)
	// decode image
	qrReader := uncode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)
	panicIf(err)
	return result.String()
}
