package item

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestHeader(t *testing.T) {
	h := `{"ConfigurationName":"Item1","ConfigurationType":"I100","header":{"name":"HP","value":200},"application":{"SkuList":[1],"ItemList":[2]},"reference":[1,2]}`
	var c Configuration

	err := json.Unmarshal([]byte(h), &c)
	if err != nil {
		t.Errorf("Parse Error %s\n", err.Error())
	}
	fmt.Printf("JSON %v %v %v %s %s\n", c.Header, c.Reference, c.Application, c.Name, c.Type)
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
	h := `{"ConfigurationName":"Item1","ConfigurationType":"I100","header":{"name":"HP","value":200},"reference":[1,2]}`
	var c Configuration

	err := json.Unmarshal([]byte(h), &c)
	if err != nil {
		t.Errorf("Parse Error %s\n", err.Error())
	}
	fmt.Printf("JSON %v %v %v %s %s\n", c.Header, c.Reference, c.Application, c.Name, c.Type)
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
