package event

import "gameclustering.com/internal/core"

const (
	CLOSE_CID            int32 = 0
	STAT_CID             int32 = 1
	LOGIN_CID            int32 = 2
	MESSAGE_CID          int32 = 3
	TOURNAMENT_CID       int32 = 4
	SUBSCRIPTION_CID     int32 = 5
	TOURNAMENT_SCORE_CID int32 = 6
	INVENTORY_CID        int32 = 7
	JOIN_CID             int32 = 8
	KICKOFF_CID          int32 = 9
	REGISTER_CID         int32 = 10

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

func CreateEvent(cid int32) core.Event {
	switch cid {
	case STAT_CID:
		return &StatEvent{}
	case LOGIN_CID:
		return &LoginEvent{}
	case MESSAGE_CID:
		return &MessageEvent{}
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
