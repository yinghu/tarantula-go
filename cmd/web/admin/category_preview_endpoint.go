package main

import (
	//"fmt"
	"fmt"
	"net/http"
	"strconv"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/item"
	"gameclustering.com/internal/util"
)

type CategoryPreview struct {
	Cat item.Category        `json:"category"`
	Ins []item.Configuration `json:"list"`
}

type CategoryPreviewer struct {
	*AdminService
}

func (s *CategoryPreviewer) AccessControl() int32 {
	return bootstrap.ADMIN_ACCESS_CONTROL
}
func (s *CategoryPreviewer) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	cid, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	cat, err := s.AdminService.ItemService().LoadCategoryWithId(cid)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	ins, err := s.AdminService.ItemService().LoadWithName(cat.Name, 10)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	for i := range ins {
		item.ItemView(ins[i], func(prop string, c item.Configuration) {
			fmt.Printf("Config %s , %s\n", prop, c.Category)
		})
	}
	w.Write(util.ToJson(CategoryPreview{Cat: cat, Ins: ins}))

}
