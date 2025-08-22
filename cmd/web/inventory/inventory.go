package main

import "gameclustering.com/internal/bootstrap"

type Inventory struct {
	Id           int32  `json:"Id"`
	SystemId     int64  `json:"SystemId,string"`
	TypeId       string `json:"string"`
	Amount       int32  `json:"Amount"`
	Rechargeable bool   `json:"Rechargeable"`
}

type InventoryItem struct {
	Id          int32 `json:"Id"`
	InventoryId int32 `json:"InventoryId"`
	ItemId      int64 `json:"ItemId,string"`
}

func main() {
	bootstrap.AppBootstrap(&InventoryService{})
}
