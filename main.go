package main

import (
	"flag"
	"amulet/web"
	"github.com/golang/glog"
)



func main() {
	flag.Parse()
	server := &web.Server{}
	server.Start()
	glog.Flush()
}
