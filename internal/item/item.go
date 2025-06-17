package item

type ItemService interface {
	Save(c Configuration) error
	LoadWithName(cname string) (Configuration, error)
	LoadWithId(cid int32) (Configuration, error)
}
