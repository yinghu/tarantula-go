package conf

import (
	"fmt"
	"testing"
)



func TestFromFile(t *testing.T) {
	env := Env{}
	err := env.Load("/etc/tarantula/presence-conf.json")
	if err!= nil{
		t.Errorf("Config error %s", err.Error())
	}
	fmt.Printf("Conf %v\n",env) 	
}