package main 
import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"unicode"
)

type Reference struct {
	BibleVersion string
	Book         string
	Chapter      int
	VerseNumber  int
	Canon        string
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

func handleAPIBeta(w http.ResponseWriter, r *http.Request) {
	// get the query
	url := strings.Split(r.URL.String(), "?")
	urlparams := strings.Split(url[1], "&")
	params := make(map[string]string)
	for _, param := range urlparams {
		thisParam := strings.Split(param, "=")
		params[thisParam[0]] = thisParam[1]
		fmt.Printf("%v\n", thisParam[1])
	}
	fmt.Printf("%v\n\n", params)
	ref, _ := getRef(params["query"], "protestant")
	fmt.Printf("Found Reference: %v\n", ref)
	verse := &Verse{
		Reference: Reference{
			BibleVersion: "KJV",
			Book:         ref.Book,
			Chapter:      1,
			VerseNumber:  1,
			Canon:        "protestant",
		},
		Language:   "en",
		Subsection: "all",
		Text:       r.URL.String(),
	}
	jsonVerse, _ := json.Marshal(verse)
	fmt.Fprintf(w, "%s", jsonVerse)
	return
}

func narrowBook(canon string, query []byte, startingNarrow []int) (string, []int) {

	var selectedBookName string

	//get or have list of books in canon
	canonFilename := fmt.Sprintf("./index/%s.canon", canon)
	file, _ := os.Open(canonFilename)
	scanner := bufio.NewScanner(file)

	var booknames []string

	for scanner.Scan() {
		booknames = append(booknames, scanner.Text())
	}

	fmt.Printf("%v\n", booknames)
	//get the *index* of the last character of the query using len
	lenIndex := len(query) - 1
	fmt.Printf("Query length: %v using query index: %v\n", len(query), lenIndex)
	//compare the *last* character in query with the indexnth character of each remaining book in startingNarrow
	fmt.Printf("Using: %v\n", startingNarrow)
	//create a new []int to include the newNarrow
	var newNarrow []int
	if len(startingNarrow) > 0 {
		for i := range startingNarrow {
			fmt.Printf("comparing %s to %s in %s\n", unicode.ToLower(rune(query[lenIndex])), unicode.ToLower(rune(booknames[startingNarrow[i]][lenIndex])), booknames[startingNarrow[i]])
			if unicode.ToLower(rune(query[lenIndex])) == unicode.ToLower(rune(booknames[startingNarrow[i]][lenIndex])) {
				fmt.Printf("found match!\n")
				newNarrow = append(newNarrow, startingNarrow[i])
				selectedBookName = booknames[startingNarrow[i]]
			}
		}
	} else {
		for i, bookname := range booknames {
			//TODO skip if not in startingNarrow
			if unicode.ToLower(rune(query[lenIndex])) == unicode.ToLower(rune(bookname[lenIndex])) {
				newNarrow = append(newNarrow, i)
				selectedBookName = bookname
			}
		}
	}
	fmt.Printf("Selected Book: %v\n", selectedBookName)
	//return empty string if int slice contains more than one item
	if len(newNarrow) > 1 {
		return "", newNarrow
	}
	var emptyintvar []int
	return selectedBookName, emptyintvar

}

func getRef(request string, canon string) (Reference, error) {

	var result Reference

	//parse the string to see if there is a word somewhere
	var firstLetterIndex int
	for i, character := range []rune(request) {
		if unicode.IsLetter(character) {
			//we found a letter
			fmt.Printf("found a letter in position %v\n", i)
			firstLetterIndex = i
			break
		}
	}
	//once we know the first place a letter occurs, pass this information on to narrow
	var reqWord []byte
	reqWord = []byte(fmt.Sprintf("%s", request[firstLetterIndex:]))
	//canon
	//first character of the first word
	//empty int slice
	var narrow []int
	var bookName string
	var i int
	i = 1
	for {
		fmt.Printf("narrow before: %v\n", narrow)
		fmt.Printf("reqWord: %v\n", reqWord[:i])
		curBookName, newNarrow := narrowBook(canon, reqWord[:i], narrow)
		narrow = newNarrow
		fmt.Printf("narrow after: %v\n", narrow)
		if len(narrow) < 2 {
			bookName = curBookName
			break
		}
		i = i + 1
		fmt.Printf("\n\n")
	}
	fmt.Printf("Final Bookname: %v\n", bookName)
	//req := unicode.ToLower(request[indexoffirstletter])
	//compare that to a slice of strings that represents the selected canon
	//get the index of each match for the first letter
	//for each index that was a match, look for the second letter
	//etc
	//if a book match is found, start looking for a reference within it

	//look up how many chapters and verses in each chapter there are for that book
	//parse request to see what reference numbers are being looked for
	//
	result.Book = bookName
	return result, nil
}
