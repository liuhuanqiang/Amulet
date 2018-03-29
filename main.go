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
	str := "Object.observe已经官方声明废弃，当然这也是在情理之中的，因为这个属性不可预测性太高。但是这并不意味着拥有一个可以观察的对象是一件坏事。事实上，可观察对象是一个非常强大的概念。别担心，"
	seg := &segment.Segment{}
	seg.Init()
	ret := seg.SplitToSegment([]byte(str))
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

