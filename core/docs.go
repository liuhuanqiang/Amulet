package core

// 作用:用于管理Doc, 索引是docId
type DocsManager struct {
	Docs  map[uint64]*Document
	TotalLen  uint64   // 总长度
	CacheLen  uint64   // 缓存的长度
}

func (this *DocsManager) Init() {
	this.Docs = make(map[uint64]*Document)
}

func (this *DocsManager) Add(document *Document) {
	this.Docs[document.DocId] = document
	this.TotalLen++
}


func (this *DocsManager) Get(docId uint64) *Document {
	doc, ok := this.Docs[docId]
	if ok {
		return doc
	} else {
		return nil
	}
}

func (this *DocsManager) Index() uint64 {
	return uint64(len(this.Docs))
}


