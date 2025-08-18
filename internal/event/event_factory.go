package event

const (
	STAT_CID         int = 1
	LOGIN_CID        int = 2
	MESSAGE_CID      int = 3
	TOURNAMENT_CID   int = 4
	SUBSCRIPTION_CID int = 5

	LOGIN_ETAG        string = "lgn:"
	MESSAGE_ETAG      string = "msg:"
	TOURNAMENT_ETAG   string = "tmt:"
	SUBSCRIPTION_ETAG string = "sub:"

	STAT_ETAG  string = "stat:"
	

	STAT_TOTAL string = "total"
)

func CreateEvent(cid int, listner EventListener) Event {
	switch cid {
	case STAT_CID:
		stat := StatEvent{}
		stat.Callback = listner
		return &stat
	case LOGIN_CID:
		login := LoginEvent{}
		login.Callback = listner
		return &login
	case MESSAGE_CID:
		message := MessageEvent{}
		message.Callback = listner
		return &message
	case TOURNAMENT_CID:
		tournament := TournamentEvent{}
		tournament.Callback = listner
		return &tournament
	case SUBSCRIPTION_CID:
		subscription := SubscriptionEvent{}
		subscription.Callback = listner
		return &subscription
	default:
		return nil
	}
}
