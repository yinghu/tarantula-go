package main

import (
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/protocol"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type LogForwarder struct {
	App *PostofficeService
}

func (s *LogForwarder) Forward(level zerolog.Level, log []byte) {
	lf := event.LogEventFactory{}
	e := protocol.LogEvent{}
	err := protojson.Unmarshal(log, &e)
	if err != nil {
		e.Level = "error"
		e.Message = err.Error()
		e.Time = timestamppb.Now()
		e.Source = "postoffice:64"
	}
	id, _ := s.App.Sequence().Id()
	t, _ := lf.FromLogEvent(&e)
	t.NodeId = s.App.NodeId()
	t.Tag = s.App.Context()
	t.Event.Id = uint64(id)
	s.App.Forward(t)
}
