package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/atotto/clipboard"
	_ "github.com/dxasu/tools/lancet/version"
)

func main() {
	timeout := flag.Duration("t", 0, "Erase clipboard after timeout.  Durations are specified like \"20s\" or \"2h45m\".  0 (default) means never erase.")
	paste := flag.Bool("p", false, "paste into stdout from clipboard")
	flag.Parse()

	if *paste {
		pasteTo()
		return
	}

	out, err := ioutil.ReadAll(os.Stdin)
	exitIf(err)
	err = clipboard.WriteAll(string(out))
	exitIf(err)
	if timeout != nil && *timeout > 0 {
		<-time.After(*timeout)
		text, err := clipboard.ReadAll()
		exitIf(err)
		if text == string(out) {
			clipboard.WriteAll("")
		}
	}
}

func exitIf(err error) {
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}

func pasteTo() {
	content, err := clipboard.ReadAll()
	exitIf(err)
	fmt.Print(content)
}
