package bootstrap

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"gameclustering.com/internal/cluster"
	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/event"
)

func AppBootstrap(service TarantulaContext) {
	f := conf.Env{}
	f.Load(service.Config())
	c := cluster.NewEtc(f.GroupName, f.PartitionNumber, f.EtcdEndpoints, cluster.Node{Name: f.NodeName, HttpEndpoint: f.HttpEndpoint, TcpEndpoint: f.TcpEndpoint})
	e := event.Endpoint{TcpEndpoint: f.TcpEndpoint, Service: service, ReadBufferSize: f.TcpReadBufferSize}
	go func() {
		c.Started.Wait()
		for v := range c.View() {
			fmt.Printf("View :%v\n", v)
		}
		err := service.Start(f, &c)
		if err != nil {
			fmt.Printf("Error %s\n", err.Error())
		}
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
		service.Shutdown()
		c.Quit <- true
		e.Close()
		signal.Stop(sigs)
		close(sigs)
	}()
	c.Join()
}
