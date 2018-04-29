package core

import (
	"github.com/golang/glog"
	"os"
	"bufio"
	"strings"
)

type StopToken struct {
	Tokens map[string]bool
}

func (this *StopToken) Init() {
	glog.Info("StopToken 初始化...")
	this.Tokens = make(map[string]bool)

	file, err := os.Open("data/stopword.txt")
	if err != nil {
		glog.Info("StopToken初始化失败, error:",err.Error())
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text != "" {
			this.Tokens[text] = true
		}
	}
}

func (this *StopToken) IsStopToken(str string) bool {
	_, ok := this.Tokens[str]
	if ok {
		return true
	} else {
		return false
	}
}