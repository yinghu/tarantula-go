package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gameclustering.com/internal/cluster"
	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/metrics"
)

var service Service

func debugging(s *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		action := r.Header.Get("Tarantula-action")
		defer func() {
			dur := time.Since(start)
			ms := metrics.ReqMetrics{Path: r.URL.Path + "/" + action, ReqTimed: dur.Milliseconds(), Node: s.Cluster.Local.Name}
			s.SaveMetrics(&ms)
		}()
		s.ServeHTTP(w, r)
	}
}

func bootstrap(f conf.Env, c *cluster.Etc) {
	service = Service{Cluster: c}
	err := service.Start(f)
	if err != nil {
		panic(err)
	}
	http.Handle("/auth", http.HandlerFunc(debugging(&service)))
	log.Fatal(http.ListenAndServe(f.HttpEndpoint, nil))
}

func main() {
	f := conf.Env{}
	f.Load("/etc/tarantula/presence-conf.json")
	c := cluster.NewEtc(f.GroupName, f.PartitionNumber, f.EtcdEndpoints, cluster.Node{Name: f.NodeName, HttpEndpoint: f.HttpEndpoint, TcpEndpoint: f.TcpEndpoint})
	e := event.Endpoint{TcpEndpoint: f.TcpEndpoint, Service: &service, ReadBufferSize: f.TcpReadBufferSize}
	go func() {
		c.Started.Wait()
		for v := range c.View() {
			fmt.Printf("View :%v\n", v)
		}
		bootstrap(f, &c)
	}()
	go func() {
		c.Started.Wait()
		e.Open()
	}()
	go func() {
		c.Started.Wait()
		fmt.Println("Wating for signal to exit ...")
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
		signal.Stop(sigs)
		c.Quit <- true
		service.Shutdown()
		e.Close()
		close(sigs)
	}()
	c.Join()
}
