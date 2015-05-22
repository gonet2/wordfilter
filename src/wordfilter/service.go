package main

import (
	"github.com/huichen/sego"
	"golang.org/x/net/context"
	pb "proto"
	"regexp"
	"strings"
	"unicode/utf8"
)

const (
	SERVICE = "[WORDFILTER]"
)

type server struct{}

func (s *server) Filter(ctx context.Context, in *pb.WordFilter_Text) (*pb.WordFilter_Text, error) {
	segments := _segmenter.Segment([]byte(in.Text))
	clean_text := in.Text
	words := sego.SegmentsToSlice(segments, false)
	for k := range words {
		if _dirty_words[strings.ToUpper(words[k])] {
			reg, _ := regexp.Compile("(?i:" + regexp.QuoteMeta(words[k]) + ")")
			replacement := strings.Repeat("▇", utf8.RuneCountInString(words[k]))
			clean_text = reg.ReplaceAllLiteralString(clean_text, replacement)
		}
	}

	return &pb.WordFilter_Text{clean_text}, nil
}
