package errdo

import (
	"fmt"
	"os"
)

func errdo.ExitIf(err error) {
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}
