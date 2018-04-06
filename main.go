package main

import (
	"flag"
	"amulet/server"
	"github.com/golang/glog"
)



func main() {
	flag.Parse()
	server := &server.Server{}
	server.Start()
	glog.Flush()
}
