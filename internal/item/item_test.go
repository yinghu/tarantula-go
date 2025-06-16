package item

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestHeader(t *testing.T) {
	h := `{"ConfigurationName":"Item1","ConfigurationType":"I100","header":{"name":"HP","value":200},"application":{"SkuList":[1,2,4],"ItemList":[1,2,4]},"reference":[1,2,4]}`
	var c Configuration

	err := json.Unmarshal([]byte(h), &c)
	if err != nil {
		t.Errorf("Error %s\n", err.Error())
	}
	fmt.Printf("JSON %v %v %v %s %s\n", c.Header, c.Reference, c.Application,c.Name,c.Type)
}
