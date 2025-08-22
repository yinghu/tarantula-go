package event

const (
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
		stat := StatEvent{}
		return &stat
	case LOGIN_CID:
		login := LoginEvent{}

		return &login
	case MESSAGE_CID:
		message := MessageEvent{}

		return &message
	case TOURNAMENT_CID:
		tournament := TournamentEvent{}

		return &tournament
	case SUBSCRIPTION_CID:
		subscription := SubscriptionEvent{}

		return &subscription
	case TOURNAMENT_JOIN_CID:
		join := TournamentJoinIndex{}
		return &join
	default:
		return nil
	}
}
