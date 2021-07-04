/*
 Package mp : error helper functions
*/

package mp

import (
	"fmt"
	"os"
)

// CheckErrorExit checks if an error occurred and exits the program if so
func CheckErrorExit(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// ExitWithError exits the program and prints the given text
func ExitWithError(text string) {
	fmt.Println(text)
	os.Exit(1)
}
