package clustering

import (
	context "context"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
)

func (c *DataServiceProvider) Run(ctx context.Context, in *protocol.Task) (*protocol.Response, error) {
	//msg := make(chan *protocol.Response, 1)
	//defer close(msg)
	//setData := SetData{Opt: in.Opt, Data: in.Data, Resp: msg}
	//c.DSet <- setData
	//resp := <-msg
	core.AppLog.Debug().Msg("running task on target node")
	c.DMessager <- &protocol.Mail{Task: in, Opt: 1}
	return &protocol.Response{Successful: true, Message: "run task"}, nil
}
