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

func (s *SocketBuffer) ReadInt8() (int8, error) {
	n, err := s.Socket.Read(s.Buffer[:1])
	if err != nil {
		return 0, err
	}
	if n != 1 {
		return 0, errors.New("less than 1 bytes")
	}
	buf := bytes.NewBuffer(s.Buffer[:4])
	var v int8
	if binary.Read(buf, binary.BigEndian, &v) != nil {
		return 0, errors.New("wrong data convert")
	}
	return v, nil
}

func (s *SocketBuffer) WriteInt8(data int8) error {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, data)
	if err != nil {
		return err
	}
	n, err := s.Socket.Write(buf.Bytes())
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("less than 1 bytes")
	}
	return nil
}

func (s *SocketBuffer) ReadInt16() (int16, error) {
	n, err := s.Socket.Read(s.Buffer[:2])
	if err != nil {
		return 0, err
	}
	if n != 2 {
		return 0, errors.New("less than 2 bytes")
	}
	buf := bytes.NewBuffer(s.Buffer[:2])
	var v int16
	if binary.Read(buf, binary.BigEndian, &v) != nil {
		return 0, errors.New("wrong data convert")
	}
	return v, nil
}

func (s *SocketBuffer) WriteInt16(data int16) error {
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
		return errors.New("less than 2 bytes")
	}
	return nil
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

func (s *SocketBuffer) ReadInt64() (int64, error) {
	n, err := s.Socket.Read(s.Buffer[:8])
	if err != nil {
		return 0, err
	}
	if n != 8 {
		return 0, errors.New("less than 8 bytes")
	}
	buf := bytes.NewBuffer(s.Buffer[:8])
	var v int64
	if binary.Read(buf, binary.BigEndian, &v) != nil {
		return 0, errors.New("wrong data convert")
	}
	return v, nil
}

func (s *SocketBuffer) WriteInt64(data int64) error {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, data)
	if err != nil {
		return err
	}
	n, err := s.Socket.Write(buf.Bytes())
	if err != nil {
		return err
	}
	if n != 8 {
		return errors.New("less than 8 bytes")
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
	if n != sz {
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
	if n != size {
		msg := fmt.Sprintf("size not matched %d :%d", n, size)
		return []byte{0}, errors.New(msg)
	}
	return s.Buffer[:n], nil
}

func (s *SocketBuffer) Write(data []byte) error {
	sz := len(data)
	n, err := s.Socket.Write(data)
	if err != nil {
		return err
	}
	if sz != n {
		msg := fmt.Sprintf("size not matched %d :%d", n, sz)
		return errors.New(msg)
	}
	return nil
}
