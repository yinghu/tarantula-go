package event

import (
	"gameclustering.com/internal/core"
)

type TournamentEvent struct {
	Id           int64 `json:"id,string"`

	TournamentId int64 `json:"TournamentId,string"`
	InstanceId   int64 `json:"InstanceId,string"`
	SystemId     int64 `json:"SystemId,string"`
	Score        int64 `json:"Score,string"`
	LastUpdated  int64 `json:"LastUpdated,string"`
	EventObj     `json:"-"`
}

func (s *TournamentEvent) ClassId() int {
	return TOURNAMENT_CID
}

func (s *TournamentEvent) ETag() string {
	return TOURNAMENT_ETAG
}

func (s *TournamentEvent) WriteKey(buff core.DataBuffer) error {
	if err := buff.WriteString(s.ETag()); err != nil {
		return err
	}
	if err := buff.WriteInt64(s.Id); err != nil {
		return err
	}
	return nil
}

func (s *TournamentEvent) ReadKey(buff core.DataBuffer) error {
	_, err := buff.ReadString()
	if err != nil {
		return err
	}
	id, err := buff.ReadInt64()
	if err != nil {
		return err
	}
	s.Id = id
	return nil
}

func (s *TournamentEvent) Write(buff core.DataBuffer) error {
	err := buff.WriteInt64(s.TournamentId)
	if err != nil {
		return err
	}
	err = buff.WriteInt64(s.InstanceId)
	if err != nil {
		return err
	}
	err = buff.WriteInt64(s.SystemId)
	if err != nil {
		return err
	}
	err = buff.WriteInt64(s.Score)
	if err != nil {
		return err
	}
	err = buff.WriteInt64(s.LastUpdated)
	if err != nil {
		return err
	}
	return nil
}
func (s *TournamentEvent) Read(buff core.DataBuffer) error {
	tournamentId, err := buff.ReadInt64()
	if err != nil {
		return err
	}
	s.TournamentId = tournamentId
	intanceId, err := buff.ReadInt64()
	if err != nil {
		return err
	}
	s.InstanceId = intanceId
	systemId, err := buff.ReadInt64()
	if err != nil {
		return err
	}
	s.SystemId = systemId
	score, err := buff.ReadInt64()
	if err != nil {
		return err
	}
	s.Score = score
	lastUpdated, err := buff.ReadInt64()
	if err != nil {
		return err
	}
	s.LastUpdated = lastUpdated
	return nil
}


func (s *TournamentEvent) Outbound(buff core.DataBuffer) error {
	err := s.Write(buff)
	if err != nil {
		s.Callback.OnError(err)
		return err
	}
	return nil
}

func (s *TournamentEvent) Inbound(buff core.DataBuffer) error {
	err := s.Read(buff)
	if err != nil {
		s.Callback.OnError(err)
		return err
	}
	s.Callback.OnEvent(s)
	return nil
}
