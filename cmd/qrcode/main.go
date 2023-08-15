package main

import (
	"fmt"
	"os"

	"github.com/dxasu/qrcode"
)

func main() {
	if len(os.Args) == 0 {
		fmt.Errorf("args empty")
		return
	}

	url := os.Args[0]
	q, _ := qrcode.New(url, qrcode.High)
}
