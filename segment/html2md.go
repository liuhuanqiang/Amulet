package segment

import (

)
import (
	"regexp"
	"github.com/golang/glog"
)

type Html2MD struct {

}

func(this *Html2MD) Init() {

}

func (this *Html2MD) Change(article string) {
	ret := article
	reg1 := regexp.MustCompile(`(?sU:<h1.*>)(.*?)(?U:</h1>)`)
	reg2 := regexp.MustCompile(`(?sU:<h2.*>)(.*?)(?U:</h2>)`)
	reg3 := regexp.MustCompile(`(?sU:<h3.*>)(.*?)(?U:</h3>)`)
	reg4 := regexp.MustCompile(`(?sU:<p.*>)(.*?)(?U:</p>)`)
	reg5 := regexp.MustCompile(`(?sU:<a href="(.*?)" >)(.*?)(?U:</a>)`)
	reg6 := regexp.MustCompile(`(?sU:<li.*>)(.*?)(?U:</li>)`)

	ret = reg1.ReplaceAllString(ret, "\n"+"# "+"$1"+"\n")
	ret = reg2.ReplaceAllString(ret, "\n"+"## "+"$1"+"\n")
	ret = reg3.ReplaceAllString(ret, "\n"+"### "+"$1"+"\n")
	ret = reg4.ReplaceAllString(ret, "$1" + "\n")
	ret = reg5.ReplaceAllString(ret, "[" + "$2" + "](" + "$1" + ")")
	ret = reg6.ReplaceAllString(ret,  "* " + "$1" + "\n")
	glog.Info(ret)
}

