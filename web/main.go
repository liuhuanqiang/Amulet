package main

import (
	"flag"
	"github.com/golang/glog"
	"amulet/web/config"
	"fmt"
	"amulet/utils"

	"amulet/web/model"
)

func main() {
	flag.Parse()
	config.LoadConfig(fmt.Sprintf("%s/%s",utils.GetCurrentDirectory(),"config.json"))
	db.GetDB().Init()
	router := &Router{}
	router.StartHttp()
	glog.Flush()
}

