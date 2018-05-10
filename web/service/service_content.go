package service

import (
	"github.com/PuerkitoBio/goquery"
	"strings"
	"net/http"
	"fmt"
	"regexp"
	"github.com/golang/glog"
	"bytes"
	"amulet/utils"
	"github.com/axgle/mahonia"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
)

type ServiceContent struct {
	UnLikeyCandidates *regexp.Regexp
	RegStyle *regexp.Regexp
	RegScript *regexp.Regexp
}


func (this *ServiceContent) GetContent(fid int, url string) (string, string) {
	title:= ""
	content := ""
	if fid == 0{
		this.getContent(url);
	}else if fid == 1 {
		title, content = this.getContent1(url)
	} else if fid == 2 {
		title, content = this.getContent2(url)
	} else if fid == 6 {
		title, content = this.getContent6(url)
	} else if fid == 7 || fid == 10  {
		title, content = this.getContent7(url)
	} else if fid == 11 {
		title, content = this.getContent11(url)
	} else if fid == 13 {
		title, content = this.getContent13(url)
	} else if fid == 14 {
		title, content = this.getContent14(url)
	} else if fid == 16 {
		title, content = this.getContent16(url)
	} else if fid == 17 {
		title, content = this.getContent17(url)
	} else if fid == 18 {
		title, content = this.getContent18(url)
	} else if fid == 19 || fid == 22 || fid == 23 {
		title, content = this.getCSDNContent(url)
	} else if fid == 27 {
		title, content = this.getContent27(url)
	} else if fid == 29 {
		title, content = this.getContent29(url)
	} else if fid == 30 {
		title, content = this.getContent30(url)
	} else if fid == 31 {
		title, content = this.getContent31(url)
	} else if fid == 32 {
		title, content = this.getContent32(url)
	} else if fid == 33 {
		title, content = this.getContent33(url)
	} else if fid == 34 {
		title, content = this.getContent34(url)
	} else if fid == 35 || fid == 49 {
		title, content = this.getNextMuseContent(url)
	} else if fid == 37 {
		title, content = this.getContent37(url)
	} else if fid == 39 {
		title, content = this.getContent39(url)
	} else if fid == 40  {
		title, content = this.getHuxContent(url)
	} else if fid == 41 {
		title, content = this.getContent41(url)
	} else if fid == 42 {
		title, content = this.getContent42(url)
	} else if fid == 43 {
		// todo
	} else if fid == 44 {
		title, content = this.getContent44(url)
	} else if fid == 45 {
		title, content = this.getContent45(url)
	} else if fid == 46 {
		title, content = this.getContent46(url)
	} else if fid == 47  {
		title, content = this.getMWebContent(url)
	} else if fid == 48 {
		title, content = this.getContent48(url)
	} else if fid == 50 || fid == 61 {
		title, content = this.getJacmanContent(url)
	} else if fid == 52 {
		// todo
	} else if fid == 55 {
		title, content = this.getOtcopressContent(url)
	} else if fid == 57 {
		title, content = this.getJekyllContent(url)
	} else if fid == 58 || fid == 72 {
		title, content = this.getYiliaContent(url)
	} else if fid == 59 {
		title, content = this.getMaterialContent(url)
	} else if fid == 60 {
		// todo
	} else if fid == 62 {
		title, content = this.getContent62(url)
	} else if fid == 63 {
		title, content = this.getContent63(url)
	} else if fid == 64 {
		title, content = this.getContent64(url)
	} else if fid == 66 {
		title, content = this.getContent66(url)
	} else if fid == 67 {
		title, content = this.getContent67(url)
	} else if fid == 68 {
		title, content = this.getContent68(url)
	} else if fid == 70 {
		title, content = this.getContent70(url)
	} else if fid == 71 {
		title, content = this.getThinkJsContent(url)
	} else if fid == 73 {
		title, content = this.getContent73(url)
	} else if fid == 74 {
		title, content = this.getThinkJsContent(url)
	} else if fid == 76 {
		title, content = this.getContent76(url)
	} else if fid == 77 {
		title, content = this.getContent77(url)
	} else if fid == 78 {
		title, content = this.getHexoContent(url)
	} else if fid == 79 {
		title, content = this.getContent79(url)
	} else if fid == 80 {
		title, content = this.getContent80(url)
	} else if fid == 82 {
		//todo
	} else if fid == 83 {
		title, content = this.getHexoContent(url)
	} else if fid == 84 {
		// todo
	} else if fid == 85 {
		title, content = this.getHuxContent(url)
	} else if fid == 86 {
		title, content = this.getHexoContent(url)
	} else if fid == 87 {
		title, content = this.getHexoContent(url)
	} else if fid == 88 {
		title, content = this.getContent88(url)
	} else {
		// fid = 54
		title, content = this.getHexoContent(url)
	}

	return title,content
}

func (this *ServiceContent) getResponse (url string) *http.Response {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, bytes.NewReader([]byte("")));
	req.Header.Set("User-Agent","Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36")
	resp, _ := client.Do(req)
	return resp
}

func (this *ServiceContent) getHtml(url string) string {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, bytes.NewReader([]byte("")));
	req.Header.Set("User-Agent","Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36")
	resp, _ := client.Do(req)
	str,_ := ioutil.ReadAll(resp.Body)
	return string(str)
}


func (this *ServiceContent) getHexoContent(url string) (string, string) {
	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title := doc.Find("article").Find(".post-title").Text()
	str, _ := doc.Find("article").Find(".post-body").Html()
	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	return strings.TrimSpace(strings.Replace(title,"\n","",0)), strings.TrimSpace(str)
}

// http://blog.imallen.wang/2018/02/07/SkipList%E7%9A%84%E5%8E%9F%E7%90%86%E4%B8%8E%E5%AE%9E%E7%8E%B0/
func (this *ServiceContent) getNextMuseContent(url string) (string, string) {
	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title,_ := doc.Find("article").Find(".post-title").Html()
	str, _ := doc.Find("article").Find(".post-body").Each(func(_ int, tag *goquery.Selection) {
		tag.Find("figure").Find(".gutter").Remove()

	}).Html()

	reg_pre := regexp.MustCompile(`(?sU:<figure.*>)(?s:.*?)(?U:</figure>)`)
	str = reg_pre.ReplaceAllStringFunc(str, func(str string) string {
		str = regexp.MustCompile(`(?U:</*figcaption.*>)`).ReplaceAllString(str, "")
		str = regexp.MustCompile(`(?U:</*span.*>)`).ReplaceAllString(str, "")
		str = regexp.MustCompile(`(?U:</*td.*>)`).ReplaceAllString(str, "")
		str = regexp.MustCompile(`(?U:</*tr.*>)`).ReplaceAllString(str, "")
		str = regexp.MustCompile(`(?U:</*table.*>)`).ReplaceAllString(str, "")
		str = regexp.MustCompile(`(?U:</*tbody.*>)`).ReplaceAllString(str, "")
		str = regexp.MustCompile(`(?U:</*figure.*>)`).ReplaceAllString(str, "")
		return  str
	})
	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}

//http://gityuan.com/2017/07/11/android_debug/
func (this *ServiceContent) getHuxContent(url string) (string, string) {
	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title := doc.Find(".post-heading").Find("h1").Text();
	str,_ := doc.Find("article").Find(".post-container").Each(func(_ int, tag *goquery.Selection) {
		tag.Find(".pager").Remove()
	}).Html()

	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}

// https://wusay.org/skiplist.html
func (this *ServiceContent) getMWebContent(url string) (string, string) {
	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title := doc.Find("article").Find(".entry-title").Text();
	str,_ := doc.Find("article").Find(".entry-content").Html()

	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}

// http://yulingtianxia.com/blog/2018/04/30/Colorful-Rounded-Rect-Dash-Border/
func (this *ServiceContent) getJacmanContent(url string) (string, string) {
	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title := doc.Find("article").Find("h1").Text();
	str,_ := doc.Find("article").Find(".article-content").Each(func(_ int, tag *goquery.Selection) {
		tag.Find("#toc").Remove()
		tag.Find("figure").Find(".gutter").Remove()
	}).Html()
	reg_pre := regexp.MustCompile(`(?sU:<figure.*>)(?s:.*?)(?U:</figure>)`)
	str = reg_pre.ReplaceAllStringFunc(str, func(str string) string {
		str = regexp.MustCompile(`(?U:</*figcaption.*>)`).ReplaceAllString(str, "")
		str = regexp.MustCompile(`(?U:</*span.*>)`).ReplaceAllString(str, "")
		str = regexp.MustCompile(`(?U:</*td.*>)`).ReplaceAllString(str, "")
		str = regexp.MustCompile(`(?U:</*tr.*>)`).ReplaceAllString(str, "")
		str = regexp.MustCompile(`(?U:</*table.*>)`).ReplaceAllString(str, "")
		str = regexp.MustCompile(`(?U:</*tbody.*>)`).ReplaceAllString(str, "")
		str = regexp.MustCompile(`(?U:</*figure.*>)`).ReplaceAllString(str, "")
		return  str
	})
	reg_br := regexp.MustCompile(`(?sU:</*br.*>)`)
	str = reg_br.ReplaceAllString(str, "\n")
	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegH1(str)
	str = utils.RegSpan(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}

func (this *ServiceContent) getOtcopressContent(url string) (string, string) {
	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title := doc.Find("article").Find(".entry-title").Text();
	str,_ := doc.Find("article").Find(".entry-content").Each(func(_ int, tag *goquery.Selection) {
		tag.Find("#toc").Remove()
		tag.Find("figure").Find(".gutter").Remove()
	}).Html()
	reg_pre := regexp.MustCompile(`(?sU:<figure.*>)(?s:.*?)(?U:</figure>)`)
	str = reg_pre.ReplaceAllStringFunc(str, func(str string) string {
		str = regexp.MustCompile(`(?U:</*figcaption.*>)`).ReplaceAllString(str, "")
		str = regexp.MustCompile(`(?U:</*span.*>)`).ReplaceAllString(str, "")
		str = regexp.MustCompile(`(?U:</*td.*>)`).ReplaceAllString(str, "")
		str = regexp.MustCompile(`(?U:</*tr.*>)`).ReplaceAllString(str, "")
		str = regexp.MustCompile(`(?U:</*table.*>)`).ReplaceAllString(str, "")
		str = regexp.MustCompile(`(?U:</*tbody.*>)`).ReplaceAllString(str, "")
		str = regexp.MustCompile(`(?U:</*figure.*>)`).ReplaceAllString(str, "")
		return  str
	})
	reg_br := regexp.MustCompile(`(?sU:</*br.*>)`)
	str = reg_br.ReplaceAllString(str, "\n")
	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegH1(str)
	str = utils.RegSpan(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}

func (this *ServiceContent) getJekyllContent(url string) (string, string) {
	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title := doc.Find("article").Find(".post-title").Text();
	str,_ := doc.Find("article").Find(".post").Html()

	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}

// "http://singsing.io/blog/fcc/advanced-pairwise/"
func (this *ServiceContent) getYiliaContent(url string) (string, string) {
	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title := doc.Find("article").Find(".article-title").Text();
	str,_ := doc.Find("article").Find(".article-entry").Each(func(_ int, tag *goquery.Selection) {
		tag.Find("#toc").Remove()
	}).Html()
	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}

// https://75team.com/post/webview-debug.html
func (this *ServiceContent) getThinkJsContent(url string) (string, string) {
	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title := doc.Find("article").Find(".title").Text();
	str,_ := doc.Find("article").Find(".entry-content").Each(func(_ int, tag *goquery.Selection) {
		tag.Find("#toc").Remove()
	}).Html()
	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}

func (this *ServiceContent) getMaterialContent(url string) (string, string) {
	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title := doc.Find(".material-post_container").Find(".article-headline-p").Text();
	str,_ := doc.Find(".material-post_container").Find("#post-content").Html()
	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}

// http://rango.swoole.com/
func (this *ServiceContent) getContent1(url string) (string, string){
	doc, _ := goquery.NewDocument(url)
	title,_ := doc.Find(".post").Find("h2").Html();
	str,_ := doc.Find(".post").Find(".content").Html()

	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	return title, strings.TrimSpace(str)
}

// http://www.laruence.com/
func (this *ServiceContent) getContent2(url string) (string, string) {
	doc, _ := goquery.NewDocument(url)
	title,_ := doc.Find(".content").Find("h1").Find("a").Html();
	str,_ := doc.Find(".content").Find(".post").Each(func(_ int, tag *goquery.Selection) {
		tag.Find("h1").Remove()
		tag.Find(".related_post_title").Remove()
		tag.Find(".related_post").Remove()
		tag.Find(".postmeta").Remove()
		tag.Find(".bds_more").Remove()
		tag.Find(".bds_qzone").Remove()
		tag.Find(".bds_tsina").Remove()
		tag.Find(".bds_tqq").Remove()
		tag.Find(".bds_renren").Remove()
		tag.Find(".shareCount").Remove()
	}).Html()

	str = utils.RegDiv(str)
	str = utils.RegH1(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	return title, strings.TrimSpace(str)
}

func (this *ServiceContent) getContent34(url string) (string, string) {
	doc, _ := goquery.NewDocument(url)
	title, _ := doc.Find("article").Find("h1").Html()
	str, _ := doc.Find("article").Find(".article-entry").Html()
	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	return strings.TrimSpace(strings.Replace(title,"\n","",0)), strings.TrimSpace(str)
}

func (this *ServiceContent) getContent6(url string) (string, string) {
	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title := doc.Find(".entry").Find(".entry-header").Text();
	str,_ := doc.Find(".entry").Find(".entry-content").Html()
	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	dec := mahonia.NewDecoder("gbk")
	return strings.TrimSpace(dec.ConvertString(title)), strings.TrimSpace(dec.ConvertString(str))
}

func (this *ServiceContent) getContent7(url string) (string, string) {
	doc, _ := goquery.NewDocument(url)
	title, _ := doc.Find("article").Find("h1").Html()
	str, _ := doc.Find("article").Find(".entry-content").Html()
	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	return strings.TrimSpace(strings.Replace(title,"\n","",0)), strings.TrimSpace(str)
}

func (this *ServiceContent) getContent11(url string) (string, string) {
	doc, _ := goquery.NewDocument(url)
	title, _ := doc.Find("h1").Html()
	str, _ := doc.Find(".container.typo").Html()
	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	return strings.TrimSpace(strings.Replace(title,"\n","",0)), strings.TrimSpace(str)
}

func (this *ServiceContent) getContent13(url string) (string, string) {
	doc, _ := goquery.NewDocument(url)
	title, _ := doc.Find("section").Find("h1").Html()
	str, _ := doc.Find("section").Each(func(_ int, tag *goquery.Selection) {
		tag.Find("h1").Remove()
		tag.Find("p").Each(func(_ int, p *goquery.Selection){
			s, _ := p.Html()
			glog.Info("str:",s)
			if strings.Contains(s, "友金所") || strings.Contains(s, "下一篇") ||strings.Contains(s, "上一篇") {
				p.Remove()
			}
		})

	}).Html()
	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	return strings.TrimSpace(strings.Replace(title,"\n","",0)), strings.TrimSpace(str)
}

func (this *ServiceContent) getContent14(url string) (string, string) {
	doc, _ := goquery.NewDocument(url)
	title := doc.Find("article").Find("h1").Text();
	str,_ := doc.Find("article").Find(".article-content").Each(func(_ int, tag *goquery.Selection) {
		tag.Find(".toc-article").Remove()
	}).Html()

	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}

func (this *ServiceContent) getContent16(url string) (string, string) {
	doc, _ := goquery.NewDocument(url)
	title := doc.Find(".content").Find(".page-header").Find("h1").Text()
	str, _ := doc.Find("#article").Html()
	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}

func (this *ServiceContent) getContent17(url string) (string, string) {

	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title,_ := doc.Find(".x-center").Find(".x-content").Find("h3").Html()
	str, _ := doc.Find(".x-center").Find(".x-article-content.x-main-content").Html()
	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}

func (this *ServiceContent) getContent18(url string) (string, string) {
	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title,_ := doc.Find(".middle.mdl-layout__content ").Find("h2").Html()
	str, _ := doc.Find(".mdl-card__supporting-text").Each(func(_ int, tag *goquery.Selection) {
		tag.Find("span.post-meta").Remove()
	}).Html()
	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}

func (this *ServiceContent) getCSDNContent(url string) (string,string){
	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title,_ := doc.Find(".title-article").Html()
	str, _ := doc.Find("article").Find(".article_content").Html()
	str = utils.RegBr(str)
	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}

func (this *ServiceContent) getContent24(url string) (string, string) {
	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title,_ := doc.Find(".post-title").Html()
	str, _ := doc.Find("article").Find(".post-content").Html()
	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}

func (this *ServiceContent) getContent27(url string) (string, string) {
	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title,_ := doc.Find("article").Find(".inner").Find(".post-title").Html()
	str, _ := doc.Find("article").Find(".post-content.inner").Html()
	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}


func (this *ServiceContent) getContent29(url string) (string, string) {
	doc, _ := goquery.NewDocument(url)
	title, _ := doc.Find("article").Find("h1").Html()
	str, _ := doc.Find("article").Find(".entry-content").Each(func(_ int, tag *goquery.Selection) {
		tag.Find("span.line-number").Remove()
		tag.Find(".line-numbers").Remove()
	}).Html()

	reg_pre := regexp.MustCompile(`(?sU:<figure.*>)(?s:.*?)(?U:</figure>)`)
	str = reg_pre.ReplaceAllStringFunc(str, func(str string) string {
		str = regexp.MustCompile(`(?U:</*figcaption.*>)`).ReplaceAllString(str, "")
		str = regexp.MustCompile(`(?U:</*span.*>)`).ReplaceAllString(str, "")
		str = regexp.MustCompile(`(?U:</*td.*>)`).ReplaceAllString(str, "")
		str = regexp.MustCompile(`(?U:</*tr.*>)`).ReplaceAllString(str, "")
		str = regexp.MustCompile(`(?U:</*table.*>)`).ReplaceAllString(str, "")
		str = regexp.MustCompile(`(?U:</*tbody.*>)`).ReplaceAllString(str, "")
		str = regexp.MustCompile(`(?U:</*figure.*>)`).ReplaceAllString(str, "")
		return  str
	})
	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	return strings.TrimSpace(strings.Replace(title,"\n","",0)), strings.TrimSpace(str)
}

// http://www.ruanyifeng.com
func (this *ServiceContent) getContent30(url string) (string, string) {
	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title,_ := doc.Find("article").Find("#page-title").Html()
	str, _ := doc.Find("article").Find("#main-content").Html()
	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}

// https://tech.meituan.com/ruby_autotest.html
func (this *ServiceContent)getContent31(url string) (string, string) {
	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title,_ := doc.Find("#post_detail").Find(".title").Html()
	str, _ := doc.Find("#post_detail").Find(".article__content").Html()
	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}

//
func (this *ServiceContent)getContent32(url string) (string, string) {
	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title := doc.Find(".content_text").Find(".title1").Text()
	str, _ := doc.Find(".content_banner").Find(".text").Html()
	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}

// https://www.barretlee.com/blog/2018/03/01/%E9%99%AA%E4%BC%B4/
func (this *ServiceContent) getContent33(url string) (string, string) {
	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title := doc.Find(".article").Find(".post-title").Text()
	str, _ := doc.Find(".article").Find(".post-content").Each(func(_ int, tag *goquery.Selection) {
		tag.Find(".shit-spider").Remove()
	}).Html()
	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}

// http://mindhacks.cn/2017/10/17/through-the-maze-11/
func (this *ServiceContent) getContent37(url string) (string, string) {
	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title := doc.Find(".style_breadcrumbs").Find(".container").Find("h1").Text()
	str, _ := doc.Find("article").Find(".entry-content").Html()
	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}

//http://www.520monkey.com/archives/1238
func (this *ServiceContent) getContent39(url string) (string, string) {
	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title := doc.Find("header").Find(".article-title").Text()
	str, _ := doc.Find("article").Each(func(_ int, tag *goquery.Selection) {
		tag.Find("#Addlike").Remove()
		tag.Find(".article-social").Remove()
	}).Html()
	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}

// http://melonteam.com/posts/chu_tan_kotlin_yi_bu_async_await/
func (this *ServiceContent) getContent41(url string) (string,string) {
	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title := doc.Find(".left").Find("h1").Text()
	str, _ := doc.Find("article").Html()

	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}

//http://cdc.tencent.com/2018/04/26/%E6%B3%9B%E5%A8%B1%E4%B9%90%E5%BE%AE%E4%BF%A1%E5%BA%97%E6%94%B9%E7%89%88-%E8%AE%BE%E8%AE%A1%E6%80%BB%E7%BB%93/
func (this *ServiceContent) getContent42(url string) (string, string) {
	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title := doc.Find(".content-title").Find("h3").Text()
	str, _ := doc.Find(".content-text").Html()

	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}


func (this *ServiceContent) getContent44(url string) (string, string) {
	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title := doc.Find("#content").Find(".article_container.row").Find("h1").Text()
	str, _ := doc.Find("#post_content").Html()

	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}

// http://fex.baidu.com/blog/2018/04/fex-weekly-30/
func (this *ServiceContent) getContent45(url string) (string, string) {
	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title := doc.Find("article").Find(".title").Text()
	str, _ := doc.Find("article").Find(".content").Html()

	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}

// http://taobaofed.org/blog/2018/03/12/long-list-in-rax/
func (this *ServiceContent) getContent46(url string) (string, string) {
	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title := doc.Find("article").Find(".article-title").Text()
	str, _ := doc.Find("article").Find(".article-entry").Html()

	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}

// http://mobilefrontier.github.io/articles/weekly-50/
func (this *ServiceContent) getContent48(url string) (string, string) {
	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title := doc.Find("article").Find(".post-title").Text()
	str, _ := doc.Find("article").Find(".post-content").Html()

	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}

// http://blog.sunnyxx.com/2016/08/13/reunderstanding-runtime-1/
func (this *ServiceContent) getContent62(url string) (string, string)  {
	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title := doc.Find("article").Find(".post-title").Text()
	str, _ := doc.Find("article").Find(".post-content").Html()

	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}

// http://blog.raozhizhen.com/post/2016-08-19
func (this *ServiceContent) getContent63(url string) (string, string) {
	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title := doc.Find("article").Find(".title").Text()
	str, _ := doc.Find("article").Find(".post").Html()

	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}

// http://msching.github.io/blog/2016/05/24/audio-in-ios-9/
func (this *ServiceContent) getContent64(url string) (string, string) {
	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title := doc.Find("article").Find(".title").Text()
	str, _ := doc.Find("article").Find(".entry-content").Html()

	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}


// http://coderyi.com/posts/weex3/
func (this *ServiceContent) getContent66(url string) (string, string) {
	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title := doc.Find(".entry-title").Text()
	str, _ := doc.Find(".entry-body").Html()

	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}

// https://casatwy.com/Advance_In_iOS11_Networking.html
func (this *ServiceContent) getContent67(url string) (string, string) {
	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title := doc.Find("article").Find("h1").Text()
	str, _ := doc.Find("article").Find(".articleContent").Html()

	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}

// http://blog.cnbang.net/writting/3565/
func (this *ServiceContent) getContent68(url string) (string, string) {
	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title := doc.Find(".post").Find(".title").Text()
	str, _ := doc.Find(".post").Find(".post_content").Html()

	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}

// http://f2e.souche.com/blog/webpackbian-yi-vuexiang-mu-sheng-cheng-de-dai-ma-tan-suo/
func (this *ServiceContent) getContent70(url string) (string, string) {
	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title := doc.Find("article").Find(".post-title").Text()
	str, _ := doc.Find("article").Find(".post-content").Each(func(_ int, tag *goquery.Selection) {
		tag.Find("p").Each(func(_ int, p *goquery.Selection){
			s, _ := p.Html()
			if  strings.Contains(s, "下一篇") ||strings.Contains(s, "上一篇") {
				p.Remove()
			}
		})
	}).Html()

	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}

//http://blog.codingplayboy.com/2018/04/15/react_native_dev/
func (this *ServiceContent) getContent73(url string) (string, string) {
	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title := doc.Find("article").Find(".entry-title").Text()
	str, _ := doc.Find("article").Find(".entry-content").Each(func(_ int, tag *goquery.Selection) {
		tag.Find("#toc_container").Remove()
		tag.Find(".jiathis_style").Remove()
	}).Html()
	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}

// https://luolei.org/terramaster-d2-310-review/
func (this *ServiceContent) getContent76(url string) (string, string) {
	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title := doc.Find(".post-info-container").Find(".post-page-title ").Text()
	str, _ := doc.Find("article").Find(".post-original").Html()
	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}

// https://lukesign.com/disable-wechat-font-adjust/
func (this *ServiceContent) getContent77(url string) (string, string) {
	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title := doc.Find("article").Find(".post-title").Text()
	str, _ := doc.Find("article").Find(".post-content").Html()
	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}

// http://ued.ctrip.com/?p=5657
func (this *ServiceContent) getContent79(url string) (string, string) {
	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title := doc.Find(".article").Find(".media-heading").Text()
	str, _ := doc.Find(".article").Find(".article-body").Html()
	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}

// https://aotu.io/notes/2018/04/24/jdindex_2017/
func (this *ServiceContent) getContent80(url string) (string, string) {
	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title := doc.Find("article").Find(".post-tit").Text()
	str, _ := doc.Find("article").Find(".post-content").Html()
	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	str = utils.RegStyle(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}

//https://blog.qiniu.com/archives/8728
func (this *ServiceContent) getContent88(url string) (string, string) {
	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title := doc.Find(".blog-content").Find(".entry-title").Text()
	str, _ := doc.Find(".blog-content").Find(".blog-html").Html()
	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	str = utils.RegStyle(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}

func (this *ServiceContent) getContent89(url string) (string, string) {
	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title := doc.Find("article").Find(".article-title").Text()
	str, _ := doc.Find(".blog-content").Find(".blog-html").Html()
	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	str = utils.RegSpan(str)
	str = utils.RegStyle(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}

func (this *ServiceContent) initRegexp(){
	this.RegStyle = regexp.MustCompile(`(?sU:<style[\s\S]*>)([\s\S]*)(?U:</style>)`)
	this.RegScript = regexp.MustCompile(`(?sU:<script[\s\S]*>)([\s\S]*)(?U:</script>)`)
	this.UnLikeyCandidates,_ = regexp.Compile(`banner|breadcrumbs|combx|comment|community|cover-wrap|disqus|extra|foot|header|legends|menu|related|remark|replies|rss|shoutbox|sidebar|skyscraper|social|sponsor|supplemental|ad-break|agegate|pagination|pager|popup|yom-remote|nav`)
}


func (this *ServiceContent) removeAndGetNext(node *html.Node) *html.Node {
	next := this.getNextNode(node, true)
	node.Parent.RemoveChild(node)
	return next
}

func (this *ServiceContent) getNextNode(node *html.Node, ignoreSelfAndKids bool) *html.Node {
	if node == nil {
		return nil
	}
	if !ignoreSelfAndKids && node.FirstChild != nil {
		return node.FirstChild
	}

	if node.NextSibling != nil {
		return node.NextSibling
	}

	for {
		if node == nil {
			break
		}
		if node != nil && node.NextSibling != nil {
			return node.NextSibling
		}
		node = node.Parent
	}
	return nil
}

func (this *ServiceContent) getContent(url string) {
	this.initRegexp()
	//doc, _ := html.Parse(this.getResponse(url).Body)
	doc := this.getHtml(url)
	reg_body := regexp.MustCompile(`(?sU:<body[\s\S]*>)([\s\S]*)(?U:</body>)`)
	body := string(reg_body.Find([]byte(doc)))
	body = this.RegScript.ReplaceAllString(body, "")
	body = this.RegStyle.ReplaceAllString(body,"")
	bodyNode,_ := html.Parse(strings.NewReader(body))

	// 1.保留body部分
	node := bodyNode.FirstChild
	for node != nil {
		matchString := ""
		for _, a := range  node.Attr {
			if a.Key == "class" || a.Key == "id" {
				matchString = fmt.Sprintf("%s %s",matchString, a.Val)
			}
		}
		if this.UnLikeyCandidates.Match([]byte(matchString)) {
			fmt.Println("node:",node.DataAtom.String() ," -- ", node.Data, "--", matchString)

			node = this.removeAndGetNext(node)
			continue
		}
		node = this.getNextNode(node, false)
	}
	fmt.Println("html:", this.renderNode(bodyNode))

}

func (this *ServiceContent)renderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	return buf.String()
}