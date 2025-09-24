package event

const (
	CLOSE_CID           int = 0
	STAT_CID            int = 1
	LOGIN_CID           int = 2
	MESSAGE_CID         int = 3
	TOURNAMENT_CID      int = 4
	SUBSCRIPTION_CID    int = 5
	TOURNAMENT_JOIN_CID int = 6
	INVENTORY_CID       int = 7

	LOGIN_ETAG        string = "lgn:"
	MESSAGE_ETAG      string = "msg:"
	TOURNAMENT_ETAG   string = "tmt:"
	SUBSCRIPTION_ETAG string = "sub:"
	INVENTORY_ETAG    string = "inv"

	STAT_ETAG string = "stat:"

	STAT_TOTAL string = "total"
)

func CreateEvent(cid int) Event {
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
	case TOURNAMENT_JOIN_CID:
		return &TournamentJoinIndex{}
	case INVENTORY_CID:
		return &InventoryEvent{}
	default:
		return nil
	}
}
