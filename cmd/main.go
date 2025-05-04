package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"syscall"
)



func main() {
	go bootstrap(":8080")
	fmt.Println("Started : ")
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	s := <-sigs
	signal.Stop(sigs)
	close(sigs)
	debug.PrintStack()
	buf := make([]byte, 1<<16)
	runtime.Stack(buf, true)
	fmt.Printf("%s", buf)
	fmt.Println("Exit : ", s)
}
