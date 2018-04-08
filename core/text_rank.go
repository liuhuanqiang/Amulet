package core

import (
	"sort"
	"github.com/golang/glog"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type TextRank struct {

}

const Default_Num  =  4

type TextScore struct {
	Text   string
	Score  float32
}

type TextScoreList []*TextScore
func (p TextScoreList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p TextScoreList) Len() int           { return len(p) }
func (p TextScoreList) Less(i, j int) bool { return p[i].Score < p[j].Score }
// 获取RankList
func (this *TextRank) GetRankList(terms []string) []string {
	words := make(map[string][]string)
	for i := 0; i < len(terms); i++ {
		s := 0
		if i - 5 > 0 {
			s = i - 5
		}
		e := i + 5
		if i + 5 > len(terms) {
			e = len(terms)
		}
		for j := s; j < e; j++ {
			if i != j && !this.isContain(words[terms[i]],terms[j]){
				words[terms[i]] = append(words[terms[i]], terms[j])
			}
		}
	}
	// 开始投票
	scores := make(map[string]float32)
	for i:= 0; i < 200; i++ {
		m := make(map[string]float32)
		for word, list := range words {
			m[word] = 1 - 0.85
			for _, w := range list {
				m[word] += 0.85 / float32(len(words[w])) * scores[w]
			}

		}
		scores = m
	}
	glog.Info("scores:", scores)
	// 排序 取前几个单词
	var scoresList TextScoreList
	for k, v := range scores {
		tmp := &TextScore{}
		tmp.Text = k
		tmp.Score = v
		scoresList = append(scoresList, tmp)
	}
	n := Default_Num
	sort.Sort(sort.Reverse(scoresList))
	if len(scores) < Default_Num {
		n = len(scores)
	}
	ret := []string{}
	for i:= 0; i < n; i++ {
		ret = append(ret,scoresList[i].Text)
	}
	return ret
}

// 判断是否已经在这个列表中了
func (this *TextRank) isContain(terms []string, term string) bool {
	for _, v := range terms {
		if v == term {
			return true
		}
	}
	return false
}

// 获取文章中简介语句
// 获取第一个p的标签, 直到后面不是p的标签结束
func (this *TextRank) GetSummary(article string) string {
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