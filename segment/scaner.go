package segment

import (
	"os"
	"github.com/golang/glog"
	"bufio"
	"fmt"
	"path/filepath"
	"strings"
	"io"
	"unicode/utf8"
	"unicode"
	"sort"
)

type Scaner struct {
	Dict  map[string]*Token
	Interpunction  map[string]bool
}

type Token struct {
	Word   	string
	Act     string
	Num     uint32
}

func (this *Scaner) Init() {
	this.Dict = make(map[string]*Token)
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
}

func (this *Scaner) GetWord(dir string) {
	paths := []string{}
	filepath.Walk(dir, func(filename string, fi os.FileInfo, err error) error { //遍历目录
		if err != nil {
			return nil
		}
		if fi.IsDir() { // 忽略目录
			return nil
		}
		paths = append(paths,filename)
		return nil
	})
	for _, path := range paths {
		dictFile, err := os.Open(path)
		if err != nil {
			glog.Info("dictionary初始化失败, error:", err.Error()," path:",path)
		}
		//glog.Info("path:",path)
		reader := bufio.NewReader(dictFile)
		var text string
		// 逐行读入分词
		for {
			_, error := fmt.Fscan(reader, &text)
			if error != io.EOF {
				if error != nil {
					glog.Info("error:", error.Error(), "  text:", text,"  path:", path)
					//return
					break
				}
				//glog.Info("text:",text)
				//处理text
				str := strings.Split(text, "/")
				if len(str) >= 2 && this.isValid(str[0]) == 1 {
					tmp := str[0]
					for v :=range this.Interpunction {
						tmp = strings.TrimSpace(strings.Replace(tmp, v, "", -1))
					}
					if len(tmp) > 0 {
						_, exist := this.Dict[tmp]
						if !exist {
							token := &Token{}
							token.Word = tmp
							token.Act = str[1]
							token.Num = 0
							this.Dict[tmp] = token
						} else {

							this.Dict[tmp].Num++
						}
					}
				}
			} else {
				break
			}
		}
		dictFile.Close()
	}
	this.Filter()
}

func (this *Scaner) WriteWord() {
	fl, err := os.OpenFile("", os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		glog.Info("error:", err.Error())
	}
	defer fl.Close()
	for _, token := range this.Dict {
		if token.Act != "w" {
			str := token.Word
			str = strings.Replace(str, "[", "", -1)
			_, err1 := fl.Write([]byte( str + "\n"))
			if err1 != nil {
				glog.Info("error:", err1.Error())
			}
		}
	}
}

//过滤数据
func (this *Scaner) Filter() {
	glog.Info("Filter")
	//dictFile, err := os.OpenFile("data/out1.txt", os.O_RDWR|os.O_CREATE,0766)
	//if err != nil {
	//	glog.Info("error:", err.Error())
	//}
	//defer dictFile.Close()
	//reader := bufio.NewReader(dictFile)
	//// 逐行读入分词
	//var text string
	//letterNum := 0
	//num := 0
	//existNum := 0
	//for {
	//	size ,_ := fmt.Fscanln(reader, &text)
	//	glog.Info("size:", size, "  text:", text, " len:", len([]rune(text)))
	//	num++
	//	if size == 0 {
	//		// 文件结束
	//		break
	//	}
	//	if len([]rune(text)) == 1 {
	//		continue
	//	}
	//	runes := []byte(text)
	//	current := 0
	//	isLetter := false
	//	//glog.Info("current:", current, " runes:",len(runes))
	//	for current < len(runes) {
	//		//glog.Info("213131313")
	//		r, size := utf8.DecodeRune(runes[current:])
	//		// 如果是英文字母
	//		//glog.Info("size:", size)
	//		if size <= 2 && (unicode.IsLetter(r) || unicode.IsNumber(r)) {
	//			isLetter = true
	//			break
	//		} else {
	//			current += size
	//		}
	//	}
	//	if isLetter {
	//		letterNum++
	//		continue
	//	}
	//	_, exist := this.Dict[text]
	//	if exist {
	//		existNum++
	//		this.Dict[text].Num++
	//		continue
	//	}
	//	token := &Token{}
	//	token.Word = text
	//	token.Num++
	//	this.Dict[text] = token
	//}
	//glog.Info("读取完毕 ", "num:", num, "  letterNum:",letterNum, "  existNum:",existNum)
	var keys []string
	for k := range this.Dict {
		keys = append(keys, k)
	}
	sort.Sort(sort.StringSlice(keys))
	fl, err := os.OpenFile("data/out3.txt", os.O_RDWR|os.O_CREATE,0766)
	if err != nil {
		glog.Info("error:", err.Error())
	}
	defer fl.Close()
	for _, str := range keys {
		tmp := fmt.Sprintf("%s %d", str, this.Dict[str].Num)
		//str = strings.Replace(str, "[", "", -1)
		_, err1 := fl.Write([]byte( tmp + "\n"))
		if err1 != nil {
			glog.Info("error:", err1.Error())
		}
	}
}

func (this *Scaner) isValid(text string) int {
	if len([]rune(text)) == 1 {
		glog.Info("text:", text)
		return 1
	}
	runes := []byte(text)
	current := 0
	isLetter := false
	//glog.Info("current:", current, " runes:",len(runes))
	for current < len(runes) {
		//glog.Info("213131313")
		r, size := utf8.DecodeRune(runes[current:])
		// 如果是英文字母
		//glog.Info("size:", size)
		if size <= 2 && (unicode.IsLetter(r) || unicode.IsNumber(r)) {
			isLetter = true
			break
		} else {
			current += size
		}
	}
	if isLetter {
		return 0
	}
	return 2
}