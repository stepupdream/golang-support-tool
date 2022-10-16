package log

import (
	"io"
	"log"
	"os"
)

func LogSetting(filename string) {
	// Open file for write/read logging. (if not, generate one)
	logfile, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	// Generate a Writer for both normal and file output.
	multiLogFile := io.MultiWriter(os.Stdout, logfile)

	// Log output settings (display date and time)
	// Adding log.Llongfile will also output the log output points.
	log.SetFlags(log.Ldate | log.Ltime)

	// Specify log output destination.
	log.SetOutput(multiLogFile)
}
