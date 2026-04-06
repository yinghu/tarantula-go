package bootstrap

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/protocol"
	"gameclustering.com/internal/util"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var TopicFactoryRegistry = make(map[string]func() core.ProtoTopicFactory)

func Register(name string, fac func() core.ProtoTopicFactory) {
	TopicFactoryRegistry[name] = fac
}

func AppBootstrap(tcx TarantulaContext) {
	Register(event.MESSAGE_TOPIC_NAME, func() core.ProtoTopicFactory { return &event.MessageEventFactory{} })
	Register(event.REGISTER_TOPIC_NAME, func() core.ProtoTopicFactory { return &event.RegisterEventFactory{} })
	Register(event.LOG_TOPIC_NAME, func() core.ProtoTopicFactory { return &event.LogEventFactory{} })
	Register(event.LOGIN_TOPIC_NAME, func() core.ProtoTopicFactory { return &event.LoginEventFactory{} })
	Register(event.REQUEST_TOPIC_NAME, func() core.ProtoTopicFactory { return &event.RequestEventFactory{} })

	f := core.Env{}
	err := f.Load(tcx.Config())
	if err != nil {
		fmt.Printf("Config not existed %s\n", err.Error())
		return
	}
	mountDir := fmt.Sprintf("%s/%s", f.HomeDir, f.GroupName)
	err = os.MkdirAll(mountDir, 0755)
	if err != nil {
		return
	}
	f.LogDir = mountDir
	go func() {
		err := tcx.Start(f)
		if err != nil {
			core.AppLog.Printf("Error %s\n", err.Error())
		}
		http.Handle("/"+tcx.Context()+"/metrics", metricsHandler(tcx.Service().Authenticator(), promhttp.Handler()))
		http.Handle("/", http.HandlerFunc(badRequest))
		core.AppLog.Fatal().Err(http.ListenAndServe(f.HttpBinding, nil))

	}()
	core.AppLog.Println("Wating for signal to exit ...")
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	core.AppLog.Println("Signal to exit")
	tcx.Shutdown()
	signal.Stop(sigs)
	close(sigs)
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
func preflight(w http.ResponseWriter, r *http.Request) {
	//core.AppLog.Debug().Msg("checking options header here")
	defer r.Body.Close()
	w.WriteHeader(http.StatusNoContent)
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
		//var stub int32 = 0
		var code int32 = 0
		defer func() {
			dur := time.Since(start)
			re := protocol.RequestEvent{Path: r.URL.Path, Method: r.Method, Duration: uint64(dur.Milliseconds()), Code: uint32(code)}
			re.DateTime = timestamppb.Now()
			re.Source = r.RemoteAddr
			rf := event.RequestEventFactory{}
			t, err := rf.FromRequestEvent(&re)
			if err != nil {
				fmt.Printf("request event error %s\n", err.Error())
				return
			}
			t.NodeId = s.NodeId()
			t.Tag = s.Context()
			id, err := s.Sequence().Id()
			if err != nil {
				fmt.Printf("request event id error %s\n", err.Error())
				return
			}
			t.Event.Id = uint64(id)
			s.Forward(t)
			//fmt.Printf("send out %s\n", re.Path)
			//ms := core.ReqMetrics{Path: r.URL.Path, ReqTimed: dur.Milliseconds(), Node: s.NodeId(), ReqId: stub, ReqCode: code}
			//s.Metrics().WebRequest(ms)
			//core.HTTP_REQUEST_METRICS.WithLabelValues(r.URL.Path).Observe(dur.Seconds())

		}()
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		if r.Method == "OPTIONS" {
			preflight(w, r)
			return
		}
		if s.AccessControl() == core.PUBLIC_ACCESS_CONTROL {
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
			//stub = session.Stub
			code = int32(ILLEGAL_ACCESS_CODE)
			illegalAccess(w, r)
			return
		}
		//stub = session.Stub
		s.Request(session, w, r)
	}
}
