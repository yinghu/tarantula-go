package bootstrap

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/item"
	"gameclustering.com/internal/util"
)

type AppClusterAdmin struct {
	event.EventService
	TarantulaService
}

func (s *AppClusterAdmin) AccessControl() int32 {
	return ADMIN_ACCESS_CONTROL
}

func (s *AppClusterAdmin) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	session := core.OnSession{Successful: true, Message: "app cluster admin [" + s.Cluster().Group() + "]"}
	defer func() {
		w.WriteHeader(http.StatusOK)
		w.Write(util.ToJson(session))
		r.Body.Close()
	}()
	cmd := r.PathValue("cmd")
	cid, err := strconv.ParseInt(r.PathValue("cid"), 10, 64)
	if err != nil {
		session = core.OnSession{Successful: false, Message: "cid should be a number"}
		return
	}
	switch cmd {
	case "join":
		var join core.Node
		json.NewDecoder(r.Body).Decode(&join)
		nd := s.convert(join)
		s.Cluster().OnJoin(nd)
		s.startEvent(nd)
	case "left":
		var left core.Node
		json.NewDecoder(r.Body).Decode(&left)
		s.Cluster().OnLeave(s.convert(left))
	case "update":
		var update item.KVUpdate
		json.NewDecoder(r.Body).Decode(&update)
		s.ItemService().Loader().Reload(update)
	case "schedule":
		var update item.KVUpdate
		json.NewDecoder(r.Body).Decode(&update)
		if s.ItemListener() == nil {
			return
		}
		s.dispatch(update)
	case "event":
		e := event.CreateEvent(int(cid), s)
		json.NewDecoder(r.Body).Decode(&e)
		s.OnEvent(e)
	default:
		core.AppLog.Printf("cmd not supported %s\n", cmd)
		session = core.OnSession{Successful: false, Message: "cmd not supported [" + cmd + "]"}
	}
}

func (s *AppClusterAdmin) convert(node core.Node) core.Node {
	gparts := strings.Split(s.Cluster().Group(), "/")
	node.Name = strings.Replace(node.Name, "admin", gparts[1], 1)
	lparts := strings.Split(s.Cluster().Local().TcpEndpoint, ":")
	rparts := strings.Split(node.TcpEndpoint, ":")
	node.TcpEndpoint = rparts[0] + ":" + rparts[1] + ":" + lparts[2]
	return node
}
func (s *AppClusterAdmin) dispatch(kv item.KVUpdate) {
	itemId, err := strconv.ParseInt(kv.Key, 10, 64)
	if err != nil {
		core.AppLog.Printf("Key should be int64 %s\n", kv.Key)
	}
	if kv.IsCreate || kv.IsModify {
		var reg item.ConfigRegistration
		err = json.Unmarshal([]byte(kv.Value), &reg)
		if err != nil {
			core.AppLog.Printf("Value should be json format %v\n", kv.Value)
			s.send(err)
			return
		}
		if reg.ItemId != itemId {
			core.AppLog.Printf("Key not matched %d : %d\n", itemId, reg.ItemId)
			s.send(fmt.Errorf("key not matched %d : %d", itemId, reg.ItemId))
			return
		}
		ins, err := s.ItemService().Loader().Load(reg.ItemId)
		if err != nil {
			s.send(err)
			return
		}
		core.AppLog.Printf("Item registered %d\n", ins.Id)
		s.ItemListener().OnRegister(ins)
		return
	}
	core.AppLog.Printf("Item released %d\n", itemId)
	ins, err := s.ItemService().Loader().Load(itemId)
	if err != nil {
		s.send(err)
		return
	}
	s.ItemListener().OnRelease(ins)
}

func (s *AppClusterAdmin) startEvent(n core.Node) {
	s.TarantulaService.NodeStarted(n)
	//ste := event.NodeStartEvent{NodeName: s.Cluster().Local().Name, StartTime: time.Now()}
	//id, _ := s.Sequence().Id()
	//ste.Id = id
	//ste.Topic("message")
	//s.Send(&ste)
}

func (s *AppClusterAdmin) send(err error) {
	msg := event.MessageEvent{Title: "error", Message: err.Error(), Source: s.Cluster().Group(), DateTime: time.Now()}
	id, _ := s.Sequence().Id()
	msg.Id = id
	msg.Topic("message")
	s.Send(&msg)
}
