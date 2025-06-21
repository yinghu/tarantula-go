package item

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestHeader(t *testing.T) {
	h := `{"ConfigurationName":"Item1","ConfigurationType":"I100","header":{"name":"HP","value":200},"application":{"SkuList":[1],"ItemList":[2]}}`
	var c Configuration

	err := json.Unmarshal([]byte(h), &c)
	if err != nil {
		t.Errorf("Parse Error %s\n", err.Error())
	}
	fmt.Printf("JSON %v %v %s %s\n", c.Header, c.Application, c.Name, c.Type)
	for k, v := range c.Application {
		fmt.Printf("%s %v\n", k, v)
		for i := range v {
			fmt.Printf("%d\n", v[i])
		}
	}
	r, err := json.Marshal(c)
	if err != nil {
		t.Errorf("Flat Error %s\n", err.Error())
	}
	fmt.Printf("JSON %s\n", string(r))
}

func TestApplication(t *testing.T) {
	h := `{"ConfigurationName":"Item1","ConfigurationType":"I100","header":{"name":"HP","value":200}}`
	var c Configuration

	err := json.Unmarshal([]byte(h), &c)
	if err != nil {
		t.Errorf("Parse Error %s\n", err.Error())
	}
	fmt.Printf("JSON %v %v %s %s\n", c.Header, c.Application, c.Name, c.Type)
	for k, v := range c.Application {
		fmt.Printf("%s %v\n", k, v)
		for i := range v {
			fmt.Printf("%d\n", v[i])
		}
	}
	r, err := json.Marshal(c)
	if err != nil {
		t.Errorf("Flat Error %s\n", err.Error())
	}
	fmt.Printf("JSON %s\n", string(r))
}

func TestCategory(t *testing.T) {
	h := `{"Scope":"commodity","Name":"GameAsset","Rechargeable":false,"Description":"game view"}`
	var c Category

	err := json.Unmarshal([]byte(h), &c)
	if err != nil {
		t.Errorf("Parse Error %s\n", err.Error())
	}
	fmt.Printf("JSON %v\n", c)
	for i := range c.Properties {
		fmt.Printf("%d %v\n", i, c.Properties[i])
	}

	r, err := json.Marshal(c)
	if err != nil {
		t.Errorf("Flat Error %s\n", err.Error())
	}
	fmt.Printf("JSON %s\n", string(r))

	props := make([]Property, 0)
	props = append(props, Property{Name: "abc1"})
	fmt.Printf("Len : %d\n", len(props))
	props = append(props, Property{Name: "abc2"})
	fmt.Printf("Len : %d\n", len(props))
	props = append(props, Property{Name: "abc3"})
	fmt.Printf("Len : %d\n", len(props))
	props = append(props, Property{Name: "abc4"})
	fmt.Printf("Len : %d\n", len(props))
	//var v any
	var v any = 9223372036854775807
	x, ok := v.(int)
	if ok {
		fmt.Printf("Data %d\n", x)
	}
	fmt.Printf("fm %d\n", x)

	var f any = 9223372036854775807.23
	y, ok := f.(float64)
	if ok {
		fmt.Printf("Data %v\n", y)
	}
	fmt.Printf("fm %v\n", y)
	var dt any = "2025-05-29T10:00"
	tm, err := time.Parse("2006-01-02T15:04", fmt.Sprintf("%v", dt))
	if err != nil {
		t.Errorf("Time parse err : %s", err.Error())
	}
	fmt.Printf("%v\n", tm)
}
