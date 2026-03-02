package main

import (
	"fmt"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
	hash "github.com/spaolacci/murmur3"
)

const (
	COMPOSIT_KEY_MAX int = 500
)

type SetData struct {
	*protocol.Data
	Opt int
	Msg chan SetRes
}

// HELPER METHODS
func (s *SetData) IndexKey() KeyIndex {
	ki := KeyIndex{Prefix: hash.Sum32(s.Key), Header: s.Header, Key: s.Key}
	return ki
}

func (s *SetData) DataKey() ([]byte, error) {
	ksz := len(s.Key)
	if ksz+20 > COMPOSIT_KEY_MAX {
		return []byte{}, fmt.Errorf("Key size overflow %d", ksz)
	}
	buffer := core.NewBuffer(COMPOSIT_KEY_MAX)
	buffer.WriteInt32(s.Header.FactoryId)
	buffer.WriteInt32(s.Header.ClassId)
	buffer.WriteInt32(int32(ksz))
	buffer.Write(s.Key)
	buffer.WriteInt64(s.Header.Revision)
	buffer.Flip()
	return buffer.Read(0)
}
