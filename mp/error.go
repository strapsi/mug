package mp

import (
	"fmt"
	"os"
)

func CheckErrorExit(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func ExitWithError(text string) {
	fmt.Println(text)
	os.Exit(1)
}
