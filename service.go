package main

import (
	"bufio"
	"os"
	"strings"
	"unicode/utf8"
	pb "wordfilter/proto"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/huichen/sego"
)

const (
	SERVICE = "[WORDFILTER]"
)

var (
	replaceTo   = "*" //"▇" // "*"
	replaceByte = []byte(strings.Repeat(replaceTo, 50))
)

type server struct {
	dirty_words map[string]bool
	segmenter   sego.Segmenter
}

func (s *server) init() {
	s.dirty_words = make(map[string]bool)

	dict_path, dirty_words_path := s.data_path()
	// 载入字典
	log.Debug("Loading Dictionary...")
	s.segmenter.LoadDictionary(dict_path)
	log.Debug("Dictionary Loaded")

	// 读取脏词库
	log.Debug("Loading Dirty Words...")
	f, err := os.Open(dirty_words_path)
	if err != nil {
		log.Panic(err)
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
	log.Debug("Dirty Words Loaded")
}

// get correct dict path from GOPATH
func (s *server) data_path() (dict_path string, dirty_words_path string) {
	paths := strings.Split(os.Getenv("GOPATH"), ":")
	for k := range paths {
		dirty_words_path = paths[k] + "/src/wordfilter/dirty.txt"
		_, err := os.Lstat(dirty_words_path)
		if err == nil {
			dict_path = paths[k] + "/src/wordfilter/dirty.txt," + paths[k] + "/src/wordfilter/dictionary.txt"
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
			clean_text = append(clean_text, replaceByte[:utf8.RuneCount(word)]...)
		} else {
			clean_text = append(clean_text, word...)
		}
	}
	return &pb.WordFilter_Text{string(clean_text)}, nil
}
