package main

import (
	//"fmt"
	"fmt"
	"net/http"
	"strconv"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type CategoryPreview struct {
	Cat core.Category        `json:"category"`
	Ins []core.Configuration `json:"list"`
}

type CategoryPreviewer struct {
	*AdminService
}

func (s *CategoryPreviewer) AccessControl() int32 {
	return core.ADMIN_ACCESS_CONTROL
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
		core.ItemValidator(ins[i], func(prop string, c core.Configuration) {
			fmt.Printf("Config %s , %s\n", prop, c.Category)
		})
	}
	w.Write(util.ToJson(CategoryPreview{Cat: cat, Ins: ins}))

}
