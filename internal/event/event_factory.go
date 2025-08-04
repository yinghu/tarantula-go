package event

const (
	LOGIN_CID      int = 1
	MESSAGE_CID    int = 2
	TOURNAMENT_CID int = 3
)

func CreateEvent(cid int, listner EventListener) Event {
	switch cid {
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
	default:
		return nil
	}
}
