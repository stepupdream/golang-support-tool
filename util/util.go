package util

import (
	"bufio"
	"fmt"
	"os"
)

func KeyWait(message string) {
	fmt.Println(message)
	bufio.NewScanner(os.Stdin).Scan()
}
