package main

import (
	"image/color"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/dxasu/qrcode"
	"github.com/spf13/cast"
	unQrCode "github.com/tuotoo/qrcode"
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
		println(unpack(file))
		return
	}

	err := qrcode.WriteColorFile(content, qrcode.Medium, length, color.Black, color.White, "qr.png")
	panicIf(err)

}

func panicIf(err error) {
	if err != nil {
		println(err)
		os.Exit(0)
	}
}

func unpack(qrCodePath string) string {
	var (
		qc  *unQrCode.Matrix
		err error
	)

	if strings.HasPrefix(qrCodePath, "http") {
		_, err = url.Parse(qrCodePath)
		panicIf(err)
		resp, err := http.Get(qrCodePath)
		panicIf(err)
		defer resp.Body.Close()
		qc, err = unQrCode.Decode(resp.Body)
		panicIf(err)
	} else {
		file, err := os.Open(qrCodePath)
		panicIf(err)
		defer file.Close()
		qc, err = unQrCode.Decode(file)
		panicIf(err)
	}
	return qc.Content
}
