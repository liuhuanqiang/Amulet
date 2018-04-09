package segment

import (
	"regexp"
	"github.com/golang/glog"
	"strings"
)

type Html2MD struct {

}

func(this *Html2MD) Init() {

}

// 过滤<script></script>
func(this *Html2MD) regScript(article string) string {
	reg_script := regexp.MustCompile(`(?sU:<script.*>)(.*?)(?U:</script>)`)
	return reg_script.ReplaceAllString(article, "")
}

// 过滤<style></style>
func (this *Html2MD) regStyle(article string) string {
	reg_style := regexp.MustCompile(`(?sU:<style.*>)(.*?)(?U:</style>)`)
	return reg_style.ReplaceAllString(article, "")
}

// 过滤<div></div>
func (this *Html2MD) regDiv(article string) string {
	reg_div := regexp.MustCompile(`(?sU:</*div.*>)`)
	return reg_div.ReplaceAllString(article, "")
}

// 过滤<ul></ul>
func (this *Html2MD) regUl(article string) string {
	reg_ul := regexp.MustCompile(`(?sU:</*ul.*>)`)
	return reg_ul.ReplaceAllString(article, "\n")
}

//过滤 <hr/>
func (this *Html2MD) regHr(article string) string {
	reg_hr := regexp.MustCompile(`(?sU:</*hr.*>)`)
	return reg_hr.ReplaceAllString(article, "")
}

func (this *Html2MD) regSmall(article string) string {
	reg_small := regexp.MustCompile(`(?sU:<small.*>)(.*?)(?U:</small>)`)
	return reg_small.ReplaceAllString(article, "")
}

func (this *Html2MD) regH(article string) string {
	ret := article
	reg_h1 := regexp.MustCompile(`(?sU:<h1.*>)(?s:.*?)(?sU:</h1>)`)
	reg_h2 := regexp.MustCompile(`(?sU:<h2.*>)(?s:.*?)(?sU:</h2>)`)
	reg_h3 := regexp.MustCompile(`(?sU:<h3.*>)(?s:.*?)(?sU:</h3>)`)

	ret = reg_h1.ReplaceAllStringFunc(ret, func(str string) string{
		reg := regexp.MustCompile(`(?U:</*.*>)`)
		str = "\n# " + reg.ReplaceAllString(str, "") + "\n"
		glog.Info("h1:", str)
		return str
	})

	ret = reg_h2.ReplaceAllStringFunc(ret, func(str string) string{
		reg := regexp.MustCompile(`(?U:</*.*>)`)
		str = "\n## " + reg.ReplaceAllString(str, "") + "\n"
		glog.Info("h2:", str)
		return str
	})

	ret = reg_h3.ReplaceAllStringFunc(ret, func(str string) string{
		reg := regexp.MustCompile(`(?U:</*.*>)`)
		str = "\n### " + reg.ReplaceAllString(str, "") + "\n"
		glog.Info("h3:", str)
		return str
	})

	return ret
}

func (this *Html2MD) regP(article string) string {
	reg_p := regexp.MustCompile(`(?sU:<p.*>)(.*?)(?U:</p>)`)
	reg_p.ReplaceAllStringFunc(article, func(str string) string{
		glog.Info("p:", str)
		return str
	})
	return reg_p.ReplaceAllString(article, "\n" + "$1" + "\n")
}

func (this *Html2MD) regLi(article string) string {
	reg_li := regexp.MustCompile(`(?sU:<li.*>)(.*?)(?sU:</li>)`)
	return reg_li.ReplaceAllString(article,  "\n* " + "$1" + "\n")
}

func (this *Html2MD) regLink(article string) string {
	reg_a := regexp.MustCompile(`(?sU:<a.*>)(?s:.*?)(?sU:</a>)`)
	return reg_a.ReplaceAllStringFunc(article, func(str string) string {
		glog.Info(str)
		reg1 := regexp.MustCompile(`href="(?s:.*?)"`)
		reg2 := regexp.MustCompile(`title="(?s:.*?)"`)
		reg3 := regexp.MustCompile(`(?U:</*a.*>)`)

		ret1 :=  strings.Split(reg1.FindString(str),`=`)
		href := ""
		if len(ret1) >= 2 {
			href = strings.Replace(ret1[1],`"`,``,-1)
		}
		ret2 := strings.Split(reg2.FindString(str),`=`)
		title := ""
		//glog.Info("ret2:",ret2, " len:", len(ret2))
		if len(ret2) >= 2 {
			title = strings.Replace(ret2[1],`"`,``,-1)
		}
		ret3 := reg3.ReplaceAllString(str, "")
		if strings.TrimSpace(href) != "" {
			if strings.TrimSpace(title) != "" {
				return "["+ strings.TrimSpace(title) + "](" + strings.TrimSpace(href) + ")"
			}

			if strings.TrimSpace(ret3) != "" {
				return "["+ strings.TrimSpace(ret3) + "](" + strings.TrimSpace(href) + ")"
			}

			return "["+ strings.TrimSpace(href) + "](" + strings.TrimSpace(href) + ")"
		} else {
			if strings.TrimSpace(ret3) != "" {
				return strings.TrimSpace(ret3)
			} else {
				return ""
			}
		}
		return ""
	})
}

func (this *Html2MD) regCode(article string) string {
	reg_pre := regexp.MustCompile(`(?sU:<pre.*>)(.*?)(?U:</pre>)`)
	return  reg_pre.ReplaceAllStringFunc(article, func(str string) string {
		//glog.Info("code:", str)
		// &lt; &gt;&amp;&quot;&copy;分别是<，>，&，"，©;的转义字符
		str = regexp.MustCompile(`(?U:</*pre.*>)`).ReplaceAllString(str, "")
		glog.Info("code:", str)
		str = strings.Replace(str, `&lt;`, `<`, -1)
		str = strings.Replace(str, `&quot;`, `"`, -1)
		str = strings.Replace(str, `&gt;`, `>`, -1)
		str = strings.Replace(str, `&amp;`, `&`, -1)
		glog.Info("code:", str)
		return "\n```\n"+ str + "\n```\n\n"
	})
}

func (this *Html2MD) regImg(article string) string {
	reg_img := regexp.MustCompile(`(?sU:<img src="(.*)".*/>)`)
	return reg_img.ReplaceAllString(article, "![" + "$1" + "](" + "$1" + ")")
}

func (this *Html2MD) regStrong(article string) string {
	reg_strong := regexp.MustCompile(`(?sU:<strong.*>)(.*?)(?U:</strong>)`)
	return reg_strong.ReplaceAllString(article, "**" + "$1" + "**" )
}


func (this *Html2MD) Change(article string) {
	ret := article

	//reg_script := regexp.MustCompile(`(?sU:<script.*>)(.*?)(?U:</script>)`)
	//reg_style := regexp.MustCompile(`(?sU:<style.*>)(.*?)(?U:</style>)`)
	//reg_div := regexp.MustCompile(`(?sU:</*div.*>)`)
	//reg_ul := regexp.MustCompile(`(?sU:</*ul.*>)`)
	//reg_hr := regexp.MustCompile(`(?sU:</*hr.*>)`)
	//reg_pre := regexp.MustCompile(`(?sU:<pre.*>)(.*?)(?U:</pre>)`)
	//reg_h1 := regexp.MustCompile(`(?sU:<h1.*>)(.*?)(?U:</h1>)`)
	//reg_h2 := regexp.MustCompile(`(?sU:<h2.*>)(.*?)(?U:</h2>)`)
	//reg_h3 := regexp.MustCompile(`(?sU:<h3.*>)(.*?)(?U:</h3>)`)
	//reg_p := regexp.MustCompile(`(?sU:<p.*>)(.*?)(?U:</p>)`)
	//reg_a := regexp.MustCompile(`(?sU:<a.*>)(?s:.*?)(?sU:</a>)`)
	//reg_li := regexp.MustCompile(`(?sU:<li.*>)(.*?)(?sU:</li>)`)
	//reg_img :=regexp.MustCompile(`(?sU:<img src="(.*)".*/>)`)
	//reg_small := regexp.MustCompile(`(?sU:<small.*>)(.*?)(?U:</small>)`)
	//
	//ret = reg_script.ReplaceAllString(ret, "")
	//ret = reg_style.ReplaceAllString(ret, "")
	//ret = reg_div.ReplaceAllString(ret, "")
	//ret = reg_hr.ReplaceAllString(ret, "")
	//ret = reg_ul.ReplaceAllString(ret, "")
	//ret = reg_pre.ReplaceAllStringFunc(ret, func(str string) string {
	//	//glog.Info("code:", str)
	//	// &lt; &gt;&amp;&quot;&copy;分别是<，>，&，"，©;的转义字符
	//	str = regexp.MustCompile(`(?U:</*pre.*>)`).ReplaceAllString(str, "")
	//	glog.Info("code:", str)
	//	str = strings.Replace(str, `&lt;`, `<`, -1)
	//	str = strings.Replace(str, `&quot;`, `"`, -1)
	//	str = strings.Replace(str, `&gt;`, `>`, -1)
	//	str = strings.Replace(str, `&amp;`, `&`, -1)
	//	glog.Info("code:", str)
	//	return "\n```\n"+ str + "\n```\n\n"
	//})
	//
	//ret = reg_h1.ReplaceAllString(ret, "\n"+"# "+"$1"+"\n")
	//ret = reg_h2.ReplaceAllString(ret, "\n"+"## "+"$1"+"\n")
	//ret = reg_h3.ReplaceAllString(ret, "\n"+"### "+"$1"+"\n")
	//ret = reg_p.ReplaceAllString(ret, "$1" + "\n")
	//ret = reg_small.ReplaceAllString(ret, "$1" + "\n")
	//ret = reg_li.ReplaceAllString(ret,  "\n* " + "$1" + "\n")
	//ret = reg_img.ReplaceAllString(ret, "![" + "$1" + "](" + "$1" + ")")
	//ret = reg_a.ReplaceAllStringFunc(ret, func(str string) string {
	//	glog.Info(str)
	//	reg1 := regexp.MustCompile(`href="(?s:.*?)"`)
	//	reg2 := regexp.MustCompile(`title="(?s:.*?)"`)
	//	reg3 := regexp.MustCompile(`(?U:</*a.*>)`)
	//
	//	ret1 :=  strings.Split(reg1.FindString(str),`=`)
	//	href := ""
	//	if len(ret1) >= 2 {
	//		href = strings.Replace(ret1[1],`"`,``,-1)
	//	}
	//	ret2 := strings.Split(reg2.FindString(str),`=`)
	//	title := ""
	//	glog.Info("ret2:",ret2, " len:", len(ret2))
	//	if len(ret2) >= 2 {
	//		title = strings.Replace(ret2[1],`"`,``,-1)
	//	}
	//	ret3 := reg3.ReplaceAllString(str, "")
	//	if strings.TrimSpace(href) != "" {
	//		if strings.TrimSpace(title) != "" {
	//			return "["+ strings.TrimSpace(title) + "](" + strings.TrimSpace(href) + ")"
	//		}
	//
	//		if strings.TrimSpace(ret3) != "" {
	//			return "["+ strings.TrimSpace(ret3) + "](" + strings.TrimSpace(href) + ")"
	//		}
	//
	//		return "["+ strings.TrimSpace(href) + "](" + strings.TrimSpace(href) + ")"
	//	} else {
	//		if strings.TrimSpace(ret3) != "" {
	//			return strings.TrimSpace(ret3)
	//		} else {
	//			return ""
	//		}
	//	}
	//	return ""
	//})
	ret = this.regScript(ret)
	ret = this.regStyle(ret)
	ret = this.regDiv(ret)
	ret = this.regH(ret)
	ret = this.regLink(ret)
	ret = this.regLi(ret)
	ret = this.regCode(ret)
	ret = this.regP(ret)
	ret = this.regImg(ret)
	ret = this.regHr(ret)
	ret = this.regUl(ret)
	ret = this.regSmall(ret)
	ret = this.regStrong(ret)
	glog.Info(ret)
}

