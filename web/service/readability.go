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
	"math"
)

type Readability struct {
	UnLikelyCandidates *regexp.Regexp
	OkMaybeItsACandidate *regexp.Regexp
	Positive *regexp.Regexp
	Negative *regexp.Regexp
}

type CandidateNode struct {
	Node *html.Node
	Score int
}

func (this *Readability) initRegexp(){
	this.UnLikelyCandidates,_ = regexp.Compile(`banner|breadcrumbs|combx|comment|community|cover-wrap|disqus|extra|foot|header|legends|menu|related|remark|replies|rss|shoutbox|sidebar|skyscraper|social|sponsor|supplemental|ad-break|agegate|pagination|pager|popup|yom-remote|nav|browsehappy|modal`)
	this.OkMaybeItsACandidate,_ = regexp.Compile(`and|article|body|column|main|shadow`)
	this.Positive, _ = regexp.Compile(`article|body|content|entry|hentry|h-entry|main|page|pagination|post|text|blog|story`)
	this.Negative,_ = regexp.Compile(`hidden|^hid$| hid$| hid |^hid |banner|combx|comment|com-|contact|foot|footer|footnote|masthead|media|meta|outbrain|promo|related|scroll|share|shoutbox|sidebar|skyscraper|sponsor|shopping|tags|tool|widget`)
}

func (this *Readability) GetContent(url string) string {

	this.initRegexp()
	htmlStr := this.getHtml(url)
	fmt.Println("htmlstr:", htmlStr)
	bodyString := this.getBody(htmlStr)
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
		if this.UnLikelyCandidates.Match([]byte(matchString)) && !this.OkMaybeItsACandidate.Match([]byte(matchString)) && node.DataAtom!= atom.A {
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
				n := child.FirstChild
				n.Parent = newNode
				for n.NextSibling != nil {
					n.NextSibling.Parent = newNode
					n = n.NextSibling
				}
				node.Parent.AppendChild(newNode)
				node.Parent.RemoveChild(node)
				node = this.getNextNode(newNode.Parent,false)
				continue
				//node = newNode
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
					LastChild:child.LastChild,
					Parent:nil,
					NextSibling:nil,
				}
				n := child.FirstChild
				n.Parent = newNode
				for n.NextSibling != nil {
					n.NextSibling.Parent = newNode
					n = n.NextSibling
				}
				node.Parent.AppendChild(newNode)
				node.Parent.RemoveChild(node)
				node = this.getNextNode(newNode.Parent,false)
				continue
				//node = newNode
			} else {
				// todo 如果<div></div>中全是文字，则用<p></p>替换
			}
		}
		node = this.getNextNode(node, false)
	}

	// 计算分数
	fmt.Println("html:", this.renderNode(bodyNode))
	node = bodyNode.FirstChild
	elementsToScore := []*html.Node{}
	for node != nil {
		if node.DataAtom == atom.Section || node.DataAtom == atom.H2 || node.DataAtom == atom.H3 || node.DataAtom == atom.H3 || node.DataAtom == atom.H4 ||
		 	node.DataAtom == atom.H5 || node.DataAtom == atom.H6 || node.DataAtom == atom.P || node.DataAtom == atom.Td || node.DataAtom == atom.Pre {
			elementsToScore = append(elementsToScore, node)
		}
		node = this.getNextNode(node, false)
	}
	candidates := make(map[*html.Node]float32)
	for _, elementToScore := range elementsToScore {
		if elementToScore.Parent == nil {
			continue
		}
		innerText := this.getInnerText(elementToScore)
		//fmt.Println("[sore]:", innerText)
		if len(innerText) < 25 {
			continue
		}
		ancestors := this.getNodeAncestors(elementToScore, 3)
		if len(ancestors) == 0 {
			continue
		}
		contentScore := float32(0)
		// 基本分数
		contentScore += 1
		// 每个逗号加一分
		contentScore += float32(len(strings.Split(innerText, ",")))
		// 每100单词加三分
		contentScore += float32(math.Min(math.Floor(float64(len(innerText))/100),float64(3)))
		//fmt.Println("[test]:1", elementToScore.DataAtom,elementToScore.Attr, elementToScore.Data)
		for level, ancestor := range ancestors {
			//fmt.Println("[test]:",ancestor.DataAtom.String(), ancestor.Attr,ancestor)
			if ancestor.Parent == nil {
				continue
			}
			if _,exist := candidates[ancestor]; !exist {
				// 如果不存在
				candidates[ancestor] = this.getNodeInitScore(ancestor)
			}
			divider := 1
			if level == 0 {
				divider = 1
			}  else if level == 1 {
				divider = 2
			} else {
				divider = level * 3
			}
			candidates[ancestor] += (contentScore/float32(divider))
		}

	}
	// 打印分数
	maxScore := float32(0)
	var topCandidate *html.Node
	for c := range candidates {

		score := candidates[c]* (1 - this.getLinkDensity(c))
		candidates[c] = score
		fmt.Println("[Candidate]:", c.DataAtom.String() ,(*c).Attr, " score:", score)
		if topCandidate == nil {
			maxScore = score
			topCandidate = c
		} else {
			if score > maxScore {
				maxScore = score
				topCandidate = c
			}
		}
	}

	if topCandidate == nil || topCandidate.DataAtom == atom.Body {
		// 如果得分最高的是body或者为空, 将body改成div就可以了
	} else {
		// 如果 寻找分数与之接近的node
		alternativeCandidateAncestors := [][]*html.Node{}
		for c, score := range candidates {
			if float32(score)*0.75 >= float32(maxScore) {
				alternativeCandidateAncestors = append(alternativeCandidateAncestors, this.getNodeAncestors(c,3))
			}
		}
		if len(alternativeCandidateAncestors) >= 3 {
			parentOfTopCandidate := topCandidate.Parent
			listsContainingThisAncestor := 0
			for i := 0; i < len(alternativeCandidateAncestors) && listsContainingThisAncestor < 3;i ++ {
				for _, an := range alternativeCandidateAncestors[i] {
					if an == parentOfTopCandidate {
						listsContainingThisAncestor ++
					}
				}
			}
			if listsContainingThisAncestor >= 3 {
				topCandidate = parentOfTopCandidate
			}
		}

	}
	if _,exist:= candidates[topCandidate]; !exist {
		candidates[topCandidate] = this.getNodeInitScore(topCandidate)
	}
	lastScore := float32(candidates[topCandidate])
	fmt.Println("[lastScore]:", lastScore)
	// 查看他的相邻节点，如果分数相似，则认为也是文章的一部分
	siblingScoreThreshold := float32(10)
	if float32(lastScore)*0.2 > siblingScoreThreshold {
		siblingScoreThreshold = lastScore * 0.2
	}
	parentOfTopCandidate := topCandidate.Parent
	child := parentOfTopCandidate.FirstChild
	fmt.Println("[sibing]1:", parentOfTopCandidate)
	for child != nil {
		if child != nil && topCandidate != child {
			fmt.Println("[sibing]:", child.Attr,child.DataAtom.String(),"  ", child)
			flag := false
			contentBonus := float32(0)
			if child.DataAtom == topCandidate.DataAtom {
				contentBonus += lastScore * 0.2
			}
			if score,exist := candidates[child]; exist && (float32(score) + contentBonus) >= siblingScoreThreshold{
				flag = true
			} else if child.DataAtom == atom.P {
				nodeContent := this.getInnerText(child)
				linkDensity := this.getLinkDensity(child)
				fmt.Println("q11:", len(nodeContent))
				reg,_ := regexp.Compile(`.( |$)`)
				if len(nodeContent) >= 80 && linkDensity < 0.25 {
					flag = true
				} else if len(nodeContent) < 80 && len(nodeContent) > 0 && linkDensity == 0 && reg.Match([]byte(nodeContent)) {
					flag = true
				}
			}
			if !flag {
				next := child.NextSibling
				parentOfTopCandidate.RemoveChild(child)
				child = next
				continue
			}
		}
		child = child.NextSibling
	}
	// 寻找分数最高的节点
	node = bodyNode.FirstChild
	for node != nil {
		if node.Type == html.TextNode  && strings.TrimSpace(strings.Replace(node.Data,"\n","",0)) == "" {
			node = this.removeAndGetNext(node)
		} else {
			node = this.getNextNode(node, false)
		}
	}
	fmt.Println("html:", this.renderNode(bodyNode))
	return this.renderNode(bodyNode)

}

func (this *Readability) getHtml(url string) string {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, bytes.NewReader([]byte("")));
	req.Header.Set("User-Agent","Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.139 Safari/537.36")
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
	//reg := regexp.MustCompile(`(?sU:</*html.*>)|(?sU:</*head.*>)|(?sU:</*body.*>)`)
	reg := regexp.MustCompile(`(?sU:</*html.*>)|(?sU:</*head.*>)`)
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

func (this *Readability) replaceNode(node *html.Node, child *html.Node) {
	// 替换两个node
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
	node.Parent.AppendChild(newNode)
	node.Parent.RemoveChild(node)
	n := newNode.FirstChild
	for n.NextSibling != nil {
		n.NextSibling.Parent = newNode
		n = n.NextSibling
	}
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

func (this *Readability) getInnerText(n *html.Node) string {
	str := ""
	//for c:= n.FirstChild; c != nil; c = c.NextSibling {
	//	if c.Type == html.TextNode {
	//		str += strings.TrimSpace(strings.Replace(c.Data,"\n","",0))
	//	}
	//}
	node := n.FirstChild
	for node != nil {
		if node.Type == html.TextNode {
			str += strings.TrimSpace(strings.Replace(node.Data,"\n","",0))
		}
		node = this.getNextNode(node, false)
	}
	return str
}

func (this *Readability) getNodeAncestors(n *html.Node, depth int) []*html.Node {
	fmt.Println("[test]:1", n.DataAtom,n.Attr, n.Data)
	parents := []*html.Node{}
	i := 0
	for n.Parent != nil && n.Parent.DataAtom != atom.Html {
		parents = append(parents, n.Parent)
		i++
		if i >= depth {
			break
		}
		n = n.Parent
	}
	return parents
}

func (this *Readability) getNodeInitScore(node *html.Node) float32 {
	switch node.DataAtom {
	case atom.Div:
		return 5
	case atom.Pre:
	case atom.Td:
	case atom.Blockquote:
		return 3
	case atom.Address:
	case atom.Ol:
	case atom.Ul:
	case atom.Dl:
	case atom.Dd:
	case atom.Dt:
	case atom.Li:
	case atom.Form:
		return -3
	case atom.H1:
	case atom.H2:
	case atom.H3:
	case atom.H4:
	case atom.H5:
	case atom.H6:
	case atom.Th:
		return -5

	}
	weight := float32(0)
	// 否则看class 和 id
	for _,attr := range node.Attr {
		if attr.Key == "class" || attr.Key == "id"{
			if this.Negative.Match([]byte(attr.Val)) {
				weight -= 25
			}
			if this.Positive.Match([]byte(attr.Val)) {
				weight += 25
			}
		}
	}
	return weight
}

func (this *Readability) getLinkDensity(node *html.Node) float32 {
	innerText := this.getInnerText(node)
	if len(innerText) == 0 {
		return 0
	}
	textLength := float32(len(innerText))
	linkLength := float32(0)
	n := node.FirstChild
	for n != nil {
		if n.DataAtom == atom.A {
			linkLength += float32(len(this.getInnerText(n)))
			n = this.getNextNode(n, true)
		} else {
			n = this.getNextNode(n, false)
		}
	}
	return linkLength/textLength
}

func (this *Readability) cleanClass(node *html.Node) {
	
}