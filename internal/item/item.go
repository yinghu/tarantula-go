package item
import (
	"gameclustering.com/internal/core"
)
type Enum struct {
	Id     int32       `json:"Id"`
	Name   string      `json:"Name"`
	Values []EnumValue `json:"Values"`
}

type EnumValue struct {
	Name  string `json:"Name"`
	Value int32  `json:"Value"`
}

type Category struct {
	Id           int32      `json:"Id"`
	Scope        string     `json:"Scope"`
	Name         string     `json:"Name"`
	Rechargeable bool       `json:"Rechargeable"`
	Description  string     `json:"Description"`
	Properties   []Property `json:"Properties"`
}

type Property struct {
	Name         string `json:"Name"`
	Type         string `json:"Type"`
	Reference    string `json:"Reference"`
	Nullable     bool   `json:"Nullable"`
	Downloadable bool   `json:"Downloadable"`
}

type Configuration struct {
	Id          int32              `json:"ItemId"`
	Name        string             `json:"ConfigurationName"`
	Type        string             `json:"ConfigurationType"`
	TypeId      string             `json:"ConfigurationTypeId"`
	Category    string             `json:"ConfigurationCategory"`
	Version     string             `json:"ConfigurationVersion"`
	Header      map[string]any     `json:"header"`
	Application map[string][]int32 `json:"application"`
}

type ItemService interface {
	core.SetUp
	SaveEnum(c Enum) error
	LoadEnum(cname string) (Enum, error)
	SaveCategory(c Category) error
	LoadCategory(cname string) (Category, error)

	Save(c Configuration) error
	LoadWithName(cname string, limit int) ([]Configuration, error)
	LoadWithId(cid int32) (Configuration, error)

	DeleteWithName(cname string) error
	DeleteWithId(cid int32) error

	Validate(c Configuration) error
	ValidateCategory(c Category) error
	ValidateEnum(c Enum) error
}
