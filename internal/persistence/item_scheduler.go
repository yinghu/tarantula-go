package persistence

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/item"
	"gameclustering.com/internal/util"
)

func (db *ItemDB) Schedule(reg item.ConfigRegistration) error {
	db.Cls.Atomic(db.Cls.Group(), func(ctx core.Ctx) error {
		ctx.Put("register:"+reg.App, string(util.ToJson(reg)))
		return nil
	})
	return nil
}

func (db *ItemDB) Unschedule(reg item.ConfigRegistration) error {
	db.Cls.Atomic(db.Cls.Group(), func(ctx core.Ctx) error {
		ctx.Put("release:"+reg.App, string(util.ToJson(reg)))
		return nil
	})
	return nil
}
