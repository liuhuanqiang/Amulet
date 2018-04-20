package segment

import (
	"regexp"
	"github.com/golang/glog"
	"strings"
	"net/http"
	"fmt"
	"io/ioutil"
	"github.com/PuerkitoBio/goquery"
)

const (
	Host_Other = 0
	Host_ZhiHu = 1
	Host_JianShu = 2
)

type Html2MD struct {

}

func(this *Html2MD) Init() {

}

func (this *Html2MD) regAnno(article string) string {
	reg := regexp.MustCompile(`(?sU:<!--).*?(?s:-->)`)
	reg.ReplaceAllStringFunc(article, func(str string) string {
		//glog.Info("str:",str)
		return ""
	})
	return reg.ReplaceAllString(article, "")
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

func (this *Html2MD) regH(article string) string {
	ret := article
	reg_h1 := regexp.MustCompile(`(?sU:<h1.*>)(?s:.*?)(?U:</h1>)`)
	reg_h2 := regexp.MustCompile(`(?sU:<h2.*>)(?s:.*?)(?U:</h2>)`)
	reg_h3 := regexp.MustCompile(`(?sU:<h3.*>)(?s:.*?)(?U:</h3>)`)

	ret = reg_h1.ReplaceAllStringFunc(ret, func(str string) string{
		reg := regexp.MustCompile(`(?sU:</*.*>)`)
		str = "\n# " + reg.ReplaceAllString(str, "") + "\n"
		//glog.Info("h1:", str)
		return str
	})

	ret = reg_h2.ReplaceAllStringFunc(ret, func(str string) string{
		reg := regexp.MustCompile(`(?sU:</*.*>)`)
		str = "\n## " + reg.ReplaceAllString(str, "") + "\n"
		//glog.Info("h2:", str)
		return str
	})

	ret = reg_h3.ReplaceAllStringFunc(ret, func(str string) string{
		reg := regexp.MustCompile(`(?sU:</*.*>)`)
		str = "\n### " + reg.ReplaceAllString(str, "") + "\n"
		//glog.Info("h3:", str)
		return str
	})

	return ret
}

func (this *Html2MD) regP(article string) string {
	reg_p := regexp.MustCompile(`(?sU:<p.*>)(?s:.*?)(?U:</p>)`)
	return reg_p.ReplaceAllStringFunc(article, func(str string) string{
		//glog.Info("p:", str)
		return regexp.MustCompile(`(?sU:</*p.*>)`).ReplaceAllString(str,"")
	})
	//return reg_p.ReplaceAllString(article, "\n" + "$2" + "\n")
}

func (this *Html2MD) regLi(article string) string {
	reg_li := regexp.MustCompile(`(?sU:<li.*>)(.*?)(?sU:</li>)`)
	return reg_li.ReplaceAllString(article,  "\n* " + "$1" + "\n")
}

func (this *Html2MD) regLink(article string, url string) string {

	reg_a := regexp.MustCompile(`(?sU:<a.*>)(?s:.*?)(?sU:</a>)`)
	return reg_a.ReplaceAllStringFunc(article, func(str string) string {
		//glog.Info(str)
		reg1 := regexp.MustCompile(`href="(?s:.*?)"`)
		reg2 := regexp.MustCompile(`title="(?s:.*?)"`)
		reg3 := regexp.MustCompile(`(?U:</*a.*>)`)

		ret1 :=  strings.Split(reg1.FindString(str),`=`)
		href := ""
		if len(ret1) >= 2 {
			href = strings.Replace(ret1[1],`"`,``,-1)
			if !strings.Contains(href, "https:") && !strings.Contains(href, "http:"){
				href = url + href
			}
		}
		ret2 := strings.Split(reg2.FindString(str),`=`)
		title := ""
		//glog.Info("ret2:",ret2, " len:", len(ret2))
		if len(ret2) >= 2 {
			title = strings.Replace(ret2[1],`"`,``,-1)
		}
		ret3 := reg3.ReplaceAllString(str, "")

		if strings.TrimSpace(href) != "" {
			s := ""
			if strings.TrimSpace(title) != "" {
				s = "["+ strings.TrimSpace(title) + "](" + strings.TrimSpace(href) + ")"
			}

			if strings.TrimSpace(ret3) != "" {
				if strings.Contains(ret3,"<img") {
					s = ret3
				} else {
					if s == "" {
						s = "["+ strings.TrimSpace(ret3) + "](" + strings.TrimSpace(href) + ")"
					}
				}
			}
			glog.Info("s",s)
			if s != "" {
				return s
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

	reg_pre := regexp.MustCompile(`(?sU:<pre.*>)(?s:.*?)(?U:</pre>)`)
	article = reg_pre.ReplaceAllStringFunc(article, func(str string) string {
		glog.Info("code:", str)
		// &lt; &gt;&amp;&quot;&copy;分别是<，>，&，"，©;的转义字符
		str = regexp.MustCompile(`(?U:</*pre.*>)`).ReplaceAllString(str, "")
		str = regexp.MustCompile(`(?U:</*code.*>)`).ReplaceAllString(str, "")
		//glog.Info("code:", str)
		str = strings.Replace(str, `&lt;`, `<`, -1)
		str = strings.Replace(str, `&quot;`, `"`, -1)
		str = strings.Replace(str, `&gt;`, `>`, -1)
		str = strings.Replace(str, `&amp;`, `&`, -1)
		//glog.Info("code:", str)
		return "\n```\n"+ str + "\n```\n\n"
	})

	reg_code := regexp.MustCompile(`(?U:<code>)(?s:.*?)(?U:</code>)`)
	article = reg_code.ReplaceAllStringFunc(article, func(str string) string {
		str = regexp.MustCompile(`(?U:</*code.*>)`).ReplaceAllString(str, "")
		return "`" + str + "`"
	})
	return article
}

func (this *Html2MD) regImg(article string, url string) string {
	reg_img := regexp.MustCompile(`(?sU:<img.*/*>)`)
	article = reg_img.ReplaceAllStringFunc(article, func(str string) string {
		//glog.Info("regImg:", str)
		reg1 := regexp.MustCompile(`src="(.*?)"`)
		ret1 :=  strings.Split(reg1.FindString(str),`=`)
		src := ""
		if len(ret1) >= 2 {
			src = strings.Replace(ret1[1],`"`,``,-1)
			if !strings.Contains(src, "https:") && !strings.Contains(src, "http:"){
				src = url + src
			}
		}
		return "![" + strings.TrimSpace(src) + "](" + strings.TrimSpace(src) + ")"
	})

	//return reg_img.ReplaceAllString(article, "![" + "$1" + "](" + "$1" + ")")
	return article
}

func (this *Html2MD) regStrong(article string) string {
	reg_strong := regexp.MustCompile(`(?sU:<strong.*>)(.*?)(?U:</strong>)`)
	return reg_strong.ReplaceAllString(article, "**" + "$1" + "**" )
}

func (this *Html2MD) regSpan(article string) string {
	reg_span := regexp.MustCompile(`(?sU:</*span.*>)`)
	return reg_span.ReplaceAllString(article, "")
}

func (this *Html2MD) regArticle(ar string) string {
	reg_article := regexp.MustCompile(`(?sU:</*article.*>)`)
	ar = reg_article.ReplaceAllStringFunc(ar, func(str string) string {
		//glog.Info("article:", str)
		return ""
	})
	return ar
}


func (this *Html2MD) regOther(article string) string {
	reg_footer := regexp.MustCompile(`(?sU:<footer.*>)(?sU:.*?)(?sU:</footer>)`)
	article = reg_footer.ReplaceAllString(article, "")

	reg_header := regexp.MustCompile(`(?sU:</*header.*>)`)
	article = reg_header.ReplaceAllString(article, "")

	reg_link := regexp.MustCompile(`(?sU:</*link.*>)`)
	article = reg_link.ReplaceAllString(article, "")

	reg_time := regexp.MustCompile(`(?sU:</*time.*>)`)
	article = reg_time.ReplaceAllString(article, "")

	reg_meta := regexp.MustCompile(`(?sU:</*meta.*>)`)
	article = reg_meta.ReplaceAllString(article, "")

	reg_i := regexp.MustCompile(`(?sU:</*i[^mg].*>)`)
	article = reg_i.ReplaceAllString(article, "")

	reg_button := regexp.MustCompile(`(?sU:<button.*>)(.*?)(?sU:</button>)`)
	article = reg_button.ReplaceAllString(article, "")

	reg_hr := regexp.MustCompile(`(?sU:</*hr.*>)`)
	article = reg_hr.ReplaceAllString(article, "")

	reg_small := regexp.MustCompile(`(?sU:<small.*>)(.*?)(?sU:</small>)`)
	article = reg_small.ReplaceAllString(article, "$1")

	reg_blockquote := regexp.MustCompile(`(?sU:<blockquote.*>)(?s:.*?)(?U:</blockquote>)`)
	article = reg_blockquote.ReplaceAllStringFunc(article, func(str string) string {
		str = regexp.MustCompile(`(?U:</*blockquote.*>)`).ReplaceAllString(str, "")
		return " > " + str
	})
	//article = reg_blockquote.ReplaceAllString(article, "\n > " + "$1" )


	reg_br := regexp.MustCompile(`(?sU:</*br.*>)`)
	article = reg_br.ReplaceAllString(article, "\n")

	reg_b := regexp.MustCompile(`(?sU:<b.*>)(.*?)(?sU:</b>)`)
	article = reg_b.ReplaceAllString(article, "\n*" + "$1")

	reg_strong := regexp.MustCompile(`(?sU:<em.*>)(.*?)(?U:</em>)`)
	article = reg_strong.ReplaceAllString(article, "**" + "$1" + "**" )
	return  article
}

func (this *Html2MD) Change(url string, host string) string {
	glog.Info("url:", url, " host:",host, " getHost:", this.getHost(url))
	return this.getJianShuContent(this.getHtml(url));
}

// 获取链接的host, 针对不同的host,截取不同的值。 比如知乎， 简书 ，博客 hexo wordprss,
func (this *Html2MD) getHost(url string) int {
	reg := regexp.MustCompile(`\.(.*?)\.`)
	ret := reg.FindAllString(url, -1)
	if len(ret) >= 1 {
		if strings.Contains(ret[0],"zhihu") {
			// 知乎
			return Host_ZhiHu
		} else if strings.Contains(ret[0],"jianshu") {
			// 简书
			return Host_JianShu
		}
	}
	return Host_Other
}

func (this *Html2MD) getHtml(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("http get error.")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("http read error")
	}
	return string(body)
}

// 获取正文，去掉广告， 推荐 等无用信息
func (this *Html2MD) getZhiHuContent(html string) string {
	reg := regexp.MustCompile(`(?sU:<article.*>)(?s:.*?)(?U:</article>)`)
	return reg.FindString(html)
}

// 获取简书的
func (this *Html2MD) getJianShuContent(html string) string {
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
	str,_ := doc.Find(".article").Html()
	// 过滤掉div
	str = this.regDiv(str)
	str = this.regAnno(str)
	str = strings.Replace(str, "\n", "", -1)
	return str
}
