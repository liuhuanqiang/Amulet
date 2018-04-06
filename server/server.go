package server

import (
	"net/http"
	"amulet/segment"
	"amulet/core"
	"encoding/json"
	"github.com/golang/glog"
	"time"
)

type Server struct {
	Seg   *segment.Segment
	TextRank  *core.TextRank
	StopToken  *core.StopToken
}

const (
	Success = 1
	Failed  = 0
)

type Resp struct {
	Code    int  `json:"code"`
	Msg     string  `json:"msg"`
	Data    interface{} 	`json:"data"`
}

// 启动一个server
func (this *Server) Start() {
	this.Seg = &segment.Segment{}
	this.Seg.Init()
	this.TextRank = &core.TextRank{}
	this.StopToken = &core.StopToken{}
	this.StopToken.Init()


	http.HandleFunc("/getRankList", this.GetRankList)
	http.HandleFunc("/cut", this.Cut)
	http.HandleFunc("/getSummary", this.GetSummary)

	glog.Info("端口8081: 启动成功")
	http.ListenAndServe(":8081", nil)
}


// 分词
func (this *Server) Cut(w http.ResponseWriter, r *http.Request) {

}

// 获取ranklist
func (this *Server) GetRankList(w http.ResponseWriter, r *http.Request) {
	// 获取文本
	startTime := time.Now()
	text := r.FormValue("text")
	terms := this.Seg.Cut([]byte(text))
	list := []string{}
	for _, term := range terms {
		if !this.StopToken.IsStopToken(term) {
			list = append(list, term)
		}
	}
	ret := &Resp{}
	ret.Data = this.TextRank.GetRankList(list)
	ret.Code = Success
	ret.Msg = "成功"
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	glog.Info("text:", text, "   cost:", time.Since(startTime))
	this.RenderJson(w,ret)
}

// 获取文章的简介
func (this *Server) GetSummary(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	article := r.FormValue("article")
	ret := &Resp{}
	ret.Code = Success
	ret.Data = this.TextRank.GetSummary(article)
	ret.Msg = "成功"
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	glog.Info("text:", ret.Data, "   cost:", time.Since(startTime))
	this.RenderJson(w,ret)
}

func (this *Server)RenderJson(w http.ResponseWriter, v interface{}) {
	bs, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(bs)
}