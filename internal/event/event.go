package event
import(
	"gameclustering.com/internal/persistence"
)

type Event interface {
	OnTopic() bool
	OnChan() chan []byte
	persistence.Persistentable
}
