package handler

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// PageData is a type for filling HTML template
type PageData struct {
	Title   string
	Isindex bool
	IsLogin bool
	Main    template.HTML
	Time    int64
	Year    int
	IsEn    bool
	IsZh    bool
}

func initPageData() *PageData {
	data := new(PageData)
	data.Isindex = false // default value
	data.Time = time.Now().Unix() >> 10
	data.Year, _, _ = time.Now().Date()
	return data
}

// BasicWebHandler is a handler for handling url whose prefix is /
func BasicWebHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// Handle static file
	staticFiles := []string{"/robots.txt", "/sitemap.xml", "/favicon.ico"}
	for _, f := range staticFiles {
		if r.URL.Path == f {
			http.StripPrefix("/", http.FileServer(http.Dir("./"))).ServeHTTP(w, r)
			return
		}
	}

	// Handle simple web and check login info
	data := initPageData()

	var twisc string
	// Get language
	langCode := lang(r)
	if langCode == ZH {
		data.IsZh = true
		twisc = "國立中興大學資通安全研究與教學中心"
	} else if langCode == EN {
		data.IsEn = true
		twisc = "Taiwan Information Security Center @NCHU"
	}

	user := CheckLoginBySession(w, r)
	data.IsLogin = user != nil

	var simpleWeb = map[string]map[int]string{
		"/about": {
			ZH: "關於",
			EN: "About",
		},
		"/development": {
			ZH: "技術研發",
			EN: "Development",
		},
		"/events": {
			ZH: "學術活動",
			EN: "Events",
		},
		"/research": {
			ZH: "研究成果",
			EN: "Research Findings",
		},
		"/official-document": {
			ZH: "辦法表格",
			EN: "Official Document",
		},
		"/privacy-statement": {
			ZH: "隱私權聲明",
			EN: "Privacy Statement",
		},
		"/ncsu-cloud": {
			ZH: "資通安全弱點通報機制",
		},
	}

	// Handle simple web
	if title, ok := simpleWeb[r.URL.Path]; ok {
		data.Title = title[langCode]
	} else {
		// Handle non simple web
		switch r.URL.Path {
		case "/":
			data.Title = twisc
			data.Isindex = true
			data.Main = RenderIndexPage(langCode)
		case "/news":
			if langCode == ZH {
				data.Title = "最新消息"
			} else {
				data.Title = "News"
			}

			if id := strings.Join(r.Form["id"], ""); id != "" {
				aid, err := strconv.ParseInt(id, 10, 64)

				if err != nil {
					NotFound(w, r)
					return
				}

				uid := ""
				if data.IsLogin {
					uid = user.ID
				}

				artInfo := GetArticleByAid(aid, uid)

				// avoid /news?id=xxx
				if artInfo == nil {
					NotFound(w, r)
					return
				}

				data.Title = artInfo.Title + " | " + twisc
				data.Main = RenderPublicArticle(artInfo)
			} else {
				data.Title += " | " + twisc
				data.Main, _ = getHTML(r.URL.Path)
			}
		case "/login":
			if CheckLoginBySession(w, r) != nil {
				http.Redirect(w, r, "/manage", 302)
				return
			}
			data.Title = "登入"
		case "/logout":
			ret := struct {
				Err string `json:"err"`
			}{}
			if err := Logout(w, r); err != nil {
				ret.Err = "登出失敗，重試，或使用瀏覽器清除 Cookie"
				json.NewEncoder(w).Encode(w)
				return
			}

			http.Redirect(w, r, "/", 302)
			return
		default:
			NotFound(w, r)
			return
		}
	}

	if r.URL.Path != "/" && r.URL.Path != "/news" {
		data.Title += " | " + twisc
		data.Main, _ = getHTMLwithLang(r.URL.Path, langCode)
	}

	t, _ := template.ParseFiles("./html/layout.gohtml")
	t.Execute(w, data)
}
