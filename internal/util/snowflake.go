package util

import (
	"errors"
	"sync"
	"time"
)

const (
	BITS_OF_NODE_NUMBER = 10
	BITS_OF_SEQUENCE    = 12
	MAX_NODE_ID         = (1 << BITS_OF_NODE_NUMBER) - 1 //1023
	MAX_SEQUENCE        = (1 << BITS_OF_SEQUENCE) - 1    //4095
)

type Snowflake struct {
	NodeId     int64
	EpochStart int64

	LastTimestamp int64
	Sequence      int64

	Lock sync.Mutex
}

func (s *Snowflake) Id() (int64, error) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	currentTimestamp := time.Now().UnixMilli() - s.EpochStart
	if currentTimestamp < s.LastTimestamp {
		return -1, errors.New("wrong system clock setting")
	}
	if currentTimestamp == s.LastTimestamp {
		s.Sequence = (s.Sequence + 1) & MAX_SEQUENCE
		if s.Sequence == 0 {
			for currentTimestamp == s.LastTimestamp {
				currentTimestamp = time.Now().UnixMilli() - s.EpochStart
			}
		}
	} else {
		s.Sequence = 0
	}
	s.LastTimestamp = currentTimestamp
	id := currentTimestamp<<(BITS_OF_NODE_NUMBER+BITS_OF_SEQUENCE) | (s.NodeId << BITS_OF_SEQUENCE) | s.Sequence
	return id, nil
}

func (s *Snowflake) Parse(snowflakeId int64) (int64, int64, int64) {
	nodeIdMask := ((1 << BITS_OF_NODE_NUMBER) - 1) << BITS_OF_SEQUENCE
	sequenceMask := (1 << BITS_OF_SEQUENCE) - 1
	timestamp := (snowflakeId >> (BITS_OF_NODE_NUMBER + BITS_OF_SEQUENCE))
	nodeId := (snowflakeId & int64(nodeIdMask)) >> BITS_OF_SEQUENCE
	sequence := snowflakeId & int64(sequenceMask)
	return timestamp, nodeId, sequence
}
