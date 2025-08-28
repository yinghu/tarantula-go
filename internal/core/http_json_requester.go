package core

type Chunk struct {
	Remaining bool
	Data      []byte
}
type JsonRequester interface {
	PostJsonSync(url string, payload any) OnSession
	PostJsonAsync(url string, payload any, ch chan Chunk)
}
