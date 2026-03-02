package clustering

import (
	"fmt"

	"gameclustering.com/internal/core"
)

const (
	SET_OPT_RECOVER int = 1
	SET_OPT_CLOSE   int = 2
	SET_OPT_CREATE  int = 3
	SET_OPT_UPDATE  int = 4
	SET_OPT_DELETE  int = 5
	SET_OPT_RESET   int = 6

	SET_OPERATOR_NUM int = 8
)

type SetRes struct {
	Suc  bool
	Code int
	Err  error
}

func (m *DataServiceProvider) runSetData(num int) {
start:
	core.AppLog.Debug().Msgf("starting set operator %d", num)
	m.DWait.Wait()
	for sd := range m.DSet {
		if sd.Opt == SET_OPT_RECOVER {
			break
		} else if sd.Opt == SET_OPT_CLOSE {
			core.AppLog.Debug().Msgf("closing set data operator %d", num)
			return
		}
		var err error
		switch sd.Opt {
		case SET_OPT_CREATE:
			err = m.create(sd)
		case SET_OPT_UPDATE:
			err = m.update(sd)
		case SET_OPT_DELETE:
			err = m.delete(sd)
		case SET_OPT_RESET:
			err = m.reset(sd)
		default:
			err = fmt.Errorf("opt not supported %d", sd.Opt)
		}

		if err != nil {
			sd.Msg <- SetRes{Err: err}

		} else {
			sd.Msg <- SetRes{Suc: true}
		}
	}
	core.AppLog.Debug().Msgf("running recovery on operator %d", num)
	sync := <-m.DPull
	core.AppLog.Warn().Msgf("pulling data from %s", sync.Remote)
	for _, h := range sync.Hashs {
		core.AppLog.Warn().Msgf("recovering data from hash %d", h)
		//m.ClientPull()
	}
	m.DWait.Done()
	goto start
}
