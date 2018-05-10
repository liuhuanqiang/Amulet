package service

import (
	"testing"
)

func TestReadability_GetContent(t *testing.T) {
	//url := "https://blog.qiniu.com/archives/8728"
	url := "https://www.skiy.net/201511103817.html"
	content := &Readability{}
	content.GetContent(url)
}