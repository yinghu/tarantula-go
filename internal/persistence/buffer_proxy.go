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

func (s *BufferProxy) WriteInt32(data int32) {
	buf := new(bytes.Buffer)
	binary.Write(buf,binary.BigEndian,data)
	s.data.PutBytes(buf.Bytes(),0,buf.Len())
}

func (s *BufferProxy) WriteString(data string) {
	slen := len(data)
	s.WriteInt32(int32(slen))
	s.data.PutBytes([]byte(data),0,slen)
}

func (s *BufferProxy) ReadInt32() int32 {
	buf := make([]byte,4)
	s.data.GetBytes(buf,0,4)
	var v int32
	binary.Read(bytes.NewBuffer(buf),binary.BigEndian,&v)
	return v
}

func (s *BufferProxy) ReadString() string {
	len := s.ReadInt32()
	buf := make([]byte,len)
	s.data.GetBytes(buf,0,int(len))
	return string(buf)
}

func (s *BufferProxy) Array() ([]byte, error) {
	len := s.data.Remaining()
	k := make([]byte, len)
	err := s.data.GetBytes(k, 0, len)
	if err != nil {
		return k, err
	}
	return k, nil
}

func (s *BufferProxy) Write(data []byte){
	s.data.PutBytes(data,0,len(data))	
}

func (s *BufferProxy) Remaining() int{
	return s.data.Remaining()
}

func (s *BufferProxy) Flip() error{
	return s.data.Flip()
}
