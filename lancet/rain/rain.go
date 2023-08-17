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
