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
		core.AppLog.Printf("Schedule :%d %v\n", seg.TournamentId, seg)
		refs, ok := conf.Reference["SegmentList"].([]item.Configuration)
		if !ok {
			core.AppLog.Printf("no segement data\n")
			return
		}
		for i := range refs {
			sg := refs[i]
			core.AppLog.Printf("segement data %d %v\n", sg.Id, sg)
		}
		//[]item.Configuration(conf.Reference["SegmentList"])
		//reference, err := json.Marshal(conf.Reference)

	}
}
func (a *TournamentService) OnRelease(conf item.Configuration) {
	core.AppLog.Printf("item release %d %s\n", conf.Id, conf.Category)
}
