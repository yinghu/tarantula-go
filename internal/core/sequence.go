package core

type Sequence interface {
	Id() (int64, error)
}
