package main

import (
	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/mj"
)

// event util
func NewMahjongErrorEvent(systemId int64, tableId int64, code int, msg string) MahjongErrorEvent {
	me := MahjongErrorEvent{TableId: tableId, Code: code, Message: msg}
	me.SystemId = systemId
	return me
}

func NewMahjongTurnEvent(systemId int64, oid int64, cmd int, t OnTurn) MahjongTurnEvent {
	mt := MahjongTurnEvent{}
	mt.SystemId = systemId
	mt.OnOId(oid)
	mt.N = MahjongPlayTurn{Cmd: cmd, CountDown: TURN_TICKER_SECONDS}
	mt.T = t
	return mt
}

func NewMahjongHandEvent(systemId int64, h mj.Hand, k []int) MahjongHandEvent {
	mh := MahjongHandEvent{H: h, K: k}
	mh.SystemId = systemId
	return mh
}
func NewMahjongEvent(systemId int64, t MahjongPlayToken) MahjongEvent {
	me := MahjongEvent{Token: t}
	me.SystemId = systemId
	return me
}
func NewMahjongSitEvent(systemId int64,tableId int64,seat int) MahjongSitEvent{
	me := MahjongSitEvent{TableId: tableId,Seat: int32(seat)}
	me.SystemId = systemId
	return me
}

func main() {
	bootstrap.AppBootstrap(&MahjongService{})
}
