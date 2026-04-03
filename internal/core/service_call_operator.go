package core

import (
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	SUCCESS        int = 0
	NO_TCP_CONNECT int = 1

	TASK_OPT_REQUEST int = 0
	TASK_OPT_CLOSE   int = 2

	TASK_RETRY_MAX int = 3
)

type Execute func(tcp *grpc.ClientConn, opt int) error

type Task struct {
	Opt     int
	Target  string
	Execute Execute
}

type ServiceCallOperator struct {
	RTask      <-chan Task
	LocalConns map[string]*grpc.ClientConn
}

func (s *ServiceCallOperator) RunTask() {
	for task := range s.RTask {
		if task.Opt == TASK_OPT_CLOSE {
			for _, c := range s.LocalConns {
				c.Close()
			}
			clear(s.LocalConns)
			return
		}
		c, ok := s.LocalConns[task.Target]
		if ok {
			if err := task.Execute(c, SUCCESS); err != nil {
				AppLog.Warn().Msgf("remove connecting from error %s", err.Error())
				delete(s.LocalConns, task.Target)
			}
			continue
		}
		opt := NO_TCP_CONNECT
		for r := range TASK_RETRY_MAX {
			AppLog.Debug().Msgf("connecting to %s", task.Target)
			cx, err := grpc.NewClient(task.Target, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err == nil {
				s.LocalConns[task.Target] = cx
				opt = TASK_OPT_REQUEST
				if err := task.Execute(cx, SUCCESS); err != nil {
					AppLog.Warn().Msgf("remove connecting from error %s", err.Error())
					delete(s.LocalConns, task.Target)
				}
				break
			}
			time.Sleep(1 * time.Second)
			AppLog.Warn().Msgf("retry to connect to %s retries : %d", task.Target, r)
		}
		if opt == NO_TCP_CONNECT {
			AppLog.Warn().Msgf("no connection from remote %s", task.Target)
			task.Execute(c, NO_TCP_CONNECT)
		}

	}
}
