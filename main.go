package main

import (
	"github.com/golang/glog"
	"flag"
	"unicode/utf8"
	"amulet/segment"
)

func main() {
	flag.Parse()
	//eg := engine.Engine{}
	//eg.Init()
	seg := &segment.Segment{}
	seg.Init()
	ret := seg.NShort("商品和服务")
	glog.Info(ret)
	glog.Flush()
}

type Text []byte

func splitTextToWords(text Text) []Text {
	output := []Text{}
	current := 0
	glog.Info("text:", text)
	length := 0
	for current < len(text) {
		r, size := utf8.DecodeRune(text[current:])
		if string(r) == "，" {
			output = append(output, text[current:current+length])
			current += length
			length = 0
		} else {
			length += size
		}

	}
	for _, v := range output {
		glog.Info(string(v))

	}
	return output
}

