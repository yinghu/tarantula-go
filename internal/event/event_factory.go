package event

import "fmt"

const (
	STAT_CID            int = 1
	LOGIN_CID           int = 2
	MESSAGE_CID         int = 3
	TOURNAMENT_CID      int = 4
	SUBSCRIPTION_CID    int = 5
	TOURNAMENT_JOIN_CID int = 6

	LOGIN_ETAG        string = "lgn:"
	MESSAGE_ETAG      string = "msg:"
	TOURNAMENT_ETAG   string = "tmt:"
	SUBSCRIPTION_ETAG string = "sub:"

	STAT_ETAG string = "stat:"

	STAT_TOTAL string = "total"
)

func CreateEvent(cid int, listner EventListener) Event {
	if listner != nil {
		fmt.Printf("Listener is ok")
	}
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
	case TOURNAMENT_JOIN_CID:
		join := TournamentJoinIndex{}
		join.Callback = listner
		return &join
	default:
		return nil
	}
}
