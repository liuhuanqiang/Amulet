package main

import (
	"github.com/golang/glog"
	"flag"
	"amulet/engine"
)

func main() {
	flag.Parse()
	eg := engine.Engine{}
	eg.Init()
	glog.Flush()
}

