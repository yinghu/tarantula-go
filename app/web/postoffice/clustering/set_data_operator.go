package clustering

import (
	"fmt"

	"gameclustering.com/internal/core"
)

const (


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
		if sd.Opt == core.SET_OPT_RECOVER {
			break
		} else if sd.Opt == core.SET_OPT_CLOSE {
			core.AppLog.Debug().Msgf("closing set data operator %d", num)
			return
		}
		var err error
		switch sd.Opt {
		case core.CREATE_DATA_REQUEST:
			err = m.create(sd)
		case core.UPDATE_DATA_REQUEST:
			err = m.update(sd)
		case core.DELETE_DATA_REQUEST:
			err = m.delete(sd)
		case core.RESET_DATA_REQUEST:
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
