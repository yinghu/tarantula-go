package event

import "gameclustering.com/internal/core"

const (
	T_SCORE_TAG string = "score"
)

type TournamentScoreIndex struct {
	TournamentId int64 `json:"TournamentId,string"`
	InstanceId   int64 `json:"InstanceId,string"`
	Score        int64 `json:"Score,string"`
	UpdateTime   int64 `json:"UpdateTime,string"`
	SystemId     int64 `json:"SystemId,string"`
	EventObj     `json:"-"`
}

func (s *TournamentScoreIndex) ClassId() int {
	return TOURNAMENT_SCORE_CID
}

func (s *TournamentScoreIndex) ETag() string {
	return TOURNAMENT_ETAG
}

func (s *TournamentScoreIndex) Distributed() bool {
	return true
}

func (s *TournamentScoreIndex) WriteKey(buff core.DataBuffer) error {
	if err := buff.WriteString(s.ETag()); err != nil {
		return err
	}
	if err := buff.WriteString(T_SCORE_TAG); err != nil {
		return err
	}
	if err := buff.WriteInt64(s.TournamentId); err != nil {
		return err
	}
	if err := buff.WriteInt64(s.InstanceId); err != nil {
		return err
	}
	if err := buff.WriteInt64(s.Score); err != nil {
		return err
	}
	if err := buff.WriteInt64(s.UpdateTime); err != nil {
		return err
	}
	if err := buff.WriteInt64(s.SystemId); err != nil {
		return err
	}
	return nil
}

func (s *TournamentScoreIndex) ReadKey(buff core.DataBuffer) error {
	_, err := buff.ReadString()
	if err != nil {
		return err
	}
	_, err = buff.ReadString()
	if err != nil {
		return err
	}
	tid, err := buff.ReadInt64()
	if err != nil {
		return err
	}
	s.TournamentId = tid
	iid, err := buff.ReadInt64()
	if err != nil {
		return err
	}
	s.InstanceId = iid
	sco, err := buff.ReadInt64()
	if err != nil {
		return err
	}
	s.Score = sco
	updateTime, err := buff.ReadInt64()
	if err != nil {
		return err
	}
	s.UpdateTime = updateTime
	sid, err := buff.ReadInt64()
	if err != nil {
		return err
	}
	s.SystemId = sid
	return nil
}

func (s *TournamentScoreIndex) Write(buff core.DataBuffer) error {
	if err := buff.WriteInt64(s.SystemId); err != nil {
		return err
	}
	return nil
}
func (s *TournamentScoreIndex) Read(buff core.DataBuffer) error {
	sid, err := buff.ReadInt64()
	if err != nil {
		return err
	}
	s.SystemId = sid
	return nil
}

func (s *TournamentScoreIndex) Outbound(buff core.DataBuffer) error {
	if err := s.WriteKey(buff); err != nil {
		s.Callback.OnError(s, err)
		return err
	}
	if err := s.Write(buff); err != nil {
		s.Callback.OnError(s, err)
		return err
	}
	return nil
}

func (s *TournamentScoreIndex) Inbound(buff core.DataBuffer) error {
	if err := s.ReadKey(buff); err != nil {
		s.Callback.OnError(s, err)
		return err
	}
	if err := s.Read(buff); err != nil {
		s.Callback.OnError(s, err)
		return err
	}
	s.Callback.OnEvent(s)
	return nil
}
