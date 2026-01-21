package grpc

import (
	context "context"
	"fmt"
	"io"
	"testing"
	"time"

	gx "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestClubClient1(t *testing.T) {
	conn, err := gx.NewClient("localhost:8888", gx.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Errorf("connect error %s", err)
	}
	defer conn.Close()
	cs := NewClubServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	c, err := cs.Club(ctx, &emptypb.Empty{})
	if err != nil {
		t.Errorf("rpc error %s", err.Error())
	}
	fmt.Printf("Club %s : %d : %d\n", c.Name, c.Id, len(c.Members))

}
func TestClubClient2(t *testing.T) {
	conn, err := gx.NewClient("localhost:8888", gx.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Errorf("connect error %s", err)
	}
	defer conn.Close()
	cs := NewClubServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	p := Person{Name: "tester", Id: 100, Age: 10, Type: MemberType_Diamond}
	cb, err := cs.Register(ctx, &p)
	if err != nil {
		t.Errorf("rpc error %s", err.Error())
	}
	fmt.Printf("Register club %v\n", cb)
	cc, err := cs.Unregister(ctx, &p)
	if err != nil {
		t.Errorf("rpc error %s", err.Error())
	}
	fmt.Printf("UnRegister club %v\n", cc)

}

func TestClubClient3(t *testing.T) {
	conn, err := gx.NewClient("localhost:8888", gx.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Errorf("connect error %s", err)
	}
	defer conn.Close()
	cs := NewClubServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	streaming, err := cs.List(ctx, &emptypb.Empty{})
	if err != nil {
		t.Errorf("rpc error %s", err.Error())
	}
	for {
		p, err := streaming.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Errorf("streaming error %s", err.Error())
			break
		}
		fmt.Printf("Stream p %v\n", p)
	}
}

func TestClubClient4(t *testing.T) {
	conn, err := gx.NewClient("localhost:8888", gx.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Errorf("connect error %s", err)
	}
	defer conn.Close()
	cs := NewClubServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	streaming, err := cs.Post(ctx)
	if err != nil {
		t.Errorf("rpc error %s", err.Error())
	}

	for i := range 5 {
		p := Person{Name: "tester", Id: 100, Age: int32(i), Type: MemberType_Diamond}
		err := streaming.Send(&p)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Errorf("client streaming error %s", err.Error())
			break
		}
		fmt.Printf("Client Stream p %v\n", &p)
	}
	c, err := streaming.CloseAndRecv()
	if err != nil {
		t.Errorf("rpc error %s", err.Error())
	}
	fmt.Printf("Done %v\n", c)
}

func TestClubClient5(t *testing.T) {
	conn, err := gx.NewClient("localhost:8888", gx.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Errorf("connect error %s", err)
	}
	defer conn.Close()
	cs := NewClubServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	streaming, err := cs.Flow(ctx)
	if err != nil {
		t.Errorf("rpc error %s", err.Error())
	}
	wt := make(chan bool)
	defer close(wt)
	go func() {
		for {
			p, err := streaming.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Printf("streaming error %s", err.Error())
				break
			}
			fmt.Printf("ddd Stream p %v\n", p)
		}
		wt <- true
	}()
	for i := range 5 {
		p := Person{Name: "tester", Id: 100, Age: int32(i), Type: MemberType_Diamond}
		err := streaming.Send(&p)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Errorf("client streaming error %s", err.Error())
			break
		}
		fmt.Printf("Client Stream p %v\n", &p)
	}
	streaming.CloseSend()
	<-wt
	fmt.Printf("bi stream done")
}
