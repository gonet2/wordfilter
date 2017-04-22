package main

import (
	"bufio"
	"os"
	"strings"
	"unicode/utf8"
	pb "wordfilter/proto"

	cli "gopkg.in/urfave/cli.v2"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/huichen/sego"
)

type server struct {
	replaceWord string
	dirty_words map[string]bool
	segmenter   sego.Segmenter
}

func (s *server) init(c *cli.Context) {
	s.replaceWord = c.String("replace-word")
	s.dirty_words = make(map[string]bool)

	dictionary := c.String("dictionary")
	dirty := c.String("dirty")

	// 载入字典
	log.Info("Loading Dictionary...")
	s.segmenter.LoadDictionary(dictionary)
	log.Info("Dictionary Loaded")

	// 读取脏词库
	log.Info("Loading Dirty Words...")
	f, err := os.Open(dirty)
	if err != nil {
		log.Fatalln(err)
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
	log.Info("Dirty Words Loaded")
}

func (s *server) Filter(ctx context.Context, in *pb.WordFilter_Text) (*pb.WordFilter_Text, error) {
	bin := []byte(in.Text)
	segments := s.segmenter.Segment(bin)
	clean_text := make([]byte, 0, len(bin))
	for _, seg := range segments {
		word := bin[seg.Start():seg.End()]
		if s.dirty_words[strings.ToUpper(string(word))] {
			clean_text = append(clean_text, []byte(strings.Repeat(s.replaceWord, utf8.RuneCount(word)))...)
		} else {
			clean_text = append(clean_text, word...)
		}
	}
	return &pb.WordFilter_Text{string(clean_text)}, nil
}
