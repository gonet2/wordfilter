package main

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"testing"
	pb "wordfilter/proto"
)

const (
	address   = "localhost:50002"
	test_text = "我操你大爷，法轮大法好"
)

func TestWordFilter(t *testing.T) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address)
	if err != nil {
		t.Fatal("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewWordFilterServiceClient(conn)

	// Contact the server and print out its response.
	r, err := c.Filter(context.Background(), &pb.WordFilter_Text{Text: test_text})
	if err != nil {
		t.Fatalf("could not query: %v", err)
	}
	t.Logf("Filtered Text: %s", r.Text)
}
