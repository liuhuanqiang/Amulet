package domain

import (
	"amulet/web/msg"
	"encoding/json"

	"amulet/web/model"
	"github.com/golang/glog"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"net/http"
	"fmt"
	"io/ioutil"
	"amulet/utils"
	"bytes"
)

type Content struct {

}

const (
	Source_Blog = 1
	Source_JianShu = 2
	Source_ZhiHu = 3
)

/**
	获取最新的文章的列表
 */
func (this *Content) GetLatestList(str string) interface{} {
	req := &msg.LatestListReq{}
	json.Unmarshal([]byte(str), req)
	glog.Info("GetLatestList:", str ,"  page:" ,req.Page)

	list := db.GetDB().GetLatestList(req.Page)
	for _,item := range list {
		item.Description = this.getSummary(item.Description)

		if item.Source == Source_ZhiHu && strings.Index(item.Title, "发表了文章") == 0{
			strings.Replace(item.Title, "发表了文章", "", 0)
		}
	}

	ret := &msg.LastetListResp{}
	ret.Current = req.Page
	ret.List = list
	return ret
}

// 获取文章中简介语句
// 获取第一个p的标签, 直到后面不是p的标签结束
func (this *Content) getSummary(article string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(article))
	if err != nil {
		glog.Info("GetSummary error:", err.Error())
	}

	ret := []string{}
	p := doc.Find("p").First()
	if len(p.Text()) > 0 {
		ret = append(ret, p.Text())
	}
	for p.Next().Is("p") {
		p = p.Next()
		if len(p.Text()) > 0 {
			ret = append(ret, p.Text())
		}
	}
	if len(ret) == 0 {
		return article
	}
	str := ""
	for _, v := range ret {
		str += v
	}
	return strings.TrimSpace(str)
}


/**
	根据linkid，获取文章内容
 */
func (this *Content) GetArticle(str string) interface{} {
	req := &msg.ArticleReq{}
	json.Unmarshal([]byte(str), req)
	glog.Info("GetArticle:", str ,"  page:" ,req.Linkid, " source:", req.Source)
	if req.Source == Source_Blog {
		fid, url := db.GetDB().GetLink("tb_blog", req.Linkid)

		resp := &msg.ArticleResp{}
		resp.Url = url
		if fid == 1 {
			resp.Title, resp.Content = this.getContent1(url)
		} else if fid == 2 {
			resp.Title, resp.Content = this.getContent2(url)
		} else if fid == 7 || fid == 10 {
			resp.Title, resp.Content = this.getContent7(url)
		} else if fid == 11 {
			resp.Title, resp.Content = this.getContent11(url)
		} else if fid == 13 {
			resp.Title, resp.Content = this.getContent13(url)
		} else if fid == 14 {
			resp.Title, resp.Content = this.getContent14(url)
		} else if fid == 16 {
			resp.Title, resp.Content = this.getContent16(url)
		} else if fid == 17 {
			resp.Title, resp.Content = this.getContent17(url)
		} else if fid == 19 || fid == 22 || fid == 23 {
			resp.Title, resp.Content = this.getCSDNContent(url)
		} else if fid == 34 {
			resp.Title, resp.Content = this.getContent34(url)
		} else {
			resp.Title, resp.Content = this.getHexoContent(url)
		}
		return resp
	} else if req.Source == Source_ZhiHu {
		_, url := db.GetDB().GetLink("tb_zhihu", req.Linkid)
		doc, _ := goquery.NewDocument(url)

		title, _ := doc.Find(".Post-Title").Html();
		str, _ := doc.Find(".Post-RichText").Html()
		str = utils.RegScript(str);
		str = utils.RegSpan(str)

		resp := &msg.ArticleResp{}
		resp.Title = title
		resp.Content = str
		resp.Url = url
		return resp
	} else if req.Source == Source_JianShu {
		_,url := db.GetDB().GetLink("tb_jianshu", req.Linkid)
		doc, _ := goquery.NewDocument(url)

		title,_ := doc.Find(".title").Html();
		str,_ := doc.Find(".show-content").Html()
		// 过滤掉div
		str = utils.RegDiv(str)
		str = utils.RegAnno(str)

		resp := &msg.ArticleResp{}
		resp.Title = title
		resp.Content = str
		resp.Url = url
		return resp
	} else {
		return nil
	}
}

func (this *Content) getResponse (url string) *http.Response {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, bytes.NewReader([]byte("")));
	req.Header.Set("User-Agent","Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36")
	resp, _ := client.Do(req)
	return resp
}

func (this *Content) getHtml(url string) string {
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

func (this *Content) GetContent(url string) (string, string) {
	return this.getContent24(url)
}

func (this *Content) getHexoContent(url string) (string, string) {
	doc, _ := goquery.NewDocument(url)
	title, _ := doc.Find("article").Find("h1").Html()
	str, _ := doc.Find("article").Find(".entry").Html()
	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	return strings.TrimSpace(strings.Replace(title,"\n","",0)), strings.TrimSpace(str)
}

// http://rango.swoole.com/
func (this *Content) getContent1(url string) (string, string){
	doc, _ := goquery.NewDocument(url)
	title,_ := doc.Find(".post").Find("h2").Html();
	str,_ := doc.Find(".post").Find(".content").Html()

	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	return title, strings.TrimSpace(str)
}

// http://www.laruence.com/
func (this *Content) getContent2(url string) (string, string) {
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

func (this *Content) getContent34(url string) (string, string) {
	doc, _ := goquery.NewDocument(url)
	title, _ := doc.Find("article").Find("h1").Html()
	str, _ := doc.Find("article").Find(".article-entry").Html()
	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	return strings.TrimSpace(strings.Replace(title,"\n","",0)), strings.TrimSpace(str)
}

func (this *Content) getContent7(url string) (string, string) {
	doc, _ := goquery.NewDocument(url)
	title, _ := doc.Find("article").Find("h1").Html()
	str, _ := doc.Find("article").Find(".entry-content").Html()
	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	return strings.TrimSpace(strings.Replace(title,"\n","",0)), strings.TrimSpace(str)
}

func (this *Content) getContent11(url string) (string, string) {
	doc, _ := goquery.NewDocument(url)
	title, _ := doc.Find("h1").Html()
	str, _ := doc.Find(".container.typo").Html()
	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	return strings.TrimSpace(strings.Replace(title,"\n","",0)), strings.TrimSpace(str)
}

func (this *Content) getContent13(url string) (string, string) {
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

func (this *Content) getContent14(url string) (string, string) {
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

func (this *Content) getContent16(url string) (string, string) {
	doc, _ := goquery.NewDocument(url)
	title := doc.Find(".content").Find(".page-header").Find("h1").Text()
	str, _ := doc.Find("#article").Html()
	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}

func (this *Content) getContent17(url string) (string, string) {

	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title,_ := doc.Find(".x-center").Find(".x-content").Find("h3").Html()
	str, _ := doc.Find(".x-center").Find(".x-article-content.x-main-content").Html()
	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}

func (this *Content) getCSDNContent(url string) (string,string){
	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title,_ := doc.Find(".title-article").Html()
	str, _ := doc.Find("article").Find(".article_content").Html()
	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}

func (this *Content) getContent24(url string) (string, string) {
	doc, _ := goquery.NewDocumentFromResponse(this.getResponse(url))
	title,_ := doc.Find(".post-title").Html()
	str, _ := doc.Find("article").Find(".post-content").Html()
	str = utils.RegDiv(str)
	str = utils.RegAnno(str)
	str = utils.RegScript(str)
	return strings.TrimSpace(title), strings.TrimSpace(str)
}