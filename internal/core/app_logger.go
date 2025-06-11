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
	file, err := os.Create(dir + "/log/tarantula.log")
	if err != nil {
		AppLog = log.New(os.Stdout, "", log.LstdFlags)
		return
	}
	AppLog = log.New(file, "", log.LstdFlags|log.Lshortfile)
	AppLog.Println("Initialized app log")
}
