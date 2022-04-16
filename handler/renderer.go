package handler

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"time"
)

const articleTemplate = "./html/article_layout.gohtml"

const calendarTemplate = "./html/calendar_layout.gohtml"
const manageTemplate = "./html/manage.gohtml"
const indexTemplate = "./html/index.gohtml"

type courseInfo struct {
	Subtitle string
	Course   []map[string]interface{}
}

type fileInfo struct {
	ClientName string
	ServerName string
	Mime       string
	Path       string
}

type articleRenderInfo struct {
	ID              int64
	User            string
	Type            string
	CreateTime      string
	PublishTime     string
	LastModified    string
	Title           string
	Content         template.HTML
	Attachment      []fileInfo
	PhotoAttachment []fileInfo // render photos by HTML <img>
}

type calendarRenderInfo struct {
	Calendar
	HaveLink bool
	ReadOnly bool
}

func convertStrm(i int) string {
	switch i {
	case 1:
		return "上"
	case 2:
		return "下"
	}
	return ""
}

func convertLevel(i int) string {
	switch i {
	case 1:
		return "一"
	case 2:
		return "二"
	case 3:
		return "三"
	case 4:
		return "四"
	}
	return ""
}

func renderDate(timestamp uint64) string {
	t := time.Unix(int64(timestamp), 0)
	return t.Format("2006-01-02")
	//(t.Format("2006-01-02 15:04:05"))
}

func renderArticleType(key string) string {
	dict := map[string]string{
		"normal":       "一般消息",
		"activity":     "演講 & 活動",
		"course":       "課程 & 招生",
		"scholarships": "獎學金",
		"recruit":      "徵才資訊",
	}

	val, ok := dict[key]
	if ok {
		return val
	}
	return ""
}

// RenderPublicArticle renders an article at url: /news/[articleID]
func RenderPublicArticle(artInfo *Article) template.HTML {
	data := new(articleRenderInfo)
	data.ID = artInfo.ID
	data.User = artInfo.User
	data.Type = renderArticleType(artInfo.Type)
	data.CreateTime = renderDate(artInfo.CreateTime)
	data.PublishTime = renderDate(artInfo.PublishTime)
	data.LastModified = renderDate(artInfo.LastModified)
	data.Title = artInfo.Title
	data.Content = template.HTML(artInfo.Content)
	data.PhotoAttachment = []fileInfo{}
	data.Attachment = []fileInfo{}

	for _, v := range artInfo.Attachment {
		var extName string
		matchNum, _ := fmt.Sscanf(v.Mime, "image/%s", &extName)
		if matchNum > 0 {
			data.PhotoAttachment = append(data.PhotoAttachment, fileInfo{
				Path: v.Path,
			})
		} else {
			data.Attachment = append(data.Attachment, fileInfo{
				ClientName: v.ClientName,
				ServerName: v.ServerName,
				Mime:       v.Mime,
				Path:       v.Path,
			})
		}
	}

	var buf bytes.Buffer
	t, err := template.ParseFiles(articleTemplate)

	if err != nil {
		log.Println("handler/renderer.go RenderPublicArticle() template error " + err.Error())
		return template.HTML("")
	}

	t.Execute(&buf, data)
	return template.HTML(buf.String())
}

// RenderPublicArticleBriefList dynamically renders article list an home page
func RenderPublicArticleBriefList(artInfoList []Article) template.HTML {
	data := new(articleRenderInfo)
	ret := ""
	for _, artInfo := range artInfoList {
		data.ID = artInfo.ID
		data.PublishTime = renderDate(artInfo.PublishTime)
		data.Title = artInfo.Title
		var buf bytes.Buffer
		t, _ := template.New("article_list_brief").Parse(`
        <div class="article">
            <div class="candy-header"><span class="single">{{.PublishTime}}</span></div>
            <a href="/news?id={{.ID}}">{{.Title}}</a>
        </div>`)
		t.Execute(&buf, data)
		ret += buf.String()
	}

	return template.HTML(ret)
}

func RenderIndexPage(lang int) template.HTML {
	t, err := template.ParseFiles(indexTemplate)
	if err != nil {
		log.Println("handler/renderer.go RenderIndexPage() template error " + err.Error())
		return template.HTML("")
	}

	// GetLatesetArticles() returns (list []Article, hasNext bool)
	artList, _ := GetLatesetArticles("public", "normal", "", 0, 7)
	var buf bytes.Buffer
	t.Execute(&buf, struct {
		ArticleListBrief template.HTML
		IsEn             bool
	}{
		ArticleListBrief: RenderPublicArticleBriefList(artList),
		IsEn:             lang == EN,
	})
	return template.HTML(buf.String())
}

func RenderMangePage(user *User) template.HTML {
	if user == nil {
		log.Println("handler/renderer.go RenderManagePage() variable user is nil")
		return template.HTML("")
	}

	t, err := template.ParseFiles(manageTemplate)
	if err != nil {
		log.Println("handler/renderer.go RenderManagePage() template error " + err.Error())
		return template.HTML("")
	}
	var buf bytes.Buffer

	t.Execute(&buf, user)
	return template.HTML(buf.String())
}
