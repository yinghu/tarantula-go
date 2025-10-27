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
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func AppBootstrap(tcx TarantulaContext) {
	f := conf.Env{}
	err := f.Load(tcx.Config())
	if err != nil {
		fmt.Printf("Config not existed %s\n", err.Error())
		return
	}
	c := cluster.CreateCluster(f, tcx)
	e := event.SocketEndpoint{Endpoint: f.Evp.TcpEndpoint, Service: tcx, OutboundEnabled: f.Evp.OutboundEnabled}
	if f.Evp.Enabled {
		go func() {
			c.Wait()
			e.Open()
		}()
	}
	go func() {
		c.Wait()
		err := tcx.Start(f, c, &e)
		if err != nil {
			core.AppLog.Printf("Error %s\n", err.Error())
		}
		view := c.View()
		for i := range view {
			core.AppLog.Printf("View :%v\n", view[i])
			c.Listener().MemberJoined(view[i])
		}
		http.Handle("/"+tcx.Context()+"/metrics", metricsHandler(tcx.Service().Authenticator(), promhttp.Handler()))
		if tcx.Context() != "admin" {
			http.Handle("/"+tcx.Context()+"/clusteradmin/{cmd}/{cid}", Logging(&AppClusterAdmin{tcx, tcx.Service()}))
			core.AppLog.Printf("Register app cluster admin endpoint %s\n", tcx.Context())
		}
		http.Handle("/", http.HandlerFunc(badRequest))
		core.AppLog.Fatal(http.ListenAndServe(f.HttpBinding, nil))

	}()
	go func() {
		c.Wait()
		core.AppLog.Println("Wating for signal to exit ...")
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
		core.AppLog.Println("Signal to exit")
		tcx.Shutdown()
		c.Quit()
		if f.Evp.Enabled {
			e.Close()
		}
		signal.Stop(sigs)
		close(sigs)
	}()
	c.Join()
}

func badRequest(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	session := core.OnSession{Successful: false, Message: "bad request [" + r.URL.Path + "]", ErrorCode: BAD_REQUEST_CODE}
	w.Write(util.ToJson(session))
}

func invalidToken(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	session := core.OnSession{Successful: false, Message: INVALID_TOKEN_MSG, ErrorCode: INVALID_TOKEN_CODE}
	w.Write(util.ToJson(session))
}
func illegalAccess(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	session := core.OnSession{Successful: false, Message: ILLEGAL_ACCESS_MSG, ErrorCode: ILLEGAL_ACCESS_CODE}
	w.Write(util.ToJson(session))
}

func metricsHandler(auth core.Authenticator, h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tkn := r.Header.Get("Authorization")
		parts := strings.Split(tkn, " ")
		if len(parts) != 2 {
			invalidToken(w, r)
			return
		}
		_, err := auth.ValidateTicket(parts[1])
		if err != nil {
			core.AppLog.Printf("metrics validation failed %s\n", err.Error())
			invalidToken(w, r)
			return
		}
		h.ServeHTTP(w, r)
	}
}

func Logging(s TarantulaApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		var stub int32 = 0
		var code int32 = 0
		defer func() {
			dur := time.Since(start)
			ms := core.ReqMetrics{Path: r.URL.Path, ReqTimed: dur.Milliseconds(), Node: s.Cluster().Local().Name, ReqId: stub, ReqCode: code}
			s.Metrics().WebRequest(ms)
			metrics.HTTP_REQUEST_METRICS.WithLabelValues(r.URL.Path).Observe(dur.Seconds())

		}()
		if s.AccessControl() == PUBLIC_ACCESS_CONTROL {
			s.Request(core.OnSession{}, w, r)
			return
		}
		tkn := r.Header.Get("Authorization")
		parts := strings.Split(tkn, " ")
		if len(parts) != 2 {
			code = int32(ILLEGAL_TOKEN_CODE)
			invalidToken(w, r)
			return
		}
		session, err := s.Authenticator().ValidateToken(parts[1])
		if err != nil {
			code = int32(INVALID_TOKEN_CODE)
			invalidToken(w, r)
			return
		}
		if session.AccessControl < s.AccessControl() {
			stub = session.Stub
			code = int32(ILLEGAL_ACCESS_CODE)
			illegalAccess(w, r)
			return
		}
		stub = session.Stub
		s.Request(session, w, r)
	}
}
