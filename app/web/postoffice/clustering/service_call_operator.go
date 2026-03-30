package clustering

import (
	"time"

	"gameclustering.com/internal/core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	SUCCESS        int = 0
	NO_TCP_CONNECT int = 1

	TASK_OPT_REQUEST int = 0
	TASK_OPT_CLOSE   int = 2
)

type Execute func(tcp *grpc.ClientConn, opt int) error

type Task struct {
	opt     int
	target  string
	execute Execute
}

type ServiceCallOperator struct {
	RTask      <-chan Task
	localConns map[string]*grpc.ClientConn
}

func (s *ServiceCallOperator) RunTask() {
	for task := range s.RTask {
		if task.opt == TASK_OPT_CLOSE {
			for _, c := range s.localConns {
				c.Close()
			}
			clear(s.localConns)
			return
		}
		c, ok := s.localConns[task.target]
		if ok {
			if err := task.execute(c, SUCCESS); err != nil {
				core.AppLog.Warn().Msgf("remove connecting from error %s", err.Error())
				delete(s.localConns, task.target)
			}
			continue
		}
		for r := range RETRY_MAX {
			c, err := grpc.NewClient(task.target, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err == nil {
				s.localConns[task.target] = c
				break
			}
			time.Sleep(1 * time.Second)
			core.AppLog.Warn().Msgf("retry to connect to %s retries : %d", task.target, r)
		}
		if c != nil {
			if err := task.execute(c, SUCCESS); err != nil {
				core.AppLog.Warn().Msgf("remove connecting from error %s", err.Error())
				delete(s.localConns, task.target)
			}
		} else {
			core.AppLog.Warn().Msgf("no connection from remote %s", task.target)
			task.execute(c, NO_TCP_CONNECT)
		}
	}
}
