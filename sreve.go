package main

import (
	"bufio"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"unicode"
)

var templates = template.Must(template.ParseFiles("index.html"))
var booknames [][3]string
var bibleData []Bible

type Bible struct {
	Version string
	Books   []Book
}

type Book struct {
	Name     string
	Chapters []Chapter
}

type Chapter struct {
	Number int
	Verses []Verse
}

type Verse struct {
	Reference  Reference
	Language   string
	Subsection string
	Text       string
}

type VerseGroup struct {
	Verses []Verse
}

type Reference struct {
	BibleVersion string
	Book         string
	Chapter      int
	VerseNumber  int
	VerseRange   int
	Canon        string
}

type Thingy struct {
	Title string
}

func handleWeb(w http.ResponseWriter, r *http.Request) {
	thing := Thingy{
		Title: "test",
	}
	page := template.Must(template.ParseFiles("index.html", "quick_nav.html"))
	page.Execute(w, thing)
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
	//get a list of books in canon
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

	//preload bibles
	//currently just english, possibly use languages as categories
	enBibleNames, _ := ioutil.ReadDir("./bibles/en/")

	//load bibles in to memory
	for _, bible := range enBibleNames {
		books, _ := ioutil.ReadDir(fmt.Sprintf("./bibles/en/%s/", bible.Name()))
		thisBible := new(Bible)
		thisBible.Version = bible.Name()
		for _, bookFile := range books {
			rawBook, _ := os.Open(fmt.Sprintf("./bibles/en/%s/%s", bible.Name(), bookFile.Name()))
			scanner = bufio.NewScanner(rawBook)

			var thisBook Book
			thisBook.Name = strings.TrimSuffix(strings.Replace(fmt.Sprintf("%s", bookFile.Name()), "_", " ", -1), ".sfm")
			var curChapter = new(Chapter)
			for scanner.Scan() {
				if strings.HasPrefix(scanner.Text(), "\\c") {
					if curChapter.Number > 0 {
						thisBook.Chapters = append(thisBook.Chapters, *curChapter)
					}
					number := scanner.Text()[3:]
					chapter := new(Chapter)
					chapter.Number, _ = strconv.Atoi(number)
					curChapter = chapter
				} else if strings.HasPrefix(scanner.Text(), "\\v") {
					thisVerse := new(Verse)
					thisReference := new(Reference)
					var verseNumberString string
					var textIndex int
					textIndex = 5
					verseNumberString = string(scanner.Text()[3])
					if unicode.IsNumber(rune(scanner.Text()[4])) {
						textIndex = 6
						verseNumberString = fmt.Sprintf("%s%s", verseNumberString, string(scanner.Text()[4]))
					}
					if unicode.IsNumber(rune(scanner.Text()[5])) {
						textIndex = 7
						verseNumberString = fmt.Sprintf("%s%s", verseNumberString, string(scanner.Text()[5]))
					}
					thisReference.BibleVersion = thisBible.Version
					thisReference.Book = thisBook.Name
					thisReference.Chapter = curChapter.Number
					thisReference.VerseNumber, _ = strconv.Atoi(verseNumberString)
					thisReference.Canon = "protestant"
					thisVerse.Reference = *thisReference
					thisVerse.Text = scanner.Text()[textIndex:]
					thisVerse.Language = "en"
					thisVerse.Subsection = "all"
					curChapter.Verses = append(curChapter.Verses, *thisVerse)
				}

			}
			thisBook.Chapters = append(thisBook.Chapters, *curChapter)
			thisBible.Books = append(thisBible.Books, thisBook)
		}
		bibleData = append(bibleData, *thisBible)
	}

	fmt.Printf("bibleData: %+v\n", bibleData)
	fmt.Printf("finished initializing.\n")
}
