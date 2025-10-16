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

func CreateAppLog(dir string) {
	fmt.Printf("Creating app log %s\n", dir)
	flag.Parse()
	err := os.MkdirAll(dir+"/log", 0755)
	if err != nil {
		AppLog = log.New(os.Stdout, "", log.LstdFlags)
		return
	}
	file, err := os.OpenFile(dir+"/log/tarantula.log", os.O_CREATE|os.O_APPEND, 0644)
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
