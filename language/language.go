package language

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"stonesrv/conf"
	"stonesrv/log"
	"stonesrv/models"
	"sync"
)

var l = new(lang)

type lang struct {
	langMap sync.Map
	words   models.Words
}

//初始化语言
func Init() {
	path := fmt.Sprintf("language/%s/lang.xml", conf.GetLanguage())

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Load languages failed!")
		return
	}
	err = xml.Unmarshal(buf, &l.words)
	if err != nil {
		log.Fatal("Load languages failed!")
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
	return "N/A"
}
