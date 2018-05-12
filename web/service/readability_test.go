package service

import (
	"testing"
)

func TestReadability_GetContent(t *testing.T) {
	//url := "https://blog.qiniu.com/archives/8728"
	//url := "https://www.skiy.net/201511103817.html"
	//url := "http://forthxu.com/blog/article/73.html"
	url := "https://cnodejs.org/topic/550bdbaffe20995b17bf8d32"

	content := &Readability{}
	content.GetContent(url)
}