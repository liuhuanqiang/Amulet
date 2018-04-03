package main

import (
	"github.com/golang/glog"
	"flag"
	"amulet/segment"
)

func main() {
	flag.Parse()
	//eg := engine.Engine{}
	//eg.Init()
	seg := &segment.Segment{}
	seg.Init()

	sc := &segment.Scaner{}
	sc.Init()
	sc.Filter()
	glog.Info("finish")
	glog.Flush()
}


