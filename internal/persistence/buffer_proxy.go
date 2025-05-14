package persistence

import (
	"bytes"
	"encoding/binary"

	buffer "github.com/0xc0d/encoding/bytebuffer"
)

type BufferProxy struct {
	data *buffer.ByteBuffer
}

func (s *BufferProxy) NewProxy(size int) {
	s.data = buffer.NewByteBuffer(size)
	s.data.SetOrder(binary.BigEndian)
}

func (s *BufferProxy) WriteBool(data bool) {
	if data {
		s.data.Put(byte(1))
	} else {
		s.data.Put(byte(0))
	}
}

func (s *BufferProxy) WriteComplex64(data complex64) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, data)
	s.data.PutBytes(buf.Bytes(), 0, buf.Len())
}

func (s *BufferProxy) WriteComplex128(data complex128) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, data)
	s.data.PutBytes(buf.Bytes(), 0, buf.Len())
}

func (s *BufferProxy) WriteFloat64(data float64) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, data)
	s.data.PutBytes(buf.Bytes(), 0, buf.Len())
}

func (s *BufferProxy) WriteFloat32(data float32) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, data)
	s.data.PutBytes(buf.Bytes(), 0, buf.Len())
}

func (s *BufferProxy) WriteInt64(data int64) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, data)
	s.data.PutBytes(buf.Bytes(), 0, buf.Len())
}

func (s *BufferProxy) WriteInt32(data int32) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, data)
	s.data.PutBytes(buf.Bytes(), 0, buf.Len())
}

func (s *BufferProxy) WriteInt16(data int16) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, data)
	s.data.PutBytes(buf.Bytes(), 0, buf.Len())
}

func (s *BufferProxy) WriteInt8(data int8) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, data)
	s.data.PutBytes(buf.Bytes(), 0, buf.Len())
}

func (s *BufferProxy) WriteString(data string) {
	slen := len(data)
	s.WriteInt32(int32(slen))
	s.data.PutBytes([]byte(data), 0, slen)
}

func (s *BufferProxy) ReadInt32() int32 {
	buf := make([]byte, 4)
	s.data.GetBytes(buf, 0, 4)
	var v int32
	binary.Read(bytes.NewBuffer(buf), binary.BigEndian, &v)
	return v
}

func (s *BufferProxy) ReadInt64() int64 {
	buf := make([]byte, 8)
	s.data.GetBytes(buf, 0, 8)
	var v int64
	binary.Read(bytes.NewBuffer(buf), binary.BigEndian, &v)
	return v
}

func (s *BufferProxy) ReadFloat32() float32 {
	buf := make([]byte, 4)
	s.data.GetBytes(buf, 0, 4)
	var v float32
	binary.Read(bytes.NewBuffer(buf), binary.BigEndian, &v)
	return v
}

func (s *BufferProxy) ReadFloat64() float64 {
	buf := make([]byte, 8)
	s.data.GetBytes(buf, 0, 8)
	var v float64
	binary.Read(bytes.NewBuffer(buf), binary.BigEndian, &v)
	return v
}

func (s *BufferProxy) ReadInt16() int16 {
	buf := make([]byte, 2)
	s.data.GetBytes(buf, 0, 2)
	var v int16
	binary.Read(bytes.NewBuffer(buf), binary.BigEndian, &v)
	return v
}

func (s *BufferProxy) ReadInt8() int8 {
	buf := make([]byte, 1)
	s.data.GetBytes(buf, 0, 1)
	var v int8
	binary.Read(bytes.NewBuffer(buf), binary.BigEndian, &v)
	return v
}

func (s *BufferProxy) ReadComplex64() complex64 {
	buf := make([]byte, 8)
	s.data.GetBytes(buf, 0, 8)
	var v complex64
	binary.Read(bytes.NewBuffer(buf), binary.BigEndian, &v)
	return v
}

func (s *BufferProxy) ReadComplex128() complex128 {
	buf := make([]byte, 16)
	s.data.GetBytes(buf, 0, 16)
	var v complex128
	binary.Read(bytes.NewBuffer(buf), binary.BigEndian, &v)
	return v
}

func (s *BufferProxy) ReadString() string {
	len := s.ReadInt32()
	buf := make([]byte, len)
	s.data.GetBytes(buf, 0, int(len))
	return string(buf)
}

func (s *BufferProxy) ReadBool() bool {
	b, _ := s.data.Get()
	return int(b) == 1
}
func (s *BufferProxy) Read() ([]byte, error) {
	len := s.data.Remaining()
	k := make([]byte, len)
	err := s.data.GetBytes(k, 0, len)
	if err != nil {
		return k, err
	}
	return k, nil
}

func (s *BufferProxy) Write(data []byte) {
	s.data.PutBytes(data, 0, len(data))
}

func (s *BufferProxy) Remaining() int {
	return s.data.Remaining()
}

func (s *BufferProxy) Flip() error {
	return s.data.Flip()
}
