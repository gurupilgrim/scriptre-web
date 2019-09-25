package main

import (
	//	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

var templates = template.Must(template.ParseFiles("index.html"))

type Thingy struct {
	Title string
}

func handleWeb(w http.ResponseWriter, r *http.Request) {
	thing := Thingy{
		Title: "test",
	}
	page := template.Must(template.ParseFiles("index.html", "quick_nav.html"))
	page.Execute(w, thing)
	//if err != nil {
	//	fmt.Fprintf(w, "%v", err.Error())
	//}
	//_ = templates.ExecuteTemplate(w, "index.html", page)
}

func loadFile(file string) (*[]byte, error) {
	fileContents, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return &fileContents, nil
}

func main() {
	http.HandleFunc("/", handleWeb)
	http.HandleFunc("/v0/", handleAPIBeta)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))

	log.Fatal(http.ListenAndServe(":80", nil))
}
