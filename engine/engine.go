package engine

import (
	"amulet/core"
	"github.com/huichen/sego"
	"amulet/db"
	"net/http"
	"amulet/utils"
	"encoding/json"
	"strings"
	"github.com/golang/glog"
	"amulet/types"
	"sort"
	"time"
	"amulet/config"
)

type Engine struct {
	Indexer 	*core.Indexer
	StopToken 	*core.StopToken
	Sego            *sego.Segmenter
	DB		*db.MysqlDB
	DocsManager     *core.DocsManager
	TextRank        *core.TextRank
}

func (this *Engine) Init() {
	startTime := time.Now()
	config.LoadConfig("config.json")
	this.Indexer = &core.Indexer{}
	this.StopToken = &core.StopToken{}
	this.Sego = &sego.Segmenter{}
	this.DB = &db.MysqlDB{}
	this.DocsManager = &core.DocsManager{}
	this.TextRank = &core.TextRank{}

	this.Indexer.Init()
	this.StopToken.Init()
	this.Sego.LoadDictionary("data/dictionary.txt,data/add.txt")
	this.DB.Init()
	this.DocsManager.Init()
	this.CreateIndex()
	glog.Info("创建索引结束... cose time:", time.Since(startTime))
	this.StartHttpServer()
}


// ToDo: 启动一个http server
func (this *Engine) StartHttpServer() {
	http.HandleFunc("/search",this.Search)
	http.ListenAndServe(":8888", nil)
}

// ToDo:
func (this *Engine) StopHttpServer() {

}

// 创建索引
func (this *Engine) CreateIndex() {
	blogs := this.DB.GetBlog()
	for _, blog := range blogs {
		document := &core.Document{}
		document.DocId = this.DocsManager.Index()
		document.Table = "blog"
		document.Content = blog.Description
		document.Id = blog.Id
		//去除空格，转义，标签
		segments := this.Sego.Segment([]byte(utils.FilterInvalidChar(blog.Description)))
		// 分词
		slice := sego.SegmentsToSlice(segments, true)
		var terms []string
		for _, item := range slice {
			item = strings.TrimSpace(item)
			if !this.StopToken.IsStopToken(item) {
				terms = append(terms, item)
			}
		}
		// 计算TextRank
		textRankScore := this.TextRank.GetRankList(terms)
		for _, item := range terms {
			// 去掉两边的空格
			keyword := &core.Keyword{}
			keyword.Text = item
			keyword.Frequency+= 1.0/float32(len(terms))
			keyword.TextRank = textRankScore[item]
			document.Keyword = append(document.Keyword, keyword)
		}
		this.DocsManager.Add(document)
		this.Indexer.AddDocument(document)
	}
	zhihus := this.DB.GetZhiHu()
	for _, zhihu := range zhihus {
		document := &core.Document{}
		document.DocId = this.DocsManager.Index()
		document.Table = "zhihu"
		document.Content = zhihu.Description
		document.Id = zhihu.Id
		//去除空格，转义，标签
		segments := this.Sego.Segment([]byte(utils.FilterInvalidChar(zhihu.Description)))
		// 分词
		slice := sego.SegmentsToSlice(segments, true)
		var terms []string
		for _, item := range slice {
			item = strings.TrimSpace(item)
			if !this.StopToken.IsStopToken(item) {
				terms = append(terms, item)
			}
		}
		for _, item := range terms {
			// 去掉两边的空格
			keyword := &core.Keyword{}
			keyword.Text = item
			keyword.Frequency+= 1.0/float32(len(terms))
			document.Keyword = append(document.Keyword, keyword)
		}
		this.DocsManager.Add(document)
		this.Indexer.AddDocument(document)
	}
	jianshus := this.DB.GetJianShu()
	for _, jianshu := range jianshus {
		document := &core.Document{}
		document.DocId = this.DocsManager.Index()
		document.Table = "jianshu"
		document.Content = jianshu.Description
		document.Id = jianshu.Id
		//去除空格，转义，标签
		segments := this.Sego.Segment([]byte(utils.FilterInvalidChar(jianshu.Description)))
		// 分词
		slice := sego.SegmentsToSlice(segments, true)
		var terms []string
		for _, item := range slice {
			item = strings.TrimSpace(item)
			if !this.StopToken.IsStopToken(item) {
				terms = append(terms, item)
			}
		}
		for _, item := range terms {
			// 去掉两边的空格
			keyword := &core.Keyword{}
			keyword.Text = item
			keyword.Frequency+= 1.0/float32(len(terms))
			document.Keyword = append(document.Keyword, keyword)
		}
		this.DocsManager.Add(document)
		this.Indexer.AddDocument(document)
	}


}

// ToDo:
func (this *Engine) Search(w http.ResponseWriter,r *http.Request) {
	keyword := r.FormValue("keyword")
	w.Header().Set("content-type", "application/json")
	glog.Info("keyword:", keyword)
	indices := this.Indexer.Search(keyword)
	if indices == nil {
		// 暂无结果
		glog.Info("暂无结果")
	} else {
		var ret types.DocumentResp
		for i, docid := range indices.DocIds {
			doc := this.DocsManager.Get(docid)
			doc.Frequency = indices.Frequencys[i]
			doc.TextRank = indices.TextRank[i]
			if doc != nil {
				ret = append(ret, doc)
			}
		}
		sort.Sort(ret)
		bs, err := json.Marshal(ret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(bs)
	}

}
