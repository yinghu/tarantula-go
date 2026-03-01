package main

import (
	"fmt"

	"gameclustering.com/internal/core"
	"gameclustering.com/tarantula/protocol"
	hash "github.com/spaolacci/murmur3"
)

type GetData struct {
	*protocol.Request
}

func (s *GetData) KeyIndex() KeyIndex {
	ki := KeyIndex{Prefix: hash.Sum32(s.Key), Header: s.Header, Key: s.Key}
	return ki
}
func (g *GetData) K() ([]byte, error) {
	ksz := len(g.Key)
	if ksz+20 > COMPOSIT_KEY_MAX {
		return []byte{}, fmt.Errorf("Key size overflow %d", ksz)
	}
	buffer := core.NewBuffer(COMPOSIT_KEY_MAX)
	buffer.WriteInt32(g.Header.FactoryId)
	buffer.WriteInt32(g.Header.ClassId)
	buffer.WriteInt32(int32(ksz))
	buffer.Write(g.Key)
	buffer.WriteInt64(g.Header.Revision)
	buffer.Flip()
	return buffer.Read(0)
}
