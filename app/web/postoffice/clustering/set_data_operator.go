package clustering

import (
	"fmt"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
)

const (
	SET_OPERATOR_NUM int = 8
)

func (m *DataServiceProvider) runSetData(num int) {
start:
	m.DWait.Wait()
	core.AppLog.Debug().Msgf("starting set operator %d", num)
	for sd := range m.DSet {
		if sd.Opt == core.SET_OPT_RECOVER {
			break
		} else if sd.Opt == core.SET_OPT_CLOSE {
			core.AppLog.Debug().Msgf("closing set data operator %d", num)
			return
		}
		switch sd.Opt {
		case core.CREATE_DATA_REQUEST:
			ki, err := m.create(sd)
			if err != nil {
				sd.Resp <- &protocol.Response{Successful: false, Message: err.Error()}
			} else {
				var data []*protocol.Data
				data = append(data, &protocol.Data{Header: &protocol.Header{Revision: ki.Header.Revision}})
				sd.Resp <- &protocol.Response{Successful: true, Data: &protocol.DataSet{List: data}}
			}
		case core.UPDATE_DATA_REQUEST:
			ki, err := m.update(sd)
			if err != nil {
				sd.Resp <- &protocol.Response{Successful: false, Message: err.Error()}
			} else {
				var data []*protocol.Data
				data = append(data, &protocol.Data{Header: &protocol.Header{Revision: ki.Header.Revision}})
				sd.Resp <- &protocol.Response{Successful: true, Data: &protocol.DataSet{List: data}}
			}
		case core.DELETE_DATA_REQUEST:
			//err = m.delete(sd)
		case core.RESET_DATA_REQUEST:
			//err = m.reset(sd)
		default:
			sd.Resp <- &protocol.Response{Successful: false, Message: fmt.Sprintf("opt not supported %d", sd.Opt)}
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
