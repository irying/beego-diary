package logging

import (
	"fmt"
	"time"
	"os"
	"log"
)

var (
	LogSavePath = "runtime/logs/"
	LogSaveName = "log"
	LogFileExt  = "log"
	TimeFormat  = "20060102"
)

func getLogFilePath() string {
	dir, _ := os.Getwd()
	path := dir + "/" + LogSavePath
	log.Printf("Log file path : %s", path)

	return fmt.Sprintf("%s", path)
}

func getLogFileFullPath() string {
	prefixPath := getLogFilePath()
	suffixPath := fmt.Sprintf("%s%s.%s", LogSaveName, time.Now().Format(TimeFormat), LogFileExt)

	return fmt.Sprintf("%s%s", prefixPath, suffixPath)
}

// todo openfile
func openLogFile(filePath string) *os.File {
	_, err := os.Stat(filePath)
	switch {
	case os.IsNotExist(err):
		mkDir(getLogFilePath())
	case os.IsPermission(err):
		log.Fatalf("Permission :%v", err)

	}

	handle, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Fail to open file :%v", err)
	}

	return handle
}

func mkDir(filePath string) {
	err := os.MkdirAll(filePath, os.ModePerm)
	if err != nil {
		panic(err)
	}
}
