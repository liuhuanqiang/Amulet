package service

import (
	"io/ioutil"
	"regexp"
	"net/http"
	"bytes"
	"golang.org/x/net/html"
	"strings"
	"fmt"
	"golang.org/x/net/html/atom"
	"math"
	"io"
)

type MaxSubSegment struct {
	Possible *regexp.Regexp
	UnLikelyCandidates *regexp.Regexp
	Negative *regexp.Regexp
	Positive *regexp.Regexp
}

type ChildrenNode struct {
	Node *html.Node
	Level int
}

func (this *MaxSubSegment) init() {
	this.Possible,_ = regexp.Compile(`article|body|column|main|shadow|entry|content|article|post|container`)
	this.Negative, _ = regexp.Compile(`hidden|^hid$| hid$| hid |^hid |banner|combx|comment|com-|contact|header|legends|menu|related|remark|replies|foot|footer|footnote|masthead|media|meta|outbrain|promo|related|scroll|share|shoutbox|sidebar|skyscraper|sponsor|shopping|tags|tool|widget|sponsor|social|pagination|popup|nav|browsehappy|advertise|jiathis|bdshare|copyright|profile`)
	this.UnLikelyCandidates,_ = regexp.Compile(`banner|breadcrumbs|combx|comment|community|cover-wrap|disqus|extra|foot|header|legends|menu|related|remark|replies|rss|shoutbox|sidebar|skyscraper|social|sponsor|supplemental|ad-break|agegate|pagination|pager|popup|yom-remote|nav|browsehappy|modal|toc|advertise|jiathis|bdshare|copyright|adsbygoogle|respond|spread|msgbar|arrow|back.+?top`)
	this.Positive, _ = regexp.Compile(`article|body|content|entry|hentry|h-entry|main|page|pagination|post|text|blog|story`)
}

func (this *MaxSubSegment) GetContent(url string) string {
	this.init()
	htmlString := this.getHtml(url)
	bodyString := this.getBody(htmlString)
	bodyNode,_ := html.Parse(strings.NewReader(bodyString))

	ScoreMap := make(map[*html.Node]float64);
	childrens := this.getChildrens(bodyNode,false)
	fmt.Println("[children]",  len(childrens))
	//第一步  先把大概的内容筛选出来
	for _,child := range childrens {
		if child.Node.DataAtom == atom.Div || child.Node.DataAtom == atom.Section || child.Node.DataAtom == atom.Article {
			matchString := ""
			for _, a := range child.Node.Attr {
				if a.Key == "class" || a.Key == "id" {
					matchString = fmt.Sprintf("%s %s", matchString, a.Val)
				}
			}
			if this.Possible.Match([]byte(matchString)) {
				// 开始打分
				if this.Negative.Match([]byte(matchString)) {
					ScoreMap[child.Node] -= float64(10)
				}

				if this.Negative.Match([]byte(matchString)) {
					ScoreMap[child.Node] += float64(20)
				}

				childs := this.getChildrens(child.Node,false)
				fmt.Println("[parent]", " level:", child.Level, " node:",child.Node.DataAtom.String(), " attr:", child.Node.Attr)
				for _, c := range childs {
					//fmt.Println("[child]", " level:", c.Level, " node:",c.Node.DataAtom.String(), " attr:", c.Node.Attr, " data:",c.Node.Data)
					divider := 1
					if c.Level == 1 {
						divider = 1
					} else if c.Level == 2 {
						divider = 2
					} else {
						divider = 3 * c.Level
					}
					if c.Node.Type == html.ElementNode {
						matchString := ""
						for _, a := range c.Node.Attr {
							if a.Key == "class" || a.Key == "id" {
								matchString = fmt.Sprintf("%s %s", matchString, a.Val)
							}
						}
						if this.Negative.Match([]byte(matchString)) {
							fmt.Println("[Negative]", " level:", c.Level, " node:",c.Node.DataAtom.String(), " attr:", c.Node.Attr, " data:",c.Node.Data, float64(10)/float64(divider))
							ScoreMap[child.Node] -= float64(10)/float64(divider)
						}

						if this.Positive.Match([]byte(matchString)) {
							fmt.Println("[Positive]", " level:", c.Level, " node:",c.Node.DataAtom.String(), " attr:", c.Node.Attr, " data:",c.Node.Data, float64(10)/float64(divider))
							ScoreMap[child.Node] += float64(20)/float64(divider)
						}

						if c.Node.DataAtom == atom.Span || c.Node.DataAtom == atom.Header || c.Node.DataAtom == atom.H1 || c.Node.DataAtom == atom.Time{
							fmt.Println("[Tag]", " level:", c.Level, " node:",c.Node.DataAtom.String(), " attr:", c.Node.Attr, " data:",c.Node.Data, float64(10)/float64(divider))
							ScoreMap[child.Node] -= float64(10)/float64(divider)
						}

						if c.Node.DataAtom == atom.Pre || c.Node.DataAtom == atom.Blockquote {
							// 一般是代码块
							fmt.Println("[Code]", " level:", c.Level, " node:",c.Node.DataAtom.String(), " attr:", c.Node.Attr, " data:",c.Node.Data, float64(10)/float64(divider))
							ScoreMap[child.Node] += float64(20)/float64(divider)
						}

					} else if c.Node.Type == html.TextNode {
						str := strings.TrimSpace(strings.Replace(c.Node.Data,"\n","",0))
						ScoreMap[child.Node] += float64(math.Min(math.Floor(float64(len(str))/50),float64(3)))/float64(divider)
						ScoreMap[child.Node] += float64(len(strings.Split(str, ","))*4)/float64(divider)
						ScoreMap[child.Node] += float64(len(strings.Split(str, "。"))*3)/float64(divider)

					}
				}
			}
		}
	}
	max := float64(0)
	var maxNode *html.Node
	for child, score := range ScoreMap {
		fmt.Println("[Text]",  child.DataAtom.String(), "  attr:", child.Attr, " score:", score)
		if score >= max {
			maxNode = child
			max = score
		}
	}
	fmt.Println("[MaxNode]",  maxNode.DataAtom.String(), "  attr:", maxNode.Attr, " score:", max)
	//第二步. 获取到最高的元素之后，开始清除没用的元素
	//maxNode.Parent = nil
	node := maxNode.FirstChild
	for node != nil {
		matchString := ""
		for _, a := range node.Attr {
			if a.Key == "class" || a.Key == "id" {
				matchString = fmt.Sprintf("%s %s", matchString, a.Val)
			}
		}
		// 第一步移除不可能是content的标签
		//fmt.Println("[matchString]:", matchString)
		if this.UnLikelyCandidates.Match([]byte(matchString)) && node.DataAtom!= atom.A {
			fmt.Println("[unlikely]:", node.DataAtom.String(), " -- ", node.Data, "--", matchString)
			node = this.removeAndGetNext(node)
			continue
		}

		node = this.getNextNode(node, false)
	}
	fmt.Println(this.renderNode(maxNode))
	fmt.Println("[MaxNode]",  maxNode.DataAtom.String(), "  attr:", maxNode.Attr, " score:", max)
	// 第三步 去掉空标签
	node = maxNode.FirstChild
	for node != nil {
		if node.DataAtom == atom.Div  {
			num := len(this.getChildrens(node, true))
			//fmt.Println("num:",num)
			if num == 0 {
			//空标签，移除掉
				fmt.Println("[Empty]", node.DataAtom.String(), node.Attr)
				node = this.removeAndGetNext(node)
				continue
			} else if num == 1 {
			// 只有一个子元素的标签，移除父标签
				fmt.Println("[Block]", node.DataAtom.String(), node.Attr)

				child := this.getChildrens(node, true)[0].Node
				//fmt.Println("[Block child]", child.DataAtom.String(), child.Attr)
				newNode := &html.Node {
					Type: child.Type,
					DataAtom:child.DataAtom,
					Data: child.Data,
					Namespace:child.Namespace,
					Attr:child.Attr,
					FirstChild: child.FirstChild,
					LastChild:child.LastChild,
					Parent:nil,
					NextSibling:nil,
				}
				n := child.FirstChild
				if n!= nil {
					if n.Parent != nil {
						n.Parent = newNode
					}
					for n.NextSibling != nil {
						n.NextSibling.Parent = newNode
						n = n.NextSibling
					}
				}
				node.Parent.AppendChild(newNode)
				node.Parent.RemoveChild(node)
				node = this.getNextNode(newNode.Parent,false)
				continue
			}

		}
		node = this.getNextNode(node, false)
	}
	//fmt.Println(this.renderNode(maxNode))
	// 第四步 清除style
	node = maxNode
	for node != nil {
		if node != nil && node.Type == html.ElementNode && node.DataAtom != atom.Html && len(node.Attr) > 0 {
			var Attr []html.Attribute
			for _, attr := range node.Attr {
				if attr.Key == "src" || attr.Key == "href" {
					Attr = append(Attr, attr)
				}
			}
			fmt.Println("[cleanStyles]:", node.DataAtom.String(), Attr , node.Attr)
			newNode := &html.Node {
				Type: node.Type,
				DataAtom:node.DataAtom,
				Data: node.Data,
				Namespace:node.Namespace,
				Attr: Attr,
				FirstChild: node.FirstChild,
				LastChild:node.LastChild,
				Parent:nil,
				NextSibling:nil,
			}
			n := node.FirstChild
			if n != nil {
				node.FirstChild.Parent = newNode
				for n.NextSibling != nil {
					n.NextSibling.Parent = newNode
					n = n.NextSibling
				}
			}
			node.Parent.InsertBefore(newNode,node)
			node.Parent.RemoveChild(node)
		}
		node = this.getNextNode(node, false)
	}

	return this.renderNode(maxNode)
}

func (this *MaxSubSegment) getHtml(url string) string {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, bytes.NewReader([]byte("")));
	req.Header.Set("User-Agent","Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.139 Safari/537.36")
	resp, _ := client.Do(req)
	str,_ := ioutil.ReadAll(resp.Body)
	defer  resp.Body.Close()
	return string(str)
}

//
func (this *MaxSubSegment) getBody(str string) string {
	reg_body := regexp.MustCompile(`(?sU:<body[\s\S]*>)([\s\S]*)(?sU:</body>)`)
	body := reg_body.FindString(str)

	// 过滤<script> </script>
	//reg_script := regexp.MustCompile(`(?sU:<script[\s\S]*>)([\s\S]*)(?sU:</script>)`)
	reg_script,_ := regexp.Compile(`<script[\S\s]+?</script>`)
	body = reg_script.ReplaceAllString(body, "")

	reg_noscript,_ := regexp.Compile(`<noscript[\S\s]+?</noscript>`)
	body = reg_noscript.ReplaceAllString(body, "")

	// 过滤<style> </style>
	reg_style := regexp.MustCompile(`(?sU:<style[\s\S]*>)([\s\S]*)(?sU:</style>)`)
	body = reg_style.ReplaceAllString(body,"")

	// 过滤注释
	reg_anno := regexp.MustCompile(`(?sU:<!).*?(?sU:>)`)
	body = reg_anno.ReplaceAllString(body,"")

	return body
}


func (this *MaxSubSegment) getChildrens(node *html.Node, ignore bool) []*ChildrenNode {
	children := []*ChildrenNode{}
	var f func(*html.Node, int)
	//level := 0
	f = func(n *html.Node, level int) {
		for c:= n.FirstChild; c != nil; c = c.NextSibling {
			if ignore && c.Type == html.TextNode && strings.TrimSpace(strings.Replace(c.Data, "\n","",0)) == "" {

			} else {
				child := & ChildrenNode{}
				child.Node = c
				child.Level = level
				children = append(children, child)
			}
			f(c,level+1)
		}
	}
	f(node,1)
	return children
}

func (this *MaxSubSegment) removeAndGetNext(node *html.Node) *html.Node {
	next := this.getNextNode(node, true)
	node.Parent.RemoveChild(node)
	return next
}

func (this *MaxSubSegment) getNextNode(node *html.Node, ignoreSelfAndKids bool) *html.Node {
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

func (this *MaxSubSegment) renderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	//reg := regexp.MustCompile(`(?sU:</*html.*>)|(?sU:</*head.*>)|(?sU:</*body.*>)`)
	reg := regexp.MustCompile(`(?sU:</*html.*>)|(?sU:</*head.*>)`)
	return reg.ReplaceAllString(buf.String(), "")
	//return buf.String()
}
