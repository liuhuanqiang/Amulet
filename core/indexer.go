package core

import "github.com/golang/glog"

// 索引器
type Indexer struct {
	TokenLen   int
	DocLen     int
	Map        map[string]*KeywordIndices
}

type Document struct {
	DocId 		uint64    	`json:"doc_id"`
	Table           string		`json:"table"`
	Id              int		`json:"id"`
	Keyword         []*Keyword	`json:"-"`
	Content         string		`json:"content"`
	Frequency	float32		`json:"frequency"`  // 这个是搜索查询用到的，后面需要去掉
	TextRank 	float32		`json:"text_rank"`  // 这个是搜索查询用到的，后面需要去掉
}

type Keyword struct {
	Text 		string	 	`json:"text"`
	Frequency  	float32		`json:"frequency"`  // 每个关键字，在这个文章中的词频
	TextRank	float32		`json:"text_rank"`  // 每个关键字，在这个文章中的TextRank
}

// 反向索引表的一行，收集了一个关键字出现的所有文档，按照DocId从小到大排序。
// 已经关键字在该文档中出现的频率。后面搜索会用到
type KeywordIndices struct {
	DocIds      []uint64
	Frequencys  []float32
	TextRank    []float32
}

func (this *Indexer) Init() {
	glog.Info("Indexer初始化...")
	this.Map = make(map[string]*KeywordIndices)
}

// 从KeywordIndices中获取第i个文档的docId
func (this *Indexer) getDocId(ti *KeywordIndices, i int) uint64 {
	return ti.DocIds[i-1]
}

func (this *Indexer) getIndexLength(ti *KeywordIndices) int {
	return len(ti.DocIds)
}

func (this *Indexer) AddDocument(document *Document) {
	// 先加入keyword
	for _, word := range document.Keyword {
		indices, ok := this.Map[word.Text]
		//glog.Info("word:", word, "   ok:", ok)
		if ok {
			// 如果关键字存在, 查找应该插入的位置
			pos, exist := this.searchIndex(indices, 1, this.getIndexLength(indices), document.DocId)
			if !exist {
				indices.DocIds = append(indices.DocIds, 0)
				copy(indices.DocIds[pos + 1:], indices.DocIds[pos:])
				indices.DocIds[pos] = document.DocId

				indices.Frequencys = append(indices.Frequencys, 0)
				copy(indices.Frequencys[pos + 1:], indices.Frequencys[pos:])
				indices.Frequencys[pos] = word.Frequency

				indices.TextRank = append(indices.TextRank, 0)
				copy(indices.TextRank[pos + 1:], indices.TextRank[pos:])
				indices.TextRank[pos] = word.TextRank

			}
		} else {
			// 不存在的话，直接添加
			ti := KeywordIndices{}
			ti.DocIds = []uint64{document.DocId}
			ti.Frequencys = []float32{word.Frequency}
			ti.TextRank = []float32{word.TextRank}
			this.Map[word.Text] = &ti
		}
	}
}

//二分法查找
func (this *Indexer) searchIndex(indices *KeywordIndices, start int, end int, docId uint64) (int,bool) {

	if start == this.getIndexLength(indices) {
		return start, false
	}
	if docId < this.getDocId(indices, start) {
		return start, false
	} else if docId == this.getDocId(indices, start) {
		return start, true
	}
	if docId > this.getDocId(indices, end) {
		return end, false
	} else if docId == this.getDocId(indices, end) {
		return end, true
	}

	for end > start {
		middle := (end + start)/2
		if docId == this.getDocId(indices,middle) {
			return middle,true
		} else if docId > this.getDocId(indices,middle) {
			start = middle
		} else {
			end = middle
		}
	}
	return end,false
}

// 所有
func (this *Indexer) Search(key string) *KeywordIndices {
	return this.Map[key]
}