package service

import (
	"testing"
)

func TestReadability_GetContent(t *testing.T) {
	//url := "https://blog.qiniu.com/archives/8728"
	//url := "https://www.skiy.net/201511103817.html"
	//url := "http://forthxu.com/blog/article/73.html"
	//url := "https://cnodejs.org/topic/550bdbaffe20995b17bf8d32"
	//url := "http://blog.battcn.com/2018/04/23/springboot/v2-config-logs/"
	url := "http://blog.devtang.com/2018/05/12/human-history/"
	content := &Readability{}
	content.GetContent(url)
}