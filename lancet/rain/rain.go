package rain

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
)

func ExitIf(err error) {
	if err != nil {
		println(err.Error())
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

func OpenBrower(uri string) error {
	// 不同平台启动指令不同
	var commands = map[string]string{
		"windows": "explorer",
		"darwin":  "open",
		"linux":   "xdg-open",
	}
	run, ok := commands[runtime.GOOS]
	if !ok {
		return fmt.Errorf("invalid platform: %s", runtime.GOOS)
	}
	cmd := exec.Command(run, uri)
	return cmd.Run()
}
