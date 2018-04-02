package segment

import (
	"os"
	"github.com/golang/glog"
	"bufio"
	"strings"
	"unicode/utf8"
	"fmt"
	"unicode"
)

type Segment struct {
	Interpunction  map[string]bool
	Dict           map[string]bool
	Graph          [][]*Node
}

type Node struct {
	From   int
	To     int
	Cost   float32
	Text   string
}

type Token struct {
	Text  	string
	Weight  float32
}

func (this *Segment) Init() {
	this.Interpunction = make(map[string]bool)

	file, err := os.Open("data/interpunction.txt")
	if err != nil {
		glog.Info("interpunction初始化失败, error:",err.Error())
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text != "" {
			this.Interpunction[text] = true
		}
	}
	glog.Info("interpunction初始化完成")
	this.Dict = make(map[string]bool)
	dictFile, err := os.Open("data/dictionary.txt")
	if err != nil {
		glog.Info("dictionary初始化失败, error:", err.Error())
	}
	defer dictFile.Close()
	reader := bufio.NewReader(dictFile)
	var text string
	var freqText string
	var pos string

	// 逐行读入分词
	for {
		size, _ := fmt.Fscanln(reader, &text, &freqText, &pos)

		if size == 0 {
			// 文件结束
			break
		} else if size < 2 {
			// 无效行
			continue
		} else if size == 2 {
			// 没有词性标注时设为空字符串
			pos = ""
		}
		this.Dict[text] = true
	}

	glog.Info("dictionary初始化完成")
}

func (this *Segment) isInterpution(str string) bool{
	if str == " " {
		return true
	}
	_, exist := this.Interpunction[str]
	if exist {
		return true
	} else {
		return false
	}
}

func (this *Segment) isTextExist(str string) bool {
	_, exist := this.Dict[str]
	if exist {
		return true
	} else {
		return false
	}
}

// 将一个文章, 分成一段一段的，
func (this *Segment) SplitToSegment(text []byte) []string{
	current := 0
	output := []string{}
	var tmp string
	inAlphanumeric := true
	alphanumericStart := 0
	for current < len(text) {
		r, size := utf8.DecodeRune(text[current:])
		// 如果是英文字母
		if size <= 2 && (unicode.IsLetter(r) || unicode.IsNumber(r)) {
			if !inAlphanumeric {
				alphanumericStart = current
				inAlphanumeric = true
			}
			if len(tmp) > 0 {
				output = append(output, tmp)
			}
			tmp = ""

		} else {
			if inAlphanumeric {
				inAlphanumeric = false
				if current != 0 {
					output = append(output, string(toLower(text[alphanumericStart:current])))
				}
			} else {
				if this.isInterpution(string(r)) {
					if len(tmp) > 0 {
						output = append(output, tmp)
					}
					tmp = ""
				} else {
					tmp += string(r)
				}
			}
		}
		current += size
	}

	if len(tmp) > 0{
		output = append(output, tmp)
	}
	if inAlphanumeric {
		if current != 0 {
			output = append(output, string(toLower(text[alphanumericStart:current])))
		}
	}
	for _, v :=range output {
		glog.Info(v)
	}
	return output
}


//最大逆向匹配算法
func (this *Segment) MaxReverse(str string) []string{
	// 将字符串转化成byte
	runes := []rune(str)
	current := 0
	right := len(runes)
	left := 0
	output := []string{}
	r, size := utf8.DecodeRune([]byte(str)[left:])
	// 如果是英文字母
	if size <= 2 && (unicode.IsLetter(r) || unicode.IsNumber(r)) {
		output = append(output, str)
		return output
	}
	length := 0
	for right > 0{
		// 取出第一个文字
		leftText := runes[current:right] // 剩下的文字
		if this.isTextExist(string(leftText)) {
			// 如果文本存在
			output = append(output, string(leftText))
			right = current
			current = 0
			length++
		} else {
			// 如果文本不存在 剩下一个汉字
			if len(leftText) <= 1 {
				output = append(output, string(leftText))
				length++
			}
			current++
		}
	}
	ret := make([]string, length)
	for i := 0; i < length; i++ {
		ret[i] = output[length - i - 1]
	}
	return ret
}

// N最短路径算法  1.找出字符串中所有可能的词，构造词切分有向无环图  2.寻找最短路径
// 分析每一个字，找出以这个字开头的所有词语
func (this *Segment) NShort(str string) [][]string {
	output := [][]string{}
	// 将字符串转换为rune
	runes := []rune(str)
	// 分析每一个字，找出以这个字开头的所有词语
	for i := 0; i < len(runes); i++ {
		list := []string{}
		for j := i + 1; j <= len(runes); j++ {
			text := runes[i:j]
			if this.isTextExist(string(text)) {
				list = append(list, string(text))
			}
		}
		output = append(output, list)
	}
	// 构建有向无环图
	this.Graph = [][]*Node{}
	if len(output) > 0 {
		for _, v := range output[0] {
			var ret []*Node
			node := &Node{}
			node.Text = v
			node.From = 0
			node.To = len([]rune(v))
			this.buildNode(node, output, len([]rune(v)), ret)
			glog.Info(ret)
		}
	}
	glog.Info(this.Graph)
	// 计算最短路径

	return output
}

//
func (this *Segment) buildNode(node *Node, output [][]string, from int, ret []*Node) {
	if from >= len(output) {
		ret = append(ret, node)
		for _,v := range ret {
			glog.Info("v:",v)
		}
		this.Graph = append(this.Graph,ret)
		return
	}
	ret = append(ret, node)
	for _, v := range output[from] {
		//glog.Info("v:",v)
		node := &Node{}
		node.Text = v
		node.From = from
		node.To = len([]rune(v)) + from
		this.buildNode(node, output, len([]rune(v)) + from, ret)
	}
}

// 将英文词转化为小写
func toLower(text []byte) []byte {
	output := make([]byte, len(text))
	for i, t := range text {
		if t >= 'A' && t <= 'Z' {
			output[i] = t - 'A' + 'a'
		} else {
			output[i] = t
		}
	}
	return output
}