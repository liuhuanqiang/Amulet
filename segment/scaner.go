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
}

type Token struct {
	Word   	string
	Act     string
}

func (this *Scaner) Init() {
	this.Dict = make(map[string]*Token)
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
				if len(str) >= 2 {
					token := &Token{}
					token.Word = str[0]
					token.Act = str[1]
					_, exist := this.Dict[str[0]]
					if !exist {
						this.Dict[str[0]] = token
					}
				}
			} else {
				break
			}
		}
		dictFile.Close()
	}
	this.WriteWord()
}

func (this *Scaner) WriteWord() {
	fl, err := os.OpenFile(Out_File, os.O_APPEND|os.O_WRONLY, 0666)
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
	dictFile, err := os.Open("data/out1.txt")
	if err != nil {
		glog.Info("error:", err.Error())
	}
	defer dictFile.Close()
	reader := bufio.NewReader(dictFile)
	// 逐行读入分词
	var text string
	letterNum := 0
	num := 0
	existNum := 0
	for {
		size ,_ := fmt.Fscanln(reader, &text)
		glog.Info("size:", size, "  text:", text, " len:", len([]rune(text)))
		num++
		if size == 0 {
			// 文件结束
			break
		}
		if len([]rune(text)) == 1 {
			continue
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
			letterNum++
			continue
		}
		_, exist := this.Dict[text]
		if exist {
			existNum++
			continue
		}
		token := &Token{}
		token.Word = text
		this.Dict[text] = token
	}
	glog.Info("读取完毕 ", "num:", num, "  letterNum:",letterNum, "  existNum:",existNum)
	var keys []string
	for k := range this.Dict {
		keys = append(keys, k)
	}
	sort.Sort(sort.StringSlice(keys))
	fl, err := os.OpenFile("data/out2.txt", os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		glog.Info("error:", err.Error())
	}
	defer fl.Close()
	for _, str := range keys {
		str = strings.Replace(str, "[", "", -1)
		_, err1 := fl.Write([]byte( str + "\n"))
		if err1 != nil {
			glog.Info("error:", err1.Error())
		}
	}
}