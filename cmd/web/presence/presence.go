package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"gameclustering.com/internal/auth"
	"gameclustering.com/internal/cluster"
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
	http.Handle("/auth", &service)
	log.Fatal(http.ListenAndServe(host, nil))
}

func main() {
	c := cluster.NewEtc("presence", []string{"192.168.1.7:2379"}, cluster.Node{Name: "a01", HttpEndpoint: "http://192.168.1.11:8080", TcpEndpoint: "tcp://192.168.1.11:5000"})
	go func() {
		c.Started.Wait()
		for v := range c.View() {
			fmt.Printf("View :%v\n", v)
		}
		bootstrap(":8080")
	}()
	go func() {
		c.Started.Wait()
		fmt.Println("Started : ")
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
		signal.Stop(sigs)
		c.Quit <- true
		close(sigs)
		service.Shutdown()
	}()
	c.Join()
}
