package domain

import (
	"amulet/web/msg"
	"encoding/json"

	"amulet/web/model"
	"github.com/golang/glog"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"amulet/utils"

	"amulet/web/service"
	"time"
	"net/http"
	"io/ioutil"
	"bytes"
)

type Content struct {
	Service  *service.ServiceContent
	Readability *service.Readability
	Mss  *service.MaxSubSegment
}

const (
	Source_Blog = 1

	Source_ZhiHu = 3
	Source_JianShu = 4
)

func (this *Content) Init() {
	this.Service = &service.ServiceContent{}
	this.Readability = &service.Readability{}
	this.Mss = &service.MaxSubSegment{}
}
/**
	获取最新的文章的列表
 */
func (this *Content) GetLatestList(str string) interface{} {
	req := &msg.LatestListReq{}
	json.Unmarshal([]byte(str), req)
	glog.Info("GetLatestList:", str ,"  page:" ,req.Page)

	list := db.GetDB().GetLatestList(req.Page, req.Timestamp)
	for _,item := range list {
		item.Description = utils.SubString(utils.ReplaceEscapeStr(this.getSummary(item.Description)),0,200)
		item.Title = utils.ReplaceEscapeStr(item.Title)
		glog.Info(utils.ReplaceEscapeStr(this.getSummary(item.Description)))
		if item.Source == Source_ZhiHu && strings.Index(item.Title, "发表了文章") == 0{
			strings.Replace(item.Title, "发表了文章", "", 0)
		}
	}

	ret := &msg.LastetListResp{}
	ret.Current = req.Page
	ret.List = list
	ret.Timestamp = int(time.Now().Unix())
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
	glog.Info("GetArticle:", str ,"  linkid:" ,req.Linkid, " source:", req.Source)
	if req.Source == Source_Blog {
		_, url := db.GetDB().GetLink("tb_blog", req.Linkid)

		resp := &msg.ArticleResp{}
		resp.Url = url
		//resp.Title, resp.Content = this.Service.GetContent(fid, url)
		resp.Title = "test"
		glog.Info("url:", url)
		resp.Content = this.Mss.GetContent(url)
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
		glog.Info("url:", url)
		doc, _ := goquery.NewDocumentFromReader(strings.NewReader(this.getHtml(url)))

		title,_ := doc.Find(".title").Html();
		str,_ := doc.Find(".show-content").Html()
		// 过滤掉div
		//str = utils.RegDiv(str)
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

func (this *Content) GetArticleByUrl(url string) interface{} {
	resp := &msg.ArticleResp{}
	resp.Url = url
	resp.Title = "test"
	resp.Content = this.Mss.GetContent(url)
	return resp
}

func (this *Content) getHtml(url string) string {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, bytes.NewReader([]byte("")));
	req.Header.Set("User-Agent","Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.139 Safari/537.36")
	resp, _ := client.Do(req)
	str,_ := ioutil.ReadAll(resp.Body)
	defer  resp.Body.Close()
	return string(str)
}