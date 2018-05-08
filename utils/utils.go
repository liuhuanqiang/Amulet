package utils

import (
	"strings"
	"regexp"
	"path/filepath"
	"os"
	"github.com/golang/glog"
)

// 过滤掉无用的字符
func FilterInvalidChar(str string) string {
	// 去掉空格
	str = strings.Replace(str," ","",-1)
	str = strings.TrimSpace(str)
	// 去掉回车
	str = strings.Replace(str,"\n", "", -1)

	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	str = re.ReplaceAllStringFunc(str, strings.ToLower)
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	str = re.ReplaceAllString(str, "")

	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	str = re.ReplaceAllString(str, "")

	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	str = re.ReplaceAllString(str, "")

	return str

}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		glog.Info("获取路径失败")
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func SubString(str string,begin int,length int) (substr string) {
	// 将字符串的转换成[]rune
	rs := []rune(str)
	lth := len(rs)

	// 简单的越界判断
	if begin < 0 {
		begin = 0
	}
	if begin >= lth {
		begin = lth
	}
	end := begin + length
	if end > lth {
		end = lth
	}

	// 返回子串
	return string(rs[begin:end])
}

// 过滤<script></script>
func RegScript(article string) string {
	reg_script := regexp.MustCompile(`(?sU:<script[\s\S]*>)([\s\S]*)(?U:</script>)`)
	return reg_script.ReplaceAllString(article, "")
}

func RegSpan(article string) string {
	reg_span := regexp.MustCompile(`(?sU:</*span.*>)`)
	return reg_span.ReplaceAllString(article, "")
}

// 过滤<div></div>
func RegDiv(article string) string {
	reg_div := regexp.MustCompile(`(?sU:</*div.*>)`)
	return reg_div.ReplaceAllString(article, "")
}

func RegAnno(article string) string {
	reg := regexp.MustCompile(`(?sU:<!--).*?(?s:-->)`)
	reg.ReplaceAllStringFunc(article, func(str string) string {
		//glog.Info("str:",str)
		return ""
	})
	return reg.ReplaceAllString(article, "")
}

func RegH1(article string) string {

	reg_div := regexp.MustCompile(`(?sU:</*h1.*>)`)
	return reg_div.ReplaceAllString(article, "")

}

func RegBr(article string) string {
	reg_br := regexp.MustCompile(`(?sU:</*br.*>)`)
	return reg_br.ReplaceAllString(article, "\n")
}

func RegStyle(article string) string {
	reg_style := regexp.MustCompile(`(?sU:<style.*>)(.*?)(?U:</style>)`)
	return reg_style.ReplaceAllString(article, "")
}

func ReplaceEscapeStr(str string) string {
	str = strings.Replace(str, "&#34;", "\\",0)
	str = strings.Replace(str, "&quot;", "\\", 0)
	str = strings.Replace(str, "&#38;", "&", 0)
	str = strings.Replace(str, "&amp;", "&", 0)
	str = strings.Replace(str, "&#39;", "'", 0)
	str = strings.Replace(str, "&#60;", "<", 0)
	str = strings.Replace(str, "&lt;", "<", 0)
	str = strings.Replace(str, "&#62;", ">", 0)
	str = strings.Replace(str, "&gt;", ">", 0)
	str = strings.Replace(str, "&#160;", " ",0)
	str = strings.Replace(str, "&nbsp;", " ", 0)
	str = strings.Replace(str, "\n", "",0)
	return  str
}