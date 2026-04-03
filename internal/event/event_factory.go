package event

import "gameclustering.com/internal/core"

const (
	CLOSE_CID            uint32 = 0
	STAT_CID             uint32 = 1
	LOGIN_CID            uint32 = 2
	MESSAGE_CID          uint32 = 3
	TOURNAMENT_CID       uint32 = 4
	SUBSCRIPTION_CID     uint32 = 5
	TOURNAMENT_SCORE_CID uint32 = 6
	INVENTORY_CID        uint32 = 7
	JOIN_CID             uint32 = 8
	KICKOFF_CID          uint32 = 9
	REGISTER_CID         uint32 = 10

	LOGIN_ETAG        string = "lgn"
	MESSAGE_ETAG      string = "msg"
	TOURNAMENT_ETAG   string = "tmt"
	SUBSCRIPTION_ETAG string = "sub"
	INVENTORY_ETAG    string = "inv"
	JOIN_ETAG         string = "join"
	KICKOFF_ETAG      string = "koff"
	REGISTER_ETAG     string = "reg"

	STAT_ETAG string = "stat"

	STAT_TOTAL string = "total"
)

func CreateEvent(cid uint32) core.Event {
	switch cid {
	case STAT_CID:
		return &StatEvent{}
	case LOGIN_CID:
		return &LoginEvent{}
	case TOURNAMENT_CID:
		return &TournamentEvent{}
	case SUBSCRIPTION_CID:
		return &SubscriptionEvent{}
	case TOURNAMENT_SCORE_CID:
		return &TournamentScoreIndex{}
	case INVENTORY_CID:
		return &InventoryEvent{}
	case JOIN_CID:
		return &JoinEvent{}
	case KICKOFF_CID:
		return &KickoffEvent{}
	case REGISTER_CID:
		return &RegisterEvent{}
	default:
		return nil
	}
}
