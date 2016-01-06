package main

import (
	"log"
	"proto"
	"testing"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	address  = "192.168.6.70:50002"
	testText = []string{
		"我操你大爷，法轮大法好",
		"Fuck you，fuck you sisters!",
		"개쌍또라이abasd",
	}
	conn *grpc.ClientConn
)

func init() {
	// Set up a connection to the server.
	_conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Printf("did not connect: %v", err)
		return
	}
	conn = _conn
}

func TestWordFilter(t *testing.T) {

	c := proto.NewWordFilterServiceClient(conn)

	// Contact the server and print out its response.
	for i := 0; i < 3; i++ {
		r, err := c.Filter(context.Background(), &proto.WordFilter_Text{Text: testText[i]})
		if err != nil {
			t.Fatalf("could not query: %v", err)
		}
		t.Logf("Filtered Text: %s", r.Text)
	}
}

func BenchmarkWordFilterb(b *testing.B) {
	c := proto.NewWordFilterServiceClient(conn)
	for i := 0; i < b.N; i++ {
		r, err := c.Filter(context.Background(), &proto.WordFilter_Text{Text: testText[i%3]})
		if err != nil {
			b.Fatal(err)
		}
		_ = r
	}
}
