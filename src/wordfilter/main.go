package main

import (
	"bufio"
	log "github.com/GameGophers/nsq-logger"
	"github.com/huichen/sego"
	"google.golang.org/grpc"
	"net"
	"os"
	pb "proto"
	"strings"
)

const (
	_port = ":50002"
)

var (
	_dict_path        = os.Getenv("GOPATH") + "/src/github.com/huichen/sego/data/dictionary.txt"
	_dirty_words_file = os.Getenv("GOPATH") + "/src/wordfilter/dirty.txt"
)

var (
	_dirty_words = make(map[string]bool)
	_segmenter   sego.Segmenter
)

func main() {
	// 载入字典
	log.Trace(SERVICE, "Loading Dictionary...")
	_segmenter.LoadDictionary(_dict_path)
	log.Trace(SERVICE, "Dictionary Loaded")

	// 读取脏词库
	log.Trace(SERVICE, "Loading Dirty Words...")
	f, err := os.Open(_dirty_words_file)
	if err != nil {
		log.Critical(SERVICE, err)
		os.Exit(-1)
	}
	defer f.Close()

	// 逐行扫描
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		word := strings.ToUpper(strings.TrimSpace(scanner.Text())) // 均处理为大写
		if word != "" {
			_dirty_words[word] = true
		}
	}
	log.Trace(SERVICE, "Dirty Words Loaded")

	// 监听
	lis, err := net.Listen("tcp", _port)
	if err != nil {
		log.Critical(SERVICE, err)
	}
	log.Info(SERVICE, "listening on ", lis.Addr())

	// 注册服务
	s := grpc.NewServer()
	pb.RegisterWordFilterServiceServer(s, &server{})

	// 开始服务
	s.Serve(lis)
}
