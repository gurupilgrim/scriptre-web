package main

import (
	"bufio"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var templates = template.Must(template.ParseFiles("index.html"))
var booknames [][3]string

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
	//preload data
	initialize()
	//handlers
	http.HandleFunc("/", handleWeb)
	http.HandleFunc("/v0/", handleAPIBeta)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))

	log.Fatal(http.ListenAndServe(":80", nil))
}

func initialize() {
	//determine which canon to use
	canon := "protestant"
	//get or have list of books in canon
	canonFilename := fmt.Sprintf("./index/%s.canon", canon)
	file, _ := os.Open(canonFilename)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var bookInfo [3]string
		rawBookInfo := strings.Split(scanner.Text(), ";")
		bookInfo[0] = rawBookInfo[0]
		bookInfo[1] = rawBookInfo[1]
		bookInfo[2] = rawBookInfo[2]
		booknames = append(booknames, bookInfo)
	}

}
