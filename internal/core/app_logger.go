package core

import (
	"flag"
	"log"
	"os"
)

var (
	AppLog *log.Logger
)

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		AppLog = log.New(os.Stdout, "", log.LstdFlags)
		return
	}
	err = os.MkdirAll(homeDir+"/log", 0755)
	if err != nil {
		AppLog = log.New(os.Stdout, "", log.LstdFlags)
		return
	}
	flag.Parse()
	file, err := os.Create(homeDir + "/log/tarantula.log")
	if err != nil {
		AppLog = log.New(os.Stdout, "", log.LstdFlags)
		return
	}
	AppLog = log.New(file, "", log.LstdFlags|log.Lshortfile)
	AppLog.Println("Initialize app log")
}
