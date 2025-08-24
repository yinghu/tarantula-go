package core

type JsonRequester interface {
	PostJsonSync(url string, payload any) OnSession
}
