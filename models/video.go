package models

import "encoding/json"

//Feed 主页订阅源
type Feed struct {
	Key             string  `json:"_key,omitempty"`
	IssueList       []Issue `json:"issueList"` //issueList
	NextPageURL     string  `json:"nextPageUrl"`
	NextPublishTime int64   `json:"nextPublishTime"`
	NewestIssueType string  `json:"newestIssueType"`
	Dialog          string  `json:"dialog"`
}

//Issue 内容源
type Issue struct {
	Key         string `json:"_key,omitempty"`
	ReleaseTime int64  `json:"releaseTime"`
	Type        string `json:"type"`
	Date        int64  `json:"date"`
	PublishTime int64  `json:"publishTime"`
	ItemList    []Item `json:"itemList"` //itemList
	Count       int64  `json:"count"`
}

//Item 内容项目
type Item struct {
	Key     string          `json:"_key,omitempty"`
	RawData json.RawMessage `json:"data"` //data
	Type    string          `json:"type"`
	Tag     string          `json:"tag"`
	ID      int64           `json:"id"`
	AdIndex int64           `json:"adIndex"`
}

//以下是各种data数据类型
//====== BEGIN ======

//BannerData 数据类型
type BannerData struct {
	Key         string   `json:"_key,omitempty"`
	DataType    string   `json:"dataType"`
	ID          int64    `json:"id"` //ID 同KEY
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Image       string   `json:"image"`
	ActionURL   string   `json:"actionUrl"`
	AdTrack     string   `json:"adTrack"`
	Shade       bool     `json:"shade"`
	Label       string   `json:"label"`
	LabelList   []string `json:"labelList"`
	Header      string   `json:"header"`
	AutoPlay    bool     `json:"autoPlay"`
}

//TextHeaderData 数据类型
type TextHeaderData struct {
	Key      string `json:"_key,omitempty"`
	DataType string `json:"dataType"`
	Text     string `json:"text"`
	Font     string `json:"font"`
	AdTrack  string `json:"adTrack"`
}

//VideoData 数据类型
type VideoData struct {
	Key               string      `json:"_key,omitempty"`
	DataType          string      `json:"dataType"`
	ID                int64       `json:"id"` //ID 同KEY
	Title             string      `json:"title"`
	Description       string      `json:"description"`
	Library           string      `json:"library"`
	Tags              []Tag       `json:"tags"`        //tag
	Consum            Consumption `json:"consumption"` //consumption
	ResourceType      string      `json:"resourceType"`
	Slogan            string      `json:"slogan"`
	Provi             Provider    `json:"provider"` //provider
	Category          string      `json:"category"`
	Auth              Author      `json:"author"` //author
	Cov               Cover       `json:"cover"`  //cover
	PlayURL           string      `json:"playUrl"`
	ThumbPlayURL      string      `json:"thumbPlayUrl"`
	Duration          int64       `json:"duration"`
	Web               WebUrl      `json:"webUrl"` //webUrl
	ReleaseTime       int64       `json:"releaseTime"`
	PlInfo            []PlayInfo  `json:"playInfo"` //playInfo
	Campaign          string      `json:"campaign"`
	WaterMarks        string      `json:"waterMarks"`
	AD                bool        `json:"ad"`
	AdTrack           string      `json:"adTrack"`
	Type              string      `json:"type"`
	TitlePgc          string      `json:"titlePgc"`
	DescriptionPgc    string      `json:"descriptionPgc"`
	Remark            string      `json:"remark"`
	IfLimitVideo      bool        `json:"ifLimitVideo"`
	SearchWeight      int64       `json:"searchWeight"`
	IDx               int64       `json:"idx"`
	ShareAdTrack      string      `json:"shareAdTrack"`
	FavoriteAdTrack   string      `json:"favoriteAdTrack"`
	WEBAdTrack        string      `json:"webAdTrack"`
	Date              int64       `json:"date"`
	Promotion         string      `json:"promotion"`
	Label             string      `json:"label"`
	LabelList         []string    `json:"labelList"`
	DescriptionEditor string      `json:"descriptionEditor"`
	Collected         bool        `json:"collected"`
	Played            bool        `json:"played"`
	Subtitles         []string    `json:"subtitles"`
	LastViewTime      string      `json:"lastViewTime"`
	Playlists         []string    `json:"playlists"`
	Src               string      `json:"src"`
}

//====== 以下类型均为VIDEO数据类型中的子类型 ======

//Tag 标记
type Tag struct {
	Key            string `json:"_key,omitempty"`
	ID             int64  `json:"id"` //ID 同KEY
	Name           string `json:"name"`
	ActionURL      string `json:"actionUrl"`
	AdTrack        string `json:"adTrack"`
	Desc           string `json:"desc"`
	BGPicture      string `json:"bgPicture"`
	HeaderImage    string `json:"headerImage"`
	TagRecType     string `json:"tagRecType"`
	ChildTagList   string `json:"childTagList"`
	ChildTagIDList string `json:"childTagIdList"`
	CommunityIndex int64  `json:"communityIndex"`
}

//Consumption 使用统计
type Consumption struct {
	Key             string `json:"_key,omitempty"`
	CollectionCount int64  `json:"collectionCount"`
	ShareCount      int64  `json:"shareCount"`
	ReplyCount      int64  `json:"replyCount"`
}

//Provider 提供者
type Provider struct {
	Key   string `json:"_key,omitempty"`
	Name  string `json:"name"`
	Alias string `json:"alias"`
	Icon  string `json:"icon"`
}

//Author 作者
type Author struct {
	Key                        string `json:"_key,omitempty"`
	ID                         int64  `json:"id"` //ID 同KEY
	Icon                       string `json:"icon"`
	Name                       string `json:"name"`
	Description                string `json:"description"`
	Link                       string `json:"link"`
	LatestReleaseTime          int64  `json:"latestReleaseTime"`
	VideoNum                   int64  `json:"videoNum"`
	AdTrack                    string `json:"adTrack"`
	Foll                       Follow `json:"follow"`
	Shie                       Shield `json:"shield"`
	ApprovedNotReadyVideoCount int64  `json:"approvedNotReadyVideoCount"`
	IfPgc                      bool   `json:"ifPgc"`
	RecSort                    int64  `json:"recSort"`
	Expert                     bool   `json:"expert"`
}

//Follow 关注
type Follow struct {
	Key      string `json:"_key,omitempty"`
	ItemType string `json:"itemType"`
	ItemID   int64  `json:"itemId"`
	Followed bool   `json:"followed"`
}

//Shield 喜欢？
type Shield struct {
	Key      string `json:"_key,omitempty"`
	ItemType string `json:"itemType"`
	ItemID   int64  `json:"itemId"`
	Shielded bool   `json:"shielded"`
}

//Cover 视频封面
type Cover struct {
	Key      string `json:"_key,omitempty"`
	Feed     string `json:"feed"`
	Detail   string `json:"detail"`
	Blurred  string `json:"blurred"`
	Sharing  string `json:"sharing"`
	Homepage string `json:"homepage"`
}

//WebUrl 网页的连接 可以不用
type WebUrl struct {
	Key      string `json:"_key,omitempty"`
	Raw      string `json:"raw"`
	ForWeibo string `json:"forWeibo"`
}

//PlayInfo 播放信息，标清/超清等，可以只用一种类型
type PlayInfo struct {
	Key     string `json:"_key,omitempty"`
	Height  int64  `json:"height"`
	Width   int64  `json:"width"`
	URLList []Url  `json:"urlList"` //Url
	Name    string `json:"name"`
	Type    string `json:"type"`
	URL     string `json:"url"`
}

//Url 超链接
type Url struct {
	Key  string `json:"_key,omitempty"`
	Name string `json:"name"`
	URL  string `json:"url"`
	Size int64  `json:"size"`
}

//====== 以上类型均为VIDEO数据类型中的子类型 ======

//====== END ======
