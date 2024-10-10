package logging

import (
	"io"
	"log"
	"os"
)

func InitLogging(logFileName string) io.Writer {
	logFile, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0640)
	if err != nil {
		panic("failed to open log file")
	}
	log.Println("Created log file!")

	return io.MultiWriter(os.Stdout, logFile)
}
