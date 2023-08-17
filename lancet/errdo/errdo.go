package errdo

import (
	"fmt"
	"os"
)

func ExitIf(err error) {
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}
