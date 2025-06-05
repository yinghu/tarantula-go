package bootstrap

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"gameclustering.com/internal/cluster"
	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/metrics"
	"gameclustering.com/internal/util"
)

func AppBootstrap(service TarantulaContext) {
	f := conf.Env{}
	err := f.Load(service.Config())
	if err != nil {
		fmt.Printf("Config not existed %s\n", err.Error())
		return
	}
	c := cluster.NewEtc(f.GroupName, f.EtcdEndpoints, cluster.Node{Name: f.NodeName, HttpEndpoint: f.HttpEndpoint, TcpEndpoint: f.Evp.TcpEndpoint})
	c.Kyl = service
	e := event.Endpoint{TcpEndpoint: f.Evp.TcpEndpoint, Service: service}
	if f.Evp.Enabled {
		go func() {
			c.Started.Wait()
			e.Open()
		}()
	}
	go func() {
		c.Started.Wait()
		for v := range c.View() {
			fmt.Printf("View :%v\n", v)
		}
		err := service.Start(f, &c)
		if err != nil {
			//panic(err)
			fmt.Printf("Error %s\n", err.Error())
		}
	}()

	go func() {
		c.Started.Wait()
		fmt.Println("Wating for signal to exit ...")
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
		service.Shutdown()
		c.Quit <- true
		if f.Evp.Enabled {
			e.Close()
		}
		signal.Stop(sigs)
		close(sigs)
	}()
	c.Join()
}

func invalidToken(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	session := core.OnSession{Successful: false, Message: "invalid token"}
	w.Write(util.ToJson(session))
}
func illegalAccess(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	session := core.OnSession{Successful: false, Message: "illegal access"}
	w.Write(util.ToJson(session))
}

func Logging(s TarantulaApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			dur := time.Since(start)
			ms := metrics.ReqMetrics{Path: r.URL.Path, ReqTimed: dur.Milliseconds(), Node: s.Cluster().Local().Name}
			s.Metrics().WebRequest(ms)
		}()
		if s.AccessControl() == PUBLIC_ACCESS_CONTROL {
			s.Request(core.OnSession{}, w, r)
			return
		}
		tkn := r.Header.Get("Authorization")
		parts := strings.Split(tkn, " ")
		if len(parts) != 2 {
			invalidToken(w, r)
			return
		}
		session, err := s.Authenticator().ValidateToken(parts[1])
		if err != nil {
			invalidToken(w, r)
			return
		}
		if session.AccessControl < s.AccessControl() {
			illegalAccess(w, r)
			return
		}
		s.Request(session, w, r)
	}
}
