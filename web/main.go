package main

import (
	"flag"
	"github.com/golang/glog"
)

func main() {
	flag.Parse()
	router := &Router{}
	router.StartHttp()
	glog.Flush()
}

