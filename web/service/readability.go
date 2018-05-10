package service

import (
	"io"
	"golang.org/x/net/html"
	"bytes"
	"regexp"
	"net/http"
	"io/ioutil"
	"fmt"
	"strings"
	"golang.org/x/net/html/atom"
)

type Readability struct {
	UnLikelyCandidates *regexp.Regexp
}

func (this *Readability) initRegexp(){
	this.UnLikelyCandidates,_ = regexp.Compile(`banner|breadcrumbs|combx|comment|community|cover-wrap|disqus|extra|foot|header|legends|menu|related|remark|replies|rss|shoutbox|sidebar|skyscraper|social|sponsor|supplemental|ad-break|agegate|pagination|pager|popup|yom-remote|nav`)
}

func (this *Readability) GetContent(url string) {

	this.initRegexp()
	bodyString := this.getBody(this.getHtml(url))
	//fmt.Println("body:", bodyString)
	if bodyString == "" {
		fmt.Println("not exist body tag")
	}

	bodyNode, _ := html.Parse(strings.NewReader(bodyString))
	// 1.保留body部分
	node := bodyNode.FirstChild
	for node != nil {
		matchString := ""
		for _, a := range node.Attr {
			if a.Key == "class" || a.Key == "id" {
				matchString = fmt.Sprintf("%s %s", matchString, a.Val)
			}
		}
		// 第一步移除不可能是content的标签
		if this.UnLikelyCandidates.Match([]byte(matchString)) {
			fmt.Println("[unlikely]:", node.DataAtom.String(), " -- ", node.Data, "--", matchString)
			node = this.removeAndGetNext(node)
			continue
		}
		node = this.getNextNode(node, false)
	}

	node = bodyNode.FirstChild
	for node != nil {
		// 第二步 去掉 div, section, header, h1, h2, h3, h4, h5, h6中内容为空的标签
		if node.DataAtom == atom.Div || node.DataAtom == atom.Section || node.DataAtom == atom.Header || node.DataAtom == atom.H1 ||
			node.DataAtom == atom.H2 || node.DataAtom == atom.H3 || node.DataAtom == atom.H4 || node.DataAtom == atom.H5 || node.DataAtom == atom.H6 {
			if this.isElementWithoutContent(node) {
				fmt.Println("[empty]:",this.renderNode(node))
				node = this.removeAndGetNext(node)
				continue
			}

		}
		node = this.getNextNode(node, false)
	}
	fmt.Println("html:", this.renderNode(bodyNode))
	node = bodyNode.FirstChild
	for node != nil {
		// 第三步 如果标签div中只有一个children, 则去掉这个div。 如果只有文字，就用p替换
		if node.DataAtom == atom.Div {
			if this.hasSinglePInsideElement(node) {
				// 如果里面只有一个p, 则去掉这个div
				fmt.Println("[p]", node.Attr)
				child := this.getNodeChildrens(node)[0]
				newNode := &html.Node {
					Type: child.Type,
					DataAtom:child.DataAtom,
					Data: child.Data,
					Namespace:child.Namespace,
					Attr:child.Attr,
					FirstChild: child.FirstChild,
					LastChild:nil,
				}
				node.Parent.AppendChild(newNode)
				node.Parent.RemoveChild(node)
				node = this.getNextNode(newNode.Parent,false)
			} else if this.hasChildBlockElement(node) {
				fmt.Println("[block]", node.Attr)
				child := this.getNodeChildrens(node)[0]
				newNode := &html.Node {
					Type: child.Type,
					DataAtom:child.DataAtom,
					Data: child.Data,
					Namespace:child.Namespace,
					Attr:child.Attr,
					FirstChild: child.FirstChild,
					LastChild:nil,
				}
				node.Parent.AppendChild(newNode)
				node.Parent.RemoveChild(node)
				node = this.getNextNode(newNode.Parent,false)
				continue
			}
		}
		node = this.getNextNode(node, false)
	}
	fmt.Println("[html]:", this.renderNode(bodyNode))
}

func (this *Readability) getHtml(url string) string {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, bytes.NewReader([]byte("")));
	req.Header.Set("User-Agent","Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36")
	resp, _ := client.Do(req)
	str,_ := ioutil.ReadAll(resp.Body)
	defer  resp.Body.Close()
	return string(str)
}

//
func (this *Readability) getBody(str string) string {
	reg_body := regexp.MustCompile(`(?sU:<body[\s\S]*>)([\s\S]*)(?sU:</body>)`)
	body := reg_body.FindString(str)

	// 过滤<script> </script>
	reg_script := regexp.MustCompile(`(?sU:<script[\s\S]*>)([\s\S]*)(?sU:</script>)`)
	body = reg_script.ReplaceAllString(body, "")

	// 过滤<style> </style>
	reg_style := regexp.MustCompile(`(?sU:<style[\s\S]*>)([\s\S]*)(?sU:</style>)`)
	body = reg_style.ReplaceAllString(body,"")

	// 过滤注释
	reg_anno := regexp.MustCompile(`(?sU:<!).*?(?s:>)`)
	body = reg_anno.ReplaceAllString(body,"")

	return body
}

func (this *Readability) renderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	reg := regexp.MustCompile(`(?sU:</*html.*>)|(?sU:</*head.*>)|(?sU:</*body.*>)`)
	return reg.ReplaceAllString(buf.String(), "")
	//return buf.String()
}

func (this *Readability) getNodeContent(n *html.Node) string {
	var buf bytes.Buffer
	for c:= n.FirstChild; c != nil; c = c.NextSibling {
		buf.WriteString(c.Data)
	}
	return buf.String()
}

func (this *Readability) getNodeChildrens(n *html.Node) []*html.Node {
	childs := []*html.Node{};
	for c:= n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type != html.ErrorNode {
			if c.Type == html.TextNode && strings.TrimSpace(strings.Replace(c.Data,"\n","",0)) == ""{
				continue
			} else {
				childs = append(childs, c)
			}

		}
	}
	return childs
}

func (this *Readability) removeAndGetNext(node *html.Node) *html.Node {
	next := this.getNextNode(node, true)
	node.Parent.RemoveChild(node)
	return next
}

func (this *Readability) getNextNode(node *html.Node, ignoreSelfAndKids bool) *html.Node {
	if node == nil {
		return nil
	}
	if !ignoreSelfAndKids && node.FirstChild != nil {
		return node.FirstChild
	}

	if node.NextSibling != nil {
		return node.NextSibling
	}

	for {
		if node == nil {
			break
		}
		if node != nil && node.NextSibling != nil {
			return node.NextSibling
		}
		node = node.Parent
	}
	return nil
}

func (this *Readability) isElementWithoutContent(node *html.Node) bool {
	if strings.TrimSpace(strings.Replace(this.getNodeContent(node),"\n","",0)) == "" {
		return true
	} else {
		return false
	}
}

func (this *Readability) hasSinglePInsideElement(node *html.Node) bool {
	childrens := this.getNodeChildrens(node)
	if len(childrens) == 1 && childrens[0].DataAtom == atom.P {
		return true
	}
	return false
}

func (this *Readability) hasChildBlockElement(node *html.Node) bool {
	childrens := this.getNodeChildrens(node)
	if len(childrens) == 1 && childrens[0].DataAtom == atom.Div {
		return true
	}
	return false
}