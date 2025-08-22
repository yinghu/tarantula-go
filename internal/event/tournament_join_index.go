package event

import "gameclustering.com/internal/core"

const (
	T_JOIN_TAG string = "join:"
)

type TournamentJoinIndex struct {
	
	TournamentId int64 `json:"TournamentId,string"`
	InstanceId   int64 `json:"InstanceId,string"`
	SystemId     int64 `json:"SystemId,string"`
	JoinTime     int64 `json:"JoinTime,string"`
	EventObj     `json:"-"`
}

func (s *TournamentJoinIndex) ClassId() int {
	return TOURNAMENT_JOIN_CID
}

func (s *TournamentJoinIndex) ETag() string {
	return TOURNAMENT_ETAG
}

func (s *TournamentJoinIndex) Distributed() bool {
	return true
}

func (s *TournamentJoinIndex) WriteKey(buff core.DataBuffer) error {
	if err := buff.WriteString(s.ETag()); err != nil {
		return err
	}
	if err := buff.WriteString(T_JOIN_TAG); err != nil {
		return err
	}
	if err := buff.WriteInt64(s.SystemId); err != nil {
		return err
	}
	if err := buff.WriteInt64(s.OnOId()); err != nil {
		return err
	}
	return nil
}

func (s *TournamentJoinIndex) ReadKey(buff core.DataBuffer) error {
	_, err := buff.ReadString()
	if err != nil {
		return err
	}
	_, err = buff.ReadString()
	if err != nil {
		return err
	}
	systemId, err := buff.ReadInt64()
	if err != nil {
		return err
	}
	s.SystemId = systemId
	id, err := buff.ReadInt64()
	if err != nil {
		return err
	}
	s.OId(id)
	return nil
}

func (s *TournamentJoinIndex) Write(buff core.DataBuffer) error {
	if err := buff.WriteInt64(s.TournamentId); err != nil {
		return err
	}
	if err := buff.WriteInt64(s.InstanceId); err != nil {
		return err
	}
	if err := buff.WriteInt64(s.JoinTime); err != nil {
		return err
	}
	return nil
}
func (s *TournamentJoinIndex) Read(buff core.DataBuffer) error {
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
	joinTime, err := buff.ReadInt64()
	if err != nil {
		return err
	}
	s.JoinTime = joinTime
	return nil
}

func (s *TournamentJoinIndex) Outbound(buff core.DataBuffer) error {
	if err := s.WriteKey(buff); err != nil {
		s.Callback.OnError(s,err)
		return err
	}
	if err := s.Write(buff); err != nil {
		s.Callback.OnError(s,err)
		return err
	}
	return nil
}

func (s *TournamentJoinIndex) Inbound(buff core.DataBuffer) error {
	if err := s.ReadKey(buff); err != nil {
		s.Callback.OnError(s,err)
		return err
	}
	if err := s.Read(buff); err != nil {
		s.Callback.OnError(s,err)
		return err
	}
	s.Callback.OnEvent(s)
	return nil
}
