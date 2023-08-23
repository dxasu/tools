package rain

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func ExitIf(err error) {
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}

func WaitCtrlC() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}

func NeedHelp() bool {
	return len(os.Args) == 2 && (os.Args[1] == "-h" || os.Args[1] == "--help")
}

// Return true if os.Stdin appears to be interactive
func IsInteractive() bool {
	fileInfo, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return fileInfo.Mode()&(os.ModeCharDevice|os.ModeCharDevice) != 0
}
