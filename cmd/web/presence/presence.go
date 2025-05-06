package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"syscall"
	//"time"

	"gameclustering.com/internal/auth"
)

var service auth.Service

//func debugging(f http.HandlerFunc) http.HandlerFunc {
	//return func(w http.ResponseWriter, r *http.Request) {
		//start := time.Now()
		//defer func() {
			//log.Println(r.URL.Path, time.Since(start))
		//}()
		//f(w, r)
	//}
//}

func bootstrap(host string) {
	service = auth.Service{}
	err := service.Start()
	if err != nil {
		panic(err)
	}
	//http.Handle("/auth", http.HandlerFunc(debugging(auth.AuthHandler)))
	http.Handle("/auth",&service)
	log.Fatal(http.ListenAndServe(host, nil))
}

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
	service.Shutdown()
	fmt.Printf("%s", buf)
	fmt.Println("Exit : ", s)
}
