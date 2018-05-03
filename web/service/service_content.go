package service

import (
	"github.com/PuerkitoBio/goquery"
	"strings"
	"net/http"
	"fmt"
	"io/ioutil"
	"regexp"
	"github.com/golang/glog"
	"bytes"
	"amulet/utils"
)

type ServiceContent struct {

}


func (this *ServiceContent) GetContent(fid int, url string) (string, string) {
	title:= ""
	content := ""
	if fid == 1 {
		title, content = this.getContent1(url)
	} else if fid == 2 {
		title, content = this.getContent2(url)
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
	} else if fid == 50 {
		title, content = this.getJacmanContent(url)
	} else {
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
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("http get error.")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("http read error")
	}
	return string(body)
}


func (this *ServiceContent) getHexoContent(url string) (string, string) {
	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title, _ := doc.Find("article").Find("post-title").Html()
	str, _ := doc.Find("article").Find(".post-body").Html()
	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
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

	return title, strings.TrimSpace(str)
}

func (this *ServiceContent) getContent34(url string) (string, string) {
	doc, _ := goquery.NewDocument(url)
	title, _ := doc.Find("article").Find("h1").Html()
	str, _ := doc.Find("article").Find(".article-entry").Html()
	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	return strings.TrimSpace(strings.Replace(title,"\n","",0)), strings.TrimSpace(str)
}

func (this *ServiceContent) getContent7(url string) (string, string) {
	doc, _ := goquery.NewDocument(url)
	title, _ := doc.Find("article").Find("h1").Html()
	str, _ := doc.Find("article").Find(".entry-content").Html()
	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	return strings.TrimSpace(strings.Replace(title,"\n","",0)), strings.TrimSpace(str)
}

func (this *ServiceContent) getContent11(url string) (string, string) {
	doc, _ := goquery.NewDocument(url)
	title, _ := doc.Find("h1").Html()
	str, _ := doc.Find(".container.typo").Html()
	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
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

func (this *ServiceContent) getCSDNContent(url string) (string,string){
	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title,_ := doc.Find(".title-article").Html()
	str, _ := doc.Find("article").Find(".article_content").Html()
	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}

func (this *ServiceContent) getContent24(url string) (string, string) {
	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title,_ := doc.Find(".post-title").Html()
	str, _ := doc.Find("article").Find(".post-content").Html()
	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}

func (this *ServiceContent) getContent27(url string) (string, string) {
	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title,_ := doc.Find("article").Find(".inner").Find(".post-title").Html()
	str, _ := doc.Find("article").Find(".post-content.inner").Html()
	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
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