package types

import "amulet/core"

type DocumentResp []*core.Document
func (docs DocumentResp) Len() int {
	return len(docs)
}
func (docs DocumentResp) Swap(i, j int) {
	docs[i], docs[j] = docs[j], docs[i]
}
func (docs DocumentResp) Less(i, j int) bool {
	return docs[i].Frequency > docs[j].Frequency
}