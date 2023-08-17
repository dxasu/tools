package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"bay.core/lancet/errdo"
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
	errdo.ExitIf(err)
	err = clipboard.WriteAll(string(out))
	errdo.ExitIf(err)
	if timeout != nil && *timeout > 0 {
		<-time.After(*timeout)
		text, err := clipboard.ReadAll()
		errdo.ExitIf(err)
		if text == string(out) {
			clipboard.WriteAll("")
		}
	}
}

func pasteTo() {
	content, err := clipboard.ReadAll()
	errdo.ExitIf(err)
	fmt.Print(content)
}
