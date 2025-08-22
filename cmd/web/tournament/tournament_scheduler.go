package main

import (
	"encoding/json"
	"time"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/item"
)

type Schedule struct {
	StartTime time.Time `json:"StartTime"`
	CloseTime time.Time `json:"CloseTime"`
	EndTime   time.Time `json:"EndTime"`
}

func (a *TournamentService) scheduleInstance(conf item.Configuration) {
	header, err := json.Marshal(conf.Header)
	if err != nil {
		core.AppLog.Printf("no header data %s\n", err.Error())
		return
	}
	ins := InstanceSchedule{TournamentService: a}
	ins.TournamentId = conf.Id
	err = json.Unmarshal(header, &ins)
	if err != nil {
		core.AppLog.Printf("no header data %s\n", err.Error())
		return
	}
	core.AppLog.Printf("Schedule :%d %v\n", ins.TournamentId, ins)
	sc, ok := conf.Reference["Schedule"].([]item.Configuration)
	if !ok {
		core.AppLog.Printf("no schedule data\n")
		return
	}
	jsc, err := json.Marshal(sc[0].Header)
	if err != nil {
		core.AppLog.Printf("no schedule data\n")
		return
	}
	err = json.Unmarshal(jsc, &ins.Schedule)
	if err != nil {
		core.AppLog.Printf("no schedule data %s\n", err.Error())
		return
	}
	core.AppLog.Printf("Schedule :%d %v\n", ins.TournamentId, ins)
	a.tournaments[ins.TournamentId] = &ins
	err = a.updateInstanceSchedule(ins)
	if err != nil {
		core.AppLog.Printf("sql err :%s\n", err.Error())
	}
	info := event.MessageEvent{Title: "info", Message: "tournament registered", Source: a.Context(), DateTime: time.Now()}
	id, _ := a.Sequence().Id()
	info.OnOId(id)
	info.OnTopic("message")
	a.Send(&info)
}

func (a *TournamentService) scheduleSegment(conf item.Configuration) {
	header, err := json.Marshal(conf.Header)
	if err != nil {
		core.AppLog.Printf("no header data %s\n", err.Error())
		return
	}
	seg := SegmentSchedule{TournamentService: a}
	seg.TournamentId = conf.Id
	err = json.Unmarshal(header, &seg)
	if err != nil {
		core.AppLog.Printf("no header data %s\n", err.Error())
		return
	}
	sc, ok := conf.Reference["Schedule"].([]item.Configuration)
	if !ok {
		core.AppLog.Printf("no schedule data\n")
		return
	}
	jsc, err := json.Marshal(sc[0].Header)
	if err != nil {
		core.AppLog.Printf("no schedule data\n")
		return
	}
	err = json.Unmarshal(jsc, &seg.Schedule)
	if err != nil {
		core.AppLog.Printf("no schedule data %s\n", err.Error())
		return
	}

	core.AppLog.Printf("Schedule :%d %v\n", seg.TournamentId, conf.Reference)
	refs, ok := conf.Reference["SegmentList"].([]item.Configuration)
	if !ok {
		core.AppLog.Printf("no segement data\n")
		return
	}
	seg.Segments = make([]*Segment, 0)
	for i := range refs {
		cf := refs[i]
		header, err := json.Marshal(cf.Header)
		if err != nil {
			continue
		}
		sg := Segment{InstanceId: cf.Id}
		err = json.Unmarshal(header, &sg)
		if err != nil {
			continue
		}
		seg.Segments = append(seg.Segments, &sg)
		core.AppLog.Printf("segement data %d %v\n", sg.InstanceId, sg)
	}
	core.AppLog.Printf("SEG SCHEDULE %s\n", seg.Name)
	//seg.Start()
	a.tournaments[seg.TournamentId] = &seg

	err = a.updateSegmentSchedule(seg)
	if err != nil {
		core.AppLog.Printf("sql err :%s\n", err.Error())
	}
}

func (a *TournamentService) releaseTournament(id int64) {
	delete(a.tournaments, id)
	a.updateSchedule(id)
}
