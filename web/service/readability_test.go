package service

import (
	"testing"
)

func TestReadability_GetContent(t *testing.T) {
	url := "https://blog.qiniu.com/archives/8728"

	content := &Readability{}
	content.GetContent(url)
}