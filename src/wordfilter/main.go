package main

import (
	"net"
	pb "wordfilter/proto"

	log "github.com/gonet2/libs/nsq-logger"
	_ "github.com/gonet2/libs/statsd-pprof"
	"google.golang.org/grpc"
)

const (
	_port = ":50002"
)

func main() {
	log.SetPrefix(SERVICE)
	// 监听
	lis, err := net.Listen("tcp", _port)
	if err != nil {
		log.Critical(err)
	}
	log.Info("listening on ", lis.Addr())

	// 注册服务
	s := grpc.NewServer()
	ins := &server{}
	ins.init()
	pb.RegisterWordFilterServiceServer(s, ins)

	// 开始服务
	s.Serve(lis)
}
