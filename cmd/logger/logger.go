package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

/*
In function design plz add below in inputs function:
- log level
- log module name
- log content
- make log retaintion by size and time creation.
*/

func Logger(logFilePath string, logPrefix string, logContents string, logLevel string) {
	startTime := time.Now()
	file, err := os.OpenFile(fmt.Sprintf("%s%s%s", logFilePath, logLevel, logPrefix), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		defer file.Close()
		log.Println("Cannot open log file with error: ", err)
	}
	log.SetOutput(file)
	log.Println(logContents)
	fmt.Println("Log write latency: ", time.Since(startTime))
}

func directoryNotExist(filePath string) {

}
