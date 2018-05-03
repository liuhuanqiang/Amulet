package main

import (
	"net/http"
	"amulet/web/domain"
	"strings"
	"reflect"
	"encoding/json"
	"github.com/golang/glog"
	"amulet/web/msg"
	"time"
)


type Router struct {
	RouterMap  map[string] Domain
}

type Domain interface {
	Init()
}

func (this *Router) StartHttp() {
	this.RouterMap = make(map[string]Domain)
	this.RouterMap["content"] = &domain.Content{}
	this.RouterMap["content"].Init()

	http.HandleFunc("/index/",this.route)
	http.ListenAndServe(":8888", nil)
}

func (this *Router)route(w http.ResponseWriter,r *http.Request) {
	startTime := time.Now()
	s := r.FormValue("s")
	param := r.FormValue("param")
	str := strings.Split(s, ".")
	if len(str) != 2 {
		// 参数不合理
		glog.Info("==> ", "s:", s, "  param:",param, " cost:", time.Since(startTime))
		this.write(w, nil, msg.HttpStatus_InvalidParam, "异常参数")
		return
	}
	domain, exist := this.RouterMap[str[0]]
	if !exist {
		// 不存在Domain
		glog.Info("==> ", "s:", s, "  param:",param, " cost:", time.Since(startTime))
		this.write(w, nil, msg.HttpStatus_InvalidDomain, "异常Domain")
		return
	}
	mtV := reflect.ValueOf(domain)
	method := mtV.MethodByName(str[1])
	if method.String() == "<invalid Value>" {
		// 不存在方法
		glog.Info("==> ", "s:", s, "  param:",param, " cost:", time.Since(startTime))
		this.write(w, nil, msg.HttpStatus_InvalidMethod, "异常方法")
		return
	}
	ret :=  method.Call([]reflect.Value{reflect.ValueOf(param)})

	this.write(w,ret[0].Interface(),msg.HttpStatus_Success, "成功")
	glog.Info("==> ", "s:", s, "  param:",param, " cost:", time.Since(startTime))
}

func (this *Router) write(w http.ResponseWriter,data interface{}, code int, info string) {
	resp := &msg.Resp{}
	resp.Data = data
	resp.Code = code
	resp.Msg = info
	bs, err := json.Marshal(resp)
	if err != nil {
		resp.Msg = err.Error()
		resp.Code = msg.HttpStatus_Exception
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(bs)
}