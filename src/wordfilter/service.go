package main

import (
	"bufio"
	log "github.com/GameGophers/nsq-logger"
	"github.com/huichen/sego"
	"golang.org/x/net/context"
	"os"
	pb "proto"
	"regexp"
	"strings"
	"unicode/utf8"
)

const (
	SERVICE = "[WORDFILTER]"
)

type server struct {
	dirty_words map[string]bool
	segmenter   sego.Segmenter
}

func (s *server) init() {
	s.dirty_words = make(map[string]bool)

	dict_path, dirty_words_path := s.data_path()
	// 载入字典
	log.Trace(SERVICE, "Loading Dictionary...")
	s.segmenter.LoadDictionary(dict_path)
	log.Trace(SERVICE, "Dictionary Loaded")

	// 读取脏词库
	log.Trace(SERVICE, "Loading Dirty Words...")
	f, err := os.Open(dirty_words_path)
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
			s.dirty_words[word] = true
		}
	}
	log.Trace(SERVICE, "Dirty Words Loaded")
}

// get correct dict path from GOPATH
func (s *server) data_path() (dict_path string, dirty_words_path string) {
	paths := strings.Split(os.Getenv("GOPATH"), ":")
	for k := range paths {
		dirty_words_path = paths[k] + "/src/wordfilter/dirty.txt"
		_, err := os.Lstat(dirty_words_path)
		if err == nil {
			dict_path = paths[k] + "/src/wordfilter/dictionary.txt"
			return
		}
	}
	return
}

func (s *server) Filter(ctx context.Context, in *pb.WordFilter_Text) (*pb.WordFilter_Text, error) {
	segments := s.segmenter.Segment([]byte(in.Text))
	clean_text := in.Text
	words := sego.SegmentsToSlice(segments, false)
	for k := range words {
		if s.dirty_words[strings.ToUpper(words[k])] {
			reg, _ := regexp.Compile("(?i:" + regexp.QuoteMeta(words[k]) + ")")
			replacement := strings.Repeat("▇", utf8.RuneCountInString(words[k]))
			clean_text = reg.ReplaceAllLiteralString(clean_text, replacement)
		}
	}

	return &pb.WordFilter_Text{clean_text}, nil
}
