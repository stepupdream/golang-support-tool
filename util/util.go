package util

import (
	"bufio"
	"fmt"
	"os"

	"github.com/cheggaaa/pb/v3"
)

func KeyWait(message string) {
	fmt.Println(message)
	bufio.NewScanner(os.Stdin).Scan()
}

func StartProgressBar(totalCount int) *pb.ProgressBar {
	bar := pb.Simple.Start(totalCount)
	bar.SetMaxWidth(80)

	return bar
}
