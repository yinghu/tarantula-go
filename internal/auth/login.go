package auth
import(
	"gameclustering.com/internal/event"
)
type Login struct {
	Name        string `json:"login"`
	Hash        string `json:"password"`
	ReferenceId int32  `json:"referenceId"`
	SystemId    int64

	event.EventObj
}
