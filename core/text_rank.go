package core


type TextRank struct {

}


// 获取RankList
func (this *TextRank) GetRankList(terms []string) map[string]float32 {
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
	return scores
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
