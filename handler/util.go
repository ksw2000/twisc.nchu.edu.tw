package handler

import (
	"bytes"
	"html/template"
	"log"
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
