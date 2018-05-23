package service

import "testing"

func TestMaxSubSegment_GetContent(t *testing.T) {
	//url := "https://imtx.me/archives/2690.html"
	//url := "http://xuzhibin.github.io/2010/11/11/text-extraction/"
	//url := "http://www.laruence.com/2018/04/08/3179.html"
	//url := "https://www.jianshu.com/p/229bdeb844dc"
	//url := "http://www.laruence.com/2018/04/08/3179.html"
	//url := "https://blog.qiniu.com/archives/8741"
	//url := "http://sec-redclub.com/archives/902/"
	//url := "http://www.haorooms.com/post/centos_git"
	//url := "http://blog.battcn.com/2018/05/22/springboot/v2-queue-rabbitmq/"
	url := "http://www.cnblogs.com/zichi/p/9068481.html"

	mss := &MaxSubSegment{}
	mss.GetContent(url)
}