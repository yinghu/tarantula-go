package item

import (
	"time"

	"gameclustering.com/internal/core"
)

const (
	GRANTABLE_ITEM = "Commodity"
)

type Enum struct {
	Id     int64       `json:"Id,string"`
	Name   string      `json:"Name"`
	Values []EnumValue `json:"Values"`
}

type EnumValue struct {
	Name  string `json:"Name"`
	Value int32  `json:"Value"`
}

type Category struct {
	Id            int64      `json:"Id,string"`
	Scope         string     `json:"Scope"`
	ScopeSequence int32      `json:"ScopeSequence"`
	Name          string     `json:"Name"`
	Rechargeable  bool       `json:"Rechargeable"`
	Description   string     `json:"Description"`
	Properties    []Property `json:"Properties"`
}

type Property struct {
	Name      string `json:"Name"`
	Type      string `json:"Type"`
	Reference string `json:"Reference"`
	Nullable  bool   `json:"Nullable"`
}

type Configuration struct {
	Id          int64               `json:"ItemId,string"`
	Name        string              `json:"ConfigurationName"`
	Type        string              `json:"ConfigurationType"`
	TypeId      string              `json:"ConfigurationTypeId"`
	Category    string              `json:"ConfigurationCategory"`
	Version     string              `json:"ConfigurationVersion"`
	Header      map[string]any      `json:"header"`
	Application map[string][]string `json:"application"`
	Reference   map[string]any      `json:"reference"`
}

func (c *Configuration) Amount(cat Category) int32 {
	if !cat.Rechargeable {
		return 1
	}
	for i := range cat.Properties {
		if cat.Properties[i].Name == "Amount" {
			v, exists := c.Header["Amount"]
			if exists {
				am, ok := v.(float64)
				if ok {
					return int32(am)
				}
				return 0
			}
		}
	}
	return 0
}

type ConfigRegistration struct {
	Id         int32     `json:"Id,string"`
	ItemId     int64     `json:"ItemId,string"`
	App        string    `json:"App"`
	Env        string    `json:"Env"`
	Scheduling bool      `json:"Scheduling"`
	StartTime  time.Time `json:"StartTime"`
	CloseTime  time.Time `json:"CloseTime"`
	EndTime    time.Time `json:"EndTime"`
}

type OnInventory struct {
	SystemId    int64  `json:"SystemId,string"`
	TypeId      string `json:"TypeId"`
	ItemId      int64  `json:"ItemId,string"`
	Source      string `json:"Source"`
	Description string `json:"Description"`
}

type Inventory struct {
	Id           int32     `json:"Id"`
	SystemId     int64     `json:"SystemId,string"`
	TypeId       string    `json:"string"`
	Amount       int32     `json:"Amount"`
	Rechargeable bool      `json:"Rechargeable"`
	UpdateTime   time.Time `json:"UpdateTime"`
	ItemId       int64     `json:"-"`
}

type InventoryItem struct {
	Id          int32 `json:"Id"`
	InventoryId int32 `json:"InventoryId"`
	ItemId      int64 `json:"ItemId,string"`
}

type InventoryManager interface {
	Reload(kv KVUpdate) error
	Load(cid int64) (Configuration, error)
	LoadCategory(name string) (Category, error)
	Grant(inv OnInventory) error
	Validate(c Configuration, validator Validator)
	Stock(r OnInventory) ([]Inventory, error)
}

type ItemService interface {
	core.SetUp
	SaveEnum(c Enum) error
	LoadEnum(cname string) (Enum, error)
	LoadEnums() ([]Enum, error)
	DeleteEnumWithId(cid int64) error

	SaveCategory(c Category) error
	LoadCategory(cname string) (Category, error)
	LoadCategoryWithId(cid int64) (Category, error)
	DeleteCategoryWithId(cid int64) error
	LoadCategories(scopeEnd int32, targetScope string) []Category

	Save(c Configuration) error
	LoadWithName(cname string, limit int) ([]Configuration, error)
	LoadWithId(cid int64) (Configuration, error)
	DeleteWithId(cid int64) error
	Register(reg ConfigRegistration) error
	Check(reg ConfigRegistration) (ConfigRegistration, error)
	Release(regId int32) error
	SaveRegistration(reg ConfigRegistration) error
	DeleteRegistration(itemId int64, app string, env string) error
	InventoryManager() InventoryManager
}

type RepoUpdate struct {
	Source string `json:"Source"`
	Target string `json:"Target"`
	Admin  string `json:"Admin"`
}

type KVUpdate struct {
	Key      string `json:"Key"`
	Value    string `json:"value"`
	core.Opt `json:"Opt"`
}

type ItemListener interface {
	OnRegister(conf Configuration)
	OnRelease(conf Configuration)
}
