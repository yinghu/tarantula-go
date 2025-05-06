package auth

import (
	//"encoding/json"
	"fmt"
	"net/http"

	"gameclustering.com/internal/persistence"
	//"gameclustering.com/internal/util"
)

type Login struct{
	Name string
	Hash string
	SystemId int64
}

type Service struct {
	Sql persistence.Postgresql
	Started bool
}

func (s *Service) Start() error {
	sql := persistence.Postgresql{Url: "postgres://postgres:password@192.168.1.7:5432/tarantula_user"}
	err := sql.Create()
	if err != nil{
		return err
	}
	s.Started = true
	return nil
}
func (s *Service) Shutdown() {
	s.Sql.Close()
	fmt.Printf("Auth service shut down\n")
}

func (s *Service) Register(login Login) error{
	
	return nil
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("test"))
}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	//util.EpochMillisecondsFromMidnight(2011, 1, 1)
	//var reg Register
	//err := json.NewDecoder(r.Body).Decode(&reg)
	//if err != nil {
	//resp, _ := json.Marshal(Response{Code: 500, Message: err.Error()})
	//w.WriteHeader(http.StatusOK)
	//w.Write(resp)
	//return
	//}
	//mq := make(chan payload, 3)
	//defer func() {
	//close(mq)
	//r.Body.Close()
	//}()
	//go func(Register) {
	//fmt.Printf("%v", reg)
	//mq <- payload{remaining: true, data: []byte(reg.Login)}
	//mq <- payload{remaining: true, data: []byte(reg.Password)}
	//mq <- payload{remaining: false, data: []byte("onRegister")}
	//}(reg)
	//w.Header().Set("tarantula-Name", "token")
	//w.WriteHeader(http.StatusOK)
	//for {
	//	m := <-mq
	//	if !m.remaining {
	//		w.Write(m.data)
	//		break
	//	}
	//	w.Write(m.data)
	//}
}
