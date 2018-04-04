package segment

import (
	"path/filepath"
	"os"
	"github.com/golang/glog"
	"github.com/axgle/mahonia"
	"bufio"
	"io"
	"strings"
	"unicode/utf8"
	"unicode"
	"fmt"
)

type Ngram struct {
	Interpunction  map[string]bool
	Dict   map[string]*Word
	Sum    uint64
}

type Word struct {
	Text   string	   // 文字内容
	Frequencys float64  // 出现的频率
	Time   uint64    // 出现的次数
}
func (this *Ngram) Init() {
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

	this.Dict = make(map[string]*Word)
	this.Sum = 0

}


func (this *Ngram) isInterpution(str string) bool{
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

// 目录文件  输出文件
func (this *Ngram) Count(dirPath string, outPath string) {
	paths := []string{}
	filepath.Walk(dirPath, func(filename string, fi os.FileInfo, err error) error {
		//遍历目录
		if fi.IsDir() {
			// 忽略目录
			return nil
		}
		paths = append(paths, filename)
		return nil
	})
	glog.Info("paths:", len(paths))

	for _, path := range paths {
		str := this.read(path)
		this.record(str)
	}

	for _, v := range this.Dict {
		this.Sum += v.Time
	}
	for _, v := range this.Dict {
		v.Frequencys = float64(v.Time)/float64(this.Sum)
		glog.Info(v.Frequencys)
	}
	this.write(outPath)
	glog.Info("sum:", this.Sum)
}

func (this *Ngram) record(str string) {
	//glog.Info("record:", str)
	bytes := []byte(str)
	current := 0
	for current < len(bytes) {
		r, size:= utf8.DecodeRune(bytes[current:])
		current += size
		if size <= 2 && (unicode.IsLetter(r) || unicode.IsNumber(r)) {
			continue
		} else {
			if this.isInterpution(string(r)) || r == 12288 || len(strings.TrimSpace(string(r))) <= 0 {
				continue
			} else {
				_, exist := this.Dict[string(r)]
				if exist {
					this.Dict[string(r)].Time++
				} else {
					word := &Word{}
					word.Time = 1
					word.Text = string(r)
					this.Dict[string(r)] = word
				}
 			}
		}
	}
}

func (this *Ngram) read(path string) string {
	glog.Info("read:", path)
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE,0766)
	if err != nil {
		glog.Info("err:", err.Error(), " path:", path)
	}
	decoder := mahonia.NewDecoder("gbk")
	reader := bufio.NewReader(decoder.NewReader(file))
	var str string
	for {
		var tmp string
		tmp,err:= reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		str += string(tmp)
	}
	file.Close()
	return str
}

func (this *Ngram) write(path string) {
	glog.Info("write:", path)
	fl, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE,0766)
	if err != nil {
		glog.Info("error:", err.Error())
		return
	}
	sum := ""
	for _, item := range this.Dict {
		str := fmt.Sprintf("%s %d\n", item.Text, item.Time)
		sum = sum + str
	}
	_, err1 := fl.Write([]byte(sum))
	if err1 != nil {
		glog.Info("error:", err1.Error())
	}
}