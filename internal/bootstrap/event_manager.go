package bootstrap

import (
	"context"
	"fmt"
	"io"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/protocol"
)

func (s *AppManager) Create(classId uint32, topic string) (core.Event, error) {
	return nil, nil
}

func (s *AppManager) VerifyTicket(ticket string) (core.OnSession, error) {
	session, err := s.auth.ValidateTicket(ticket)
	if err != nil {
		return session, err
	}
	if session.AccessControl < core.ADMIN_ACCESS_CONTROL {
		return session, fmt.Errorf("admin access control required %d", session.AccessControl)
	}
	return session, nil
}

func (s *AppManager) OnEvent(e core.Event) {
	core.AppLog.Debug().Msgf("event %v", e)
}

func (s *AppManager) OnError(e core.Event, err error) {

}

func (s *AppManager) Publish(e core.Event) error {
	k, v, err := core.Export(e, 200)
	if err != nil {
		return err
	}
	data := protocol.Data{Key: k, Value: v, Header: &protocol.Header{FactoryId: e.FactoryId(), ClassId: e.ClassId()}}
	dsp := protocol.NewPostofficeServiceClient(s.rpc)
	resp, err := dsp.Publish(context.Background(), &protocol.Topic{Tag: e.ETag(), Name: e.Topic(), Event: &data})
	if err != nil {
		return err
	}
	core.AppLog.Debug().Msgf("topic publish %v", resp)
	return nil
}
func (s *AppManager) List(query core.Query) {
	req := core.DataRequest{Opt: core.QUERY_DATA_REQUEST, Criteria: query}
	req.Async = query.QCc()
	s.Cluster().Request(req)
}

func (s *AppManager) Subscribe(topic string, listener core.EventListener) error {
	dsp := protocol.NewPostofficeServiceClient(s.rpc)
	resp, err := dsp.Subscribe(context.Background(), &protocol.Topic{NodeId: s.NodeId(), Tag: s.Context(), Name: topic})
	if err != nil {
		return err
	}
	core.AppLog.Debug().Msgf("topic registered %v", resp)
	return nil
}

func (s *AppManager) Unsubscribe(topic string) error {
	dsp := protocol.NewPostofficeServiceClient(s.rpc)
	resp, err := dsp.Unsubscribe(context.Background(), &protocol.Topic{Tag: s.Context(), Name: topic})
	if err != nil {
		return err
	}
	core.AppLog.Debug().Msgf("topic unregistered %v", resp)
	return nil
}

func (c *AppManager) receive() {
	dsp := protocol.NewPostofficeServiceClient(c.rpc)
	stream, err := dsp.Receive(context.Background(), &protocol.Topic{NodeId: c.NodeId(), Tag: c.Context()})
	if err != nil {
		core.AppLog.Warn().Msgf("rpc connection error %s", err.Error())
		return
	}
	for c.running {
		resp, err := stream.Recv()
		if err == io.EOF {
			core.AppLog.Debug().Msgf("eof %s", err.Error())
			break
		}
		if err != nil {
			core.AppLog.Warn().Msgf("streaming error %s", err.Error())
			break
		}
		data := resp.Event
		e := event.CreateEvent(data.Header.ClassId)
		err = core.Import(e, data.Key, data.Value, 200)
		if err != nil {
			core.AppLog.Debug().Msgf("event parse error %s", err.Error())
		} else {
			c.OnEvent(e)
		}
	}
	core.AppLog.Warn().Msg("rpc closed from remote")
}
