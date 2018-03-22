package utils

import (
	"strings"
	"regexp"
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