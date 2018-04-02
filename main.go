package main

import (
	"github.com/golang/glog"
	"flag"
	"amulet/segment"
	"path/filepath"
	"os"
	"bufio"
	"io"
	"github.com/axgle/mahonia"
	"time"
)

func main() {
	flag.Parse()
	//eg := engine.Engine{}
	//eg.Init()
	startTime := time.Now()
	seg := &segment.Segment{}
	seg.Init()

	ng := &segment.Ngrams{}
	ng.Init()

	paths:= []string{}
	filepath.Walk("./data/sougou/", func(path string, f os.FileInfo, err error) error {
		if ( f == nil ) {return err}
		if f.IsDir() {return nil}
		paths = append(paths,path)
		return nil
	})
	glog.Info("length:",len(paths))
	for _, path := range paths {
		//path := paths[0]
		fi, err := os.Open(path)
		if err != nil {
			glog.Info("error:",err.Error(), " path:", path)
			return
		}

		decoder := mahonia.NewDecoder("gbk")
		br := bufio.NewReader(decoder.NewReader(fi))
		sum := ""
		for {
			a, _, c := br.ReadLine()
			if c == io.EOF {
				break
			}
			sum += string(a)
		}
		fi.Close()
		ret := seg.SplitToSegment([]byte(sum))
		for _, v := range ret {
			ng.Start(v)
		}
	}
	glog.Info("211321")
	ng.Judge()
	glog.Info("cost:", time.Since(startTime))
	glog.Flush()
}


