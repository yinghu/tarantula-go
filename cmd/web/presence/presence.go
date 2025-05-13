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
	"gameclustering.com/internal/conf"
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

func bootstrap(f conf.Env) {
	service = auth.Service{}
	err := service.Start(f)
	if err != nil {
		panic(err)
	}
	//http.Handle("/auth", http.HandlerFunc(debugging(auth.AuthHandler)))
	http.Handle("/auth", &service)
	log.Fatal(http.ListenAndServe(f.HttpEndpoint, nil))
}

func main() {
	f := conf.Env{}
	f.Load("/etc/tarantula/presence-conf.json")
	c := cluster.NewEtc(f.GroupName, f.PartitionNumber, f.EtcdEndpoints, cluster.Node{Name: f.NodeName, HttpEndpoint: f.HttpEndpoint, TcpEndpoint: f.TcpEndpoint})
	go func() {
		c.Started.Wait()
		for v := range c.View() {
			fmt.Printf("View :%v\n", v)
		}
		bootstrap(f)
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
