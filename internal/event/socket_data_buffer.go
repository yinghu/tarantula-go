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

func (s *SocketBuffer) WriteInt32(data int32) error {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, data)
	if err != nil {
		return err
	}
	n, err := s.Socket.Write(buf.Bytes())
	if err != nil {
		return err
	}
	if n != 4 {
		return errors.New("less than 4 bytes")
	}
	return nil
}

func (s *SocketBuffer) ReadInt() (int, error) {
	n, err := s.Socket.Read(s.Buffer[:4])
	if err != nil {
		return 0, err
	}
	if n != 4 {
		return 0, errors.New("less than 4 bytes")
	}
	buf := bytes.NewBuffer(s.Buffer[:4])
	var v int
	if binary.Read(buf, binary.BigEndian, &v) != nil {
		return 0, errors.New("wrong data convert")
	}
	return v, nil
}

func (s *SocketBuffer) WriteInt(data int) error {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, data)
	if err != nil {
		return err
	}
	n, err := s.Socket.Write(buf.Bytes())
	if err != nil {
		return err
	}
	if n != 4 {
		return errors.New("less than 4 bytes")
	}
	return nil
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

func (s *SocketBuffer) WriteString(data string) error {
	sz := len(data)
	err := s.WriteInt32(int32(sz))
	if err != nil {
		return err
	}
	n, err := s.Socket.Write([]byte(data))
	if err != nil {
		return err
	}
	if n != int(sz) {
		msg := fmt.Sprintf("size not matched %d : %d", n, sz)
		return errors.New(msg)
	}
	return nil
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
