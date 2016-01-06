package main

import (
	pb "proto"
	"testing"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address   = "localhost:50002"
	test_text = "angel dust,k.u.n.t, alabama black	snake,我操你大爷함대줄래，al qaeda法轮大法好,alabama black snake  ass hole fuck you ass hole, ass hole"
)

func TestWordFilter(t *testing.T) {
	opt := grpc.WithInsecure()
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, opt)
	if err != nil {
		t.Fatalf("did not connect: %v", err)
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
