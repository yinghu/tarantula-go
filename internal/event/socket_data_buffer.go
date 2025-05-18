package event

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"

	"gameclustering.com/internal/core"
)

type SocketBuffer struct {
	Socket net.Conn
	Buffer []byte
	core.DataBufferHook
}

func (s *SocketBuffer) ReadInt32() (int32, error) {
	n, err := s.Socket.Read(s.Buffer[:4])
	if err != nil {
		return 0, err
	}
	if n != 4 {
		return 0, errors.New("less than 4 bytes")
	}
	buf := bytes.NewBuffer(s.Buffer[:4])
	var v int32
	if binary.Read(buf, binary.BigEndian, &v) != nil {
		return 0, errors.New("wrong data convert")
	}
	return v, nil
}
func (s *SocketBuffer) ReadString() (string, error) {
	sz, err := s.ReadInt32()
	if err != nil {
		return "", err
	}
	n, err := s.Socket.Read(s.Buffer[:sz])
	if err != nil {
		return "", err
	}
	if n != int(sz) {
		msg := fmt.Sprintf("size not matched %d : %d", n, sz)
		return "", errors.New(msg)
	}
	return string(s.Buffer[:sz]), nil
}

func (s *SocketBuffer) Read(size int) ([]byte, error) {
	n, err := s.Socket.Read(s.Buffer[:size])
	if err != nil {
		return []byte{0}, err
	}
	if n != int(size) {
		msg := fmt.Sprintf("size not matched %d :%d", n, size)
		return []byte{0}, errors.New(msg)
	}
	return s.Buffer[:n], nil
}
