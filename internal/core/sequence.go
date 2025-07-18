package core

type Sequence interface {
	Id() (int64, error)
	Parse(snowflakeId int64) (int64, int64, int64)
}
