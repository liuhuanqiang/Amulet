package domain

import (
	"testing"
	"fmt"
)

func TestContent_GetContent(t *testing.T) {
	//url := "http://www.trinea.cn/jobs/%E6%8E%A8%E8%8D%90-3-%E4%B8%AA%E7%AE%80%E5%8E%86%E6%A8%A1%E6%9D%BF%E5%8F%8A-2-%E5%A4%A7%E5%8A%A0%E5%88%86%E6%8A%80%E5%B7%A7/"
	//url := "https://blog.ibireme.com/2017/09/01/diary/"
	//url := "http://stormzhang.com/2017/07/25/instant-apps-vs-wechat-application/"
	//url := "http://blog.devtang.com/2018/04/16/operation-light-summary/"
	//url := "http://justjavac.com/javascript/2018/02/09/re-talk-about-reactdom-render.html"
	//url := "https://www.liaoxuefeng.com/article/001509844125769eafbb65df0a04430a2d010a24a945bfa000"
	//url := "https://blog.csdn.net/solstice/article/details/8493251"
	url := "https://drakeet.me/pure-writer-3/"
	content := &Content{}
	title, des := content.GetContent(url)
	fmt.Println(title)
	fmt.Println(des)
}
