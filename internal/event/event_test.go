package event

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"testing"

	"gameclustering.com/internal/core"
)

var W sync.WaitGroup = sync.WaitGroup{}

type sampleFactory struct {
}

func (s *sampleFactory) Create(classId int, ticket string) (Event, error) {
	fmt.Printf("%d %s\n", classId, ticket)
	return &sampleEvent{}, errors.New("error")
}
func (s *sampleFactory) OnEvent(e Event) {

}

func (s *sampleFactory) OnError(e error) {
	fmt.Printf("err %s\n", e.Error())
	W.Done()
}
func (s *sampleFactory) VerifyTicket(ticket string) error {
	return nil
}

func (s *sampleFactory) Send(e Event) error {
	return nil
}
func (s *sampleFactory) View(e Query)  {
}

type sampleEvent struct {
	Id   int64
	Str  string
	I32  int32
	B    bool
	C64  complex64
	C128 complex128
	I8   int8
	I16  int16
	I64  int64
	F32  float32
	F64  float64

	Name string

	core.PersistentableObj
	EventObj
}

func (s *sampleEvent) ClassId() int {
	return 100
}
func (s *sampleEvent) WriteKey(value core.DataBuffer) error {
	value.WriteInt64(s.Id)
	return nil
}

func (s *sampleEvent) Write(value core.DataBuffer) error {
	value.WriteInt64(s.I64)
	value.WriteInt32(s.I32)
	value.WriteInt16(s.I16)
	value.WriteInt8(s.I8)
	value.WriteString(s.Str)
	value.WriteBool(s.B)
	value.WriteComplex64(s.C64)
	value.WriteComplex128(s.C128)
	value.WriteFloat32(s.F32)
	value.WriteFloat64(s.F64)
	return nil
}
func (s *sampleEvent) Read(value core.DataBuffer) error {
	i64, err := value.ReadInt64()
	if err != nil {
		return err
	}
	s.I64 = i64
	i32, err := value.ReadInt32()
	if err != nil {
		return err
	}
	s.I32 = i32
	i16, err := value.ReadInt16()
	if err != nil {
		return err
	}
	s.I16 = i16
	i8, err := value.ReadInt8()
	if err != nil {
		return err
	}
	s.I8 = i8
	str, err := value.ReadString()
	if err != nil {
		return err
	}
	s.Str = str
	b, err := value.ReadBool()
	if err != nil {
		return err
	}
	s.B = b

	c64, err := value.ReadComplex64()
	if err != nil {
		return err
	}
	s.C64 = c64

	c128, err := value.ReadComplex128()
	if err != nil {
		return err
	}
	s.C128 = c128
	f32, err := value.ReadFloat32()
	if err != nil {
		return err
	}
	s.F32 = f32
	f64, err := value.ReadFloat64()
	if err != nil {
		return err
	}
	s.F64 = f64
	return nil
}

func (s *sampleEvent) streaming(c Chunk) {
	fmt.Printf("REV : %s\n", string(c.Data))
	//s.listener <- c
}

func (s *sampleEvent) Inbound(buff core.DataBuffer) error {
	s.Read(buff)
	for {
		sz, err := buff.ReadInt32()
		if err != nil {
			s.streaming(Chunk{true, []byte{0}})
			break
		}
		if sz == 0 {
			s.streaming(Chunk{true, []byte{0}})
			break
		}
		pd, err := buff.Read(int(sz))
		if err != nil {
			s.streaming(Chunk{true, []byte{0}})
			break
		}
		s.streaming(Chunk{false, pd})
	}
	buff.WriteInt32(100)
	buff.WriteString("bye")
	W.Done()
	return nil
}

func (s *sampleEvent) Outbound(buff core.DataBuffer) error {
	s.Write(buff)
	buff.WriteInt32(12)
	buff.Write([]byte("login passed"))
	buff.WriteInt32(12)
	buff.Write([]byte("login passed"))
	buff.WriteInt32(0)
	buff.ReadInt32()
	buff.ReadString()
	return nil
}

func TestEndpoint(t *testing.T) {
	core.CreateTestLog()
	tcp := TcpEndpoint{Endpoint: "tcp://localhost:5000", Service: &sampleFactory{}}
	go tcp.Open()

	W.Add(1)
	sample1 := sampleEvent{Id: 200}
	sample1.I64 = 64
	sample1.I32 = 32
	sample1.I16 = 16
	sample1.I8 = 8
	sample1.Str = "hello"
	sample1.B = true
	sample1.C64 = 164
	sample1.C128 = 228
	sample1.F32 = 12.09
	sample1.F64 = 64.09
	soc := SocketPublisher{Remote: "tcp://localhost:5000"}
	soc.Publish(&sample1, "ticket12123131")
	W.Wait()
	tcp.Close()
}

func createEvent() Event {
	sub := SubscriptionEvent{}
	return &sub
}
func TestEventJson(t *testing.T) {
	sub := SubscriptionEvent{App: "presence", Name: "ban"}
	data, err := json.Marshal(sub)
	if err != nil {
		t.Errorf("json error %s", err.Error())
	}
	e := createEvent()
	json.Unmarshal(data, e)

	pb, isSub := e.(*SubscriptionEvent)
	if !isSub {
		t.Errorf("json error %s", err.Error())
	}
	if pb.App != sub.App {
		t.Errorf("app not there %s", err.Error())
	}
}
