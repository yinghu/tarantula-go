package auth

import (
	"encoding/json"
	"net/http"

	"gameclustering.com/internal/util"
)

type Register struct {
	Login       string `json:"login"`
	Password    string `json:"password"`
	ReferenceId int    `json:"referenceId"`
}

type payload struct {
	remaining bool
	data      []byte
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}


func AuthHandler(w http.ResponseWriter, r *http.Request) {
	util.EpochMillisecondsFromMidnight(2011,1,1)
	var reg Register
	err := json.NewDecoder(r.Body).Decode(&reg)
	if err != nil {
		resp, _ := json.Marshal(Response{Code: 500, Message: err.Error()})
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
		return
	}
	mq := make(chan payload, 3)
	defer func() {
		close(mq)
		r.Body.Close()
	}()
	go func(Register) {
		//fmt.Printf("%v", reg)
		mq <- payload{remaining: true, data: []byte(reg.Login)}
		mq <- payload{remaining: true, data: []byte(reg.Password)}
		mq <- payload{remaining: false, data: []byte("onRegister")}
	}(reg)
	w.Header().Set("tarantula-Name", "token")
	w.WriteHeader(http.StatusOK)
	for {
		m := <-mq
		if !m.remaining {
			w.Write(m.data)
			break
		}
		w.Write(m.data)
	}
}
