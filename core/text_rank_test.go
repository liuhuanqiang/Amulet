package core

import "testing"

func TestTextRank_GetRankList(t *testing.T) {
	terms := []string{"程序员", "英文", "程序", "开发", "维护", "专业", "人员", "程序员", "分为", "程序", "设计", "人员", "程序", "编码", "人员", "界限", "特别", "中国", "软件", "人员", "分为", "程序员", "高级", "程序员", "系统", "分析员", "项目", "经理"}
	Rank := &TextRank{}
	Rank.GetRankList(terms)
}
