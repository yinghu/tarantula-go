package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Cnf struct {
	GroupName string `json:"GroupName"`
	NodeName  string `json:"NodeName"`
	NodeId    int64  `json:"NodeId"`
}

func main() {
	fmt.Printf("Profile %s\n", "start")
	conf, err := os.Open("/etc/tarantula/profile_conf.json")
	if err != nil {
		fmt.Printf("Load %s\n", err.Error())
	}
	defer conf.Close()
	var f = Cnf{}
	data, _ := io.ReadAll(conf)
	json.Unmarshal(data, &f)
	fmt.Printf("Load : %s,%s\n", f.GroupName, f.NodeName)
	log.Fatal(http.ListenAndServe(":8090", nil))
}
