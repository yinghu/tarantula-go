package item

type ItemService interface {
	Save(c Configuration) error
	Load(c Configuration) error
}
