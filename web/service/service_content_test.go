package service

import (
	"testing"
	"fmt"
)

func TestServiceContent_GetContent(t *testing.T) {
	//url := "http://www.trinea.cn/jobs/%E6%8E%A8%E8%8D%90-3-%E4%B8%AA%E7%AE%80%E5%8E%86%E6%A8%A1%E6%9D%BF%E5%8F%8A-2-%E5%A4%A7%E5%8A%A0%E5%88%86%E6%8A%80%E5%B7%A7/"
	//url := "https://blog.ibireme.com/2017/09/01/diary/"
	//url := "http://stormzhang.com/2017/07/25/instant-apps-vs-wechat-application/"
	//url := "http://blog.devtang.com/2018/04/16/operation-light-summary/"
	//url := "http://justjavac.com/javascript/2018/02/09/re-talk-about-reactdom-render.html"
	//url := "https://www.liaoxuefeng.com/article/001509844125769eafbb65df0a04430a2d010a24a945bfa000"
	//url := "https://blog.csdn.net/solstice/article/details/8493251"
	//url := "https://drakeet.me/pure-writer-3/"
	//url := "http://hencoder.com/ui-2-3/"
	//url := "https://droidyue.com/blog/2018/04/01/do-not-bother-the-user-when-app-crash-in-a-background-state/"
	//url := "http://www.ruanyifeng.com/blog/2018/04/weekly-issue-1.html"
	//url := "https://tech.meituan.com/ruby_autotest.html"
	//url := "http://www.alloyteam.com/2018/04/gka-optimize/"
	//url := "https://www.barretlee.com/blog/2018/03/01/%E9%99%AA%E4%BC%B4/"
	//url := "http://blog.imallen.wang/2018/02/07/SkipList%E7%9A%84%E5%8E%9F%E7%90%86%E4%B8%8E%E5%AE%9E%E7%8E%B0/"
	//url := "http://mindhacks.cn/2017/10/17/through-the-maze-11/"
	//url := "http://www.wjdiankong.cn/archives/1125"
	//url := "http://gityuan.com/2017/07/11/android_debug/"
	//url := "http://melonteam.com/posts/chu_tan_kotlin_yi_bu_async_await/"
	//url := "http://cdc.tencent.com/2018/04/26/%E6%B3%9B%E5%A8%B1%E4%B9%90%E5%BE%AE%E4%BF%A1%E5%BA%97%E6%94%B9%E7%89%88-%E8%AE%BE%E8%AE%A1%E6%80%BB%E7%BB%93/"
	//url := "https://isux.tencent.com/articles/tencent-docs.html"
	//url := "http://qqfe.org/archives/371"
	//url := "http://fex.baidu.com/blog/2018/04/fex-weekly-30/"
	//url := "http://taobaofed.org/blog/2018/03/12/long-list-in-rax/"
	//url := "https://wusay.org/skiplist.html"
	//url := "http://yulingtianxia.com/blog/2018/03/31/Track-Block-Arguments-of-Objective-C-Method/"
	//url := "http://blog.leichunfeng.com/blog/2017/02/20/talking-about-the-decompression-of-the-image-in-ios/"
	//url := "http://blog.startry.com/2017/01/10/2016-conclusion/"
	//url := "http://www.iosxxx.com/blog/2017-09-29-OpenGLES%E4%B9%8B%E5%9F%BA%E7%A1%80%E7%AF%87.html"
	//url := "http://blog.raozhizhen.com/post/2016-08-19"
	//url := "http://msching.github.io/blog/2016/05/24/audio-in-ios-9/"
	//url := "http://coderyi.com/posts/weex3/"
	//url := "https://casatwy.com/Advance_In_iOS11_Networking.html"
	//url := "http://blog.cnbang.net/writting/3565/"
	//url := "http://f2e.souche.com/blog/webpackbian-yi-vuexiang-mu-sheng-cheng-de-dai-ma-tan-suo/"
	//url := "https://75team.com/post/webview-debug.html"
	//url := "http://singsing.io/blog/fcc/advanced-pairwise/"
	//url := "https://imququ.com/post/reduce-ttfb-on-thinkjs3-website.html"
	//url := "https://luolei.org/terramaster-d2-310-review/"
	//url := "https://lukesign.com/disable-wechat-font-adjust/"
	//url := "http://ued.ctrip.com/?p=5657"
	//url := "https://aotu.io/notes/2018/04/24/jdindex_2017/"
	//url := "http://kaedea.com/2017/04/09/android/global-accessing-context/"
	//url := "https://blog.codingnow.com/2018/05/frostpunk.html"
	//url := "https://www.phodal.com/blog/serverless-wechat-app-login/"
	//url := "http://jm.taobao.org/2018/05/02/post-20180502/"
	url := "https://blog.qiniu.com/archives/8728"

	content := &ServiceContent{}
	title, des := content.GetContent(88,url)
	fmt.Println(title)
	fmt.Println(des)
}
