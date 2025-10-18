package core

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	AppLog *log.Logger
)

func CreateAppLog(dir string, truncated bool) {
	fmt.Printf("Creating app log %s\n", dir)
	flag.Parse()
	err := os.MkdirAll(dir+"/log", 0755)
	if err != nil {
		AppLog = log.New(os.Stdout, "", log.LstdFlags)
		return
	}
	opt := os.O_WRONLY | os.O_CREATE | os.O_APPEND
	if truncated {
		opt = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	}
	file, err := os.OpenFile(dir+"/log/tarantula.log", opt, 0644)
	if err != nil {
		AppLog = log.New(os.Stdout, "", log.LstdFlags)
		return
	}
	AppLog = log.New(file, "", log.LstdFlags|log.Lshortfile)
	AppLog.Println("Initialized app log")
}

func CreateTestLog() {
	flag.Parse()
	AppLog = log.New(os.Stdout, "", log.LstdFlags)
	AppLog.Println("Initialized app log")
}
