package handler

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"regexp"
)

const (
	ZH = iota
	EN
)

func getHTML(fileName string) (template.HTML, error) {
	t, err := template.ParseFiles("./html" + fileName + ".html")
	if err != nil {
		log.Println(err)
		return template.HTML("error try again"), err
	}

	var buf bytes.Buffer
	t.Execute(&buf, nil)
	return template.HTML(buf.String()), nil
}

func getHTMLwithLang(fileName string, lang int) (template.HTML, error) {
	t, err := template.ParseFiles("./html" + fileName + ".html")
	if err != nil {
		log.Println(err)
		return template.HTML("error try again"), err
	}

	var buf bytes.Buffer
	t.Execute(&buf, struct {
		IsEn bool
	}{
		IsEn: lang == EN,
	})
	return template.HTML(buf.String()), nil
}

func lang(r *http.Request) int {
	if r.Method == "GET" {
		cookieLang, _ := r.Cookie("lang")
		if cookieLang != nil {
			switch cookieLang.Value {
			case "zh":
				return ZH
			case "en":
				return EN
			}
		}

		headerLang := r.Header.Get("Accept-Language")
		if is, _ := regexp.MatchString("^zh.*", headerLang); is {
			return ZH
		}
	}
	return EN
}
