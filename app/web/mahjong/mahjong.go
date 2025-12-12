package main

import (
	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/mj"
)

// event util
func NewMahjongErrorEvent(systemId int64, tableId int64, code int, msg string) MahjongErrorEvent {
	me := MahjongErrorEvent{Code: code, Message: msg}
	me.TableId = tableId
	me.SystemId = systemId
	return me
}

func NewMahjongTurnEvent(systemId int64, oid int64, cmd int, t OnTurn, p OnStop) MahjongTurnEvent {
	mt := MahjongTurnEvent{}
	mt.SystemId = systemId
	mt.OnOId(oid)
	if systemId == 0 {
		mt.N = MahjongPlayTurn{Cmd: cmd, CountDown: AUTO_TURN_TICKER_SECONDS}
	} else {
		mt.N = MahjongPlayTurn{Cmd: cmd, CountDown: PLAYER_TURN_TICKER_SECONDS}
	}
	mt.T = t
	mt.P = p
	return mt
}

func NewMahjongHandEvent(systemId int64, h mj.Hand, f, k []int) MahjongHandEvent {
	mh := MahjongHandEvent{H: h, F: f, K: k}
	mh.SystemId = systemId
	return mh
}
func NewMahjongEvent(systemId int64, t MahjongPlayToken) MahjongEvent {
	me := MahjongEvent{Token: t}
	me.SystemId = systemId
	return me
}
func NewMahjongSitEvent(systemId int64, tableId int64, seat int) MahjongSitEvent {
	me := MahjongSitEvent{Seat: int32(seat)}
	me.TableId = tableId
	me.SystemId = systemId
	return me
}
func NewMahjongKongEvent(systemId int64, knog []int) MahjongKongEvent {
	me := MahjongKongEvent{Knog: knog}
	me.SystemId = systemId
	return me
}

func NewMahjongClaimEvent(systemId int64, tableId int64, seat int, claimed bool, formed []mj.Meld) MahjongClaimEvent {
	me := MahjongClaimEvent{Seat: int32(seat), Claimed: claimed, Formed: formed}
	me.SystemId = systemId
	me.TableId = tableId
	return me
}

func NewMahjongDiscardEvent(seat int, dropSeat int, drop mj.Tile, opts []mj.Meld) MahjongDiscardEvent {
	me := MahjongDiscardEvent{Seat: seat, DropSeat: dropSeat, Drop: drop, Opts: opts}
	return me
}

func main() {
	bootstrap.AppBootstrap(&MahjongService{})
}
