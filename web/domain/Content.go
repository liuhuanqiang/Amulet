package domain

import (
	"amulet/web/msg"
)

type Content struct {

}

func (this *Content) GetLatestList(str string) interface{} {
	ret := &msg.LatestResp{}
	ret.Current = 1

	// db操作

	return ret
}