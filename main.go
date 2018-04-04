package main

import (
	"github.com/golang/glog"
	"flag"
	"amulet/segment"
	"os"
	"path/filepath"
	"github.com/axgle/mahonia"
	"bufio"
	"io"
	"sync"
	"strings"
	"time"
	"amulet/core"
)


var seg = &segment.Segment{}
var chs = make(chan *Elem,6)
var wg sync.WaitGroup
var t = 0

type Elem struct {
	Path   string
	strs   string
}
func read(path string) string {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE,0766)
	if err != nil {
		glog.Info("err:", err.Error(), " path:", path)
	}
	decoder := mahonia.NewDecoder("gbk")
	reader := bufio.NewReader(decoder.NewReader(file))
	var str string
	for {
		var tmp string
		tmp,err:= reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		str += string(tmp)
	}
	file.Close()
	return str
}

// 写文件
func write(strs [][]string, path string) {
	path = strings.Replace(path, "sougou", "sougouhandle", -1)
	fl, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE,0766)
	if err != nil {
		glog.Info("error:", err.Error())
		return
	}
	sum := ""
	for _, v := range strs {
		sum = sum + "\n"
		for _, s := range v {
			sum = sum + "\\" + s
		}
	}
	_, err1 := fl.Write([]byte(sum))
	if err1 != nil {
		glog.Info("error:", err1.Error())
	}
}


func Start() {
	t1 := time.Now() // get current time
	seg.Init()
	dirPath := "data/sougou/"
	paths := []string{}
	filepath.Walk(dirPath, func(filename string, fi os.FileInfo, err error) error {
		//遍历目录
		if fi.IsDir() {
			// 忽略目录
			return nil
		}
		paths = append(paths, filename)
		return nil
	})
	glog.Info("paths:", len(paths))
	wg.Add(1)

	go func() {
		for _, path := range paths {
			str := read(path)
			e := &Elem{}
			e.Path = path
			e.strs = str
			chs <- e
		}
	}()

	go func() {
		for v := range chs {
			//glog.Info("读取")
			t++
			write(seg.Cut([]byte(v.strs)),v.Path)
			if t >= len(paths) {
				wg.Done()
			}
		}
	}()
	wg.Wait()
	glog.Info("finish:", time.Since(t1))
}

func main() {
	flag.Parse()
	//eg := engine.Engine{}
	//eg.Init()
	seg := &segment.Segment{}
	seg.Init()
	st := &core.StopToken{}
	st.Init()
	str := "60烽火实验室 第一章 勒索中的“代刷"
	//glog.Info("str:", str)
	ret := seg.Cut([]byte(str))
	glog.Info("ret:",ret)
	list := []string{}
	for _, v := range ret {
		for _,s := range v {
			if  len(s) > 0 && !st.IsStopToken(s) {
				list = append(list, s)
			}
		}
	}
	tr := &core.TextRank{}
	glog.Info(tr.GetRankList(list))

	glog.Flush()
}
