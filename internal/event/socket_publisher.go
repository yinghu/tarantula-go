package event

import "net"

type SocketPublisher struct {
	Proto      string
	Remote     string
	BufferSize int
}

func (s *SocketPublisher) Publish(e Event) error {
	conn, err := net.Dial(s.Proto, s.Remote)
	if err != nil {
		return err
	}
	defer conn.Close()
	buffer := SocketBuffer{Socket: conn, Buffer: make([]byte, s.BufferSize)}
	buffer.WriteInt(e.ClassId())
	return nil
}
