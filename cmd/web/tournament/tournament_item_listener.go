package main

import (
	"encoding/json"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/item"
)

func (a *TournamentService) OnRegister(conf item.Configuration) {
	core.AppLog.Printf("item reigster %d %s\n", conf.Id, conf.Category)
	if conf.Category == "InstanceSchedule" {
		header, err := json.Marshal(conf.Header)
		if err != nil {
			core.AppLog.Printf("no header data %s\n", err.Error())
			return
		}
		ins := InstanceSchedule{}
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
		jsc, err := json.Marshal(sc)
		if err != nil {
			core.AppLog.Printf("no schedule data\n")
			return
		}
		err = json.Unmarshal(jsc, &ins.Schedule)
		if err != nil {
			core.AppLog.Printf("no schedule data\n")
			return
		}
		core.AppLog.Printf("Schedule :%d %v\n", ins.TournamentId, ins)
		return
	}
	if conf.Category == "SegmentSchedule" {
		header, err := json.Marshal(conf.Header)
		if err != nil {
			core.AppLog.Printf("no header data %s\n", err.Error())
			return
		}
		seg := SegementSchedule{}
		seg.TournamentId = conf.Id
		err = json.Unmarshal(header, &seg)
		if err != nil {
			core.AppLog.Printf("no header data %s\n", err.Error())
			return
		}
		core.AppLog.Printf("Schedule :%d %v\n", seg.TournamentId, conf.Reference)
		refs, ok := conf.Reference["SegmentList"].([]item.Configuration)
		if !ok {
			core.AppLog.Printf("no segement data\n")
			return
		}
		seg.Segments = make([]Segment, 0)
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
			seg.Segments = append(seg.Segments, sg)
			core.AppLog.Printf("segement data %d %v\n", sg.InstanceId, sg)
		}
		core.AppLog.Printf("SEG SCHEDULE %v\n", seg)
	}
}
func (a *TournamentService) OnRelease(conf item.Configuration) {
	core.AppLog.Printf("item release %d %s\n", conf.Id, conf.Category)
}
