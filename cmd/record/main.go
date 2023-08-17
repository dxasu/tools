package main

import (
	"fmt"
	"os"
	"time"

	"os/exec"
	"path/filepath"

	"bay.core/lancet/rain"
	_ "github.com/dxasu/tools/lancet/version"
	"github.com/spf13/cast"

	"github.com/dxasu/corekey"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] != "-o" && os.Args[1] != "-l" {
		println(`args empty
-o for open data of record
-l 1h record start`)
		return
	}
	if os.Args[1] == "-o" {
		path := filepath.FromSlash(corekey.GetAppDataPath())
		exec.Command("explorer", path).Run()
		return
	}
	var long time.Duration = 3600
	if len(os.Args) > 2 {
		long = cast.ToDuration(os.Args[2]) / 1e9
	}
	corekey.PcListen(fmt.Sprintf("core_dump_v%d_linux.tmp", long), 0)
	println("Wait Ctrl + C")
	rain.WaitCtrlC()
}
