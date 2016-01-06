package main

import (
	"bufio"
	"os"
	pb "proto"
	"strings"
	"unicode/utf8"

	"golang.org/x/net/context"

	log "github.com/gonet2/libs/nsq-logger"
	"github.com/huichen/sego"
)

const (
	SERVICE = "[WORDFILTER]"
)

var (
	replaceTo    = "*" //"▇" // "*"
	replaceByte  = []byte(strings.Repeat(replaceTo, 50))
	replaceLenth = len(replaceTo)
)

type server struct {
	dirty_words map[string]bool
	segmenter   sego.Segmenter
}

func (s *server) init() {
	s.dirty_words = make(map[string]bool)

	dict_path, dirty_words_path := s.data_path()
	// 载入字典
	log.Trace("Loading Dictionary...")
	s.segmenter.LoadDictionary(dict_path)
	log.Trace("Dictionary Loaded")

	// 读取脏词库
	log.Trace(SERVICE, "Loading Dirty Words...")
	f, err := os.Open(dirty_words_path)
	if err != nil {
		log.Critical(err)
		os.Exit(-1)
	}
	defer f.Close()

	// 逐行扫描
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		words := strings.Split(strings.ToUpper(strings.TrimSpace(scanner.Text())), " ") // 均处理为大写
		if words[0] != "" {
			s.dirty_words[words[0]] = true
		}
	}
	log.Trace("Dirty Words Loaded")
}

// get correct dict path from GOPATH
func (s *server) data_path() (dict_path string, dirty_words_path string) {
	paths := strings.Split(os.Getenv("GOPATH"), ":")
	for k := range paths {
		dirty_words_path = paths[k] + "/dirty.txt"
		_, err := os.Lstat(dirty_words_path)
		if err == nil {
			dict_path = paths[k] + "/dirty.txt," + paths[k] + "/dictionary.txt"
			return
		}
	}
	return
}

func (s *server) Filter(ctx context.Context, in *pb.WordFilter_Text) (*pb.WordFilter_Text, error) {
	bin := []byte(in.Text)
	segments := s.segmenter.Segment(bin)
	clean_text := make([]byte, 0, len(bin))
	for _, seg := range segments {
		word := bin[seg.Start():seg.End()]
		if s.dirty_words[strings.ToUpper(string(word))] {
			clean_text = append(clean_text, replaceByte[:replaceLenth*utf8.RuneCount(word)]...)
			//replacement := strings.Repeat(replaceTo, utf8.RuneCount(word))
			//clean_text = append(clean_text, []byte(replacement)...)
		} else {
			clean_text = append(clean_text, word...)
		}
	}
	return &pb.WordFilter_Text{string(clean_text)}, nil
}
