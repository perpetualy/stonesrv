package language

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"stonesrv/conf"
	"stonesrv/log"
	"stonesrv/models"
	"strings"
	"sync"
)

var l = new(lang)

type lang struct {
	langMap sync.Map
	words   models.Words
}

//Init 初始化语言
func Init(path string) {
	languatepath := fmt.Sprintf("./language/%s/lang.xml", conf.GetLanguage())
	if strings.Compare(path, "") != 0 {
		languatepath = fmt.Sprintf("%s/%s/lang.xml", path, conf.GetLanguage())
	}
	buf, err := ioutil.ReadFile(languatepath)
	if err != nil {
		log.Panic("Load languages failed!")
		return
	}
	err = xml.Unmarshal(buf, &l.words)
	if err != nil {
		log.Panic("Load languages failed!")
		return
	}
	for _, w := range l.words.Word {
		l.langMap.Store(w.Code, w.Text)
	}
}

//GetText 获取文本
func GetText(code int) string {
	if itext, ok := l.langMap.Load(code); ok {
		return itext.(string)
	}
	return "NONE"
}
