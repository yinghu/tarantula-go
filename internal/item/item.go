package item

import (
	"gameclustering.com/internal/core"
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

type ItemLoader interface {
	Load(cid int64) (Configuration, error)
}

type ItemService interface {
	core.SetUp
	SaveEnum(c Enum) error
	LoadEnum(cname string) (Enum, error)
	LoadEnums() ([]Enum, error)

	SaveCategory(c Category) error
	LoadCategory(cname string) (Category, error)
	LoadCategoryWithId(cid int64) (Category, error)
	LoadCategories(scopeEnd int32, targetScope string) []Category

	Save(c Configuration) error
	LoadWithName(cname string, limit int) ([]Configuration, error)
	LoadWithId(cid int64) (Configuration, error)

	DeleteWithName(cname string) error
	DeleteWithId(cid int64) error

	//ValidateEnum(c Enum) error
	Loader() ItemLoader
}

type KVUpdate struct {
	Key      string `json:"Key"`
	Value    string `json:"value"`
	core.Opt `json:"Opt"`
}

type ItemListener interface {
	OnUpdated(kv KVUpdate)
}
