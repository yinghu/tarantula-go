package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gameclustering.com/internal/persistence"
	"gameclustering.com/internal/util"
)

type Chunk struct {
	Remaining bool
	Data      []byte
}

type Login struct {
	Name        string `json:"login"`
	Hash        string `json:"password"`
	ReferenceId int32  `json:"referenceId"`
	SystemId    int64
}

type Service struct {
	Sql     persistence.Postgresql
	Sfk     util.Snowflake
	Tkn     util.Jwt
	Ciph    util.Cipher
	Started bool
	NodeId uint64
	DatabaseURL string
}

func (s *Service) Start() error {
	s.Sfk = util.NewSnowflake(int64(s.NodeId),util.EpochMillisecondsFromMidnight(2020, 1, 1))
	s.Tkn = util.Jwt{Alg: "SHS256"}
	s.Tkn.HMac()
	ci := util.Cipher{Ksz: 32}
	er := ci.AesGcm()
	if er != nil {
		return er
	}
	s.Ciph = ci
	sql := persistence.Postgresql{Url:s.DatabaseURL}
	err := sql.Create()
	if err != nil {
		return err
	}
	s.Sql = sql
	s.Started = true
	fmt.Printf("Auth service started\n")
	return nil
}
func (s *Service) Shutdown() {
	s.Sql.Close()
	fmt.Printf("Auth service shut down\n")
}

func (s *Service) Register(login *Login) error {
	id, _ := s.Sfk.Id()
	login.SystemId = id
	hash, _ := util.Hash(login.Hash)
	login.Hash = hash
	s.Sql.Exec("INSERT INTO login (name,hash,system_id,reference_id) VALUES($1,$2,$3,$4)", login.Name, login.Hash, login.SystemId, login.ReferenceId)
	return nil
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	action := r.Header.Get("Tarantula-action")
	defer func() {
		r.Body.Close()
	}()
	w.WriteHeader(http.StatusOK)
	switch action {
	case "onRegister":
		var login Login
		json.NewDecoder(r.Body).Decode(&login)
		err := s.Register(&login)
		if err != nil {
			w.Write([]byte(err.Error()))
		} else {
			w.Write([]byte("success"))
		}
	case "onLogin":
		
	default:
		w.Write([]byte("not supported"))
	}
}
