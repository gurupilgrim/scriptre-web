package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
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
	}
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
	fmt.Printf("\n\n")
	return
}

func narrowBook(canon string, query []byte, startingNarrow []int) (int, []int) {

	var selectedBook int
	var bookIsSelected bool

	//get the *index* of the last character of the query using len
	lenIndex := len(query) - 1
	fmt.Printf("query is %v long\n", len(query))
	//compare the *last* character in query with the indexnth character of each remaining book in startingNarrow
	//create a new []int to include the newNarrow
	var newNarrow []int
	//if staringNarrow is empty it will only run the else here
	if len(startingNarrow) > 0 {
		for i := range startingNarrow {
			if unicode.ToLower(rune(query[lenIndex])) == unicode.ToLower(rune(booknames[startingNarrow[i]][1][lenIndex])) {
				curBookWeight, _ := strconv.Atoi(booknames[startingNarrow[i]][2])
				selBookWeight, _ := strconv.Atoi(booknames[selectedBook][2])

				newNarrow = append(newNarrow, startingNarrow[i])
				if bookIsSelected == false {
					selectedBook = startingNarrow[i]
					bookIsSelected = true
				} else if curBookWeight > selBookWeight {
					selectedBook = startingNarrow[i]
				}
			}
		}
	} else {
		for i, bookname := range booknames {
			if unicode.ToLower(rune(query[lenIndex])) == unicode.ToLower(rune(bookname[1][lenIndex])) {
				curBookWeight, _ := strconv.Atoi(booknames[i][2])
				selBookWeight, _ := strconv.Atoi(booknames[selectedBook][2])

				newNarrow = append(newNarrow, i)
				if bookIsSelected == false {
					selectedBook = i
					bookIsSelected = true
				} else if curBookWeight > selBookWeight {
					selectedBook = i
				}
			}
		}
	}
	return selectedBook, newNarrow

}

func narrowPrefix(canon string, prefix int) (int, []int) {
	var newNarrow []int
	var bookIsSelected bool
	var selectedBook int

	if prefix == 1 {
		for i, bookname := range booknames {
			if bookname[0] == "I" {
				newNarrow = append(newNarrow, i)
				curBookWeight, _ := strconv.Atoi(booknames[i][2])
				selBookWeight, _ := strconv.Atoi(booknames[selectedBook][2])

				if bookIsSelected == false {
					selectedBook = i
					bookIsSelected = true
				} else if curBookWeight > selBookWeight {
					selectedBook = i
				}
			}
		}
	}

	if prefix == 2 {
		for i, bookname := range booknames {
			if bookname[0] == "II" {
				newNarrow = append(newNarrow, i)
				curBookWeight, _ := strconv.Atoi(booknames[i][2])
				selBookWeight, _ := strconv.Atoi(booknames[selectedBook][2])

				if bookIsSelected == false {
					selectedBook = i
					bookIsSelected = true
				} else if curBookWeight > selBookWeight {
					selectedBook = i
				}
			}
		}
	}

	if prefix == 3 {
		for i, bookname := range booknames {
			if bookname[0] == "III" {
				newNarrow = append(newNarrow, i)
				curBookWeight, _ := strconv.Atoi(booknames[i][2])
				selBookWeight, _ := strconv.Atoi(booknames[selectedBook][2])

				if bookIsSelected == false {
					selectedBook = i
					bookIsSelected = true
				} else if curBookWeight > selBookWeight {
					selectedBook = i
				}
			}
		}
	}

	return selectedBook, newNarrow
}

func getRef(request string, canon string) (Reference, error) {

	var result Reference
	//get the right canon
	fmt.Printf("encoded request [[%v]]\n", request)
	request, _ = url.QueryUnescape(request)
	fmt.Printf("decoded request [[%v]]\n", request)

	//clean up the string a bit by removing anything other than recognized characters
	var firstLetterIndex int
	for i, character := range []rune(request) {
		if unicode.IsLetter(character) {
			//we found a letter
			fmt.Printf("found a letter in position %v\n", i)
			firstLetterIndex = i
			break
		}
	}
	var firstNumberIndex int
	for i, character := range []rune(request) {
		if unicode.IsNumber(character) {
			//we found a number
			fmt.Printf("found a number in position %v\n", i)
			firstNumberIndex = i
			break
		}
	}

	var narrow []int

	//determine if we have a prefix (1, I, First, etc)
	//do we have a number before a letter?
	if firstNumberIndex < firstLetterIndex {
		firstNumberString := string(request[firstNumberIndex])
		firstNumber, _ := strconv.Atoi(firstNumberString)
		fmt.Printf("looks like we have a number prefix: %v\n", firstNumber)

		if firstNumber < 4 {
			fmt.Printf("number is less than 4\n")
			if !(unicode.IsNumber(rune(request[firstNumberIndex+1]))) {
				fmt.Printf("we have just one number. passing %v on\n", firstNumber)
				_, narrow = narrowPrefix(canon, firstNumber)
				fmt.Printf("narrow is now %v\n", narrow)
			}
		}
	}
	//if not, is the first letter an i?
	//if so, what is the i followed by?
	//if not, is the first letter an f?
	//if so, is it the word first?
	//if not, is the first letter an s?
	//if so, is it the word second?
	//if not, is the first letter a t?
	//if so, is it the word third?

	//now we need to narrow this down to a specific book
	var reqWord []byte
	reqWord = []byte(fmt.Sprintf("%s", request[firstLetterIndex:]))

	var bookNameIndex int
	for i, _ := range []rune(request) {
		fmt.Printf("sending %v to narrow with length of %v\n", reqWord[:i+1], len(reqWord[:i+1]))
		curBook, newNarrow := narrowBook(canon, reqWord[:i+1], narrow)
		narrow = newNarrow
		bookNameIndex = curBook
		if len(narrow) < 2 {
			break
		}
		i = i + 1
		if i == (len(request) - 1) {
			break
		}
	}
	fmt.Printf("Final Bookname: %v\n", booknames[bookNameIndex])
	//req := unicode.ToLower(request[indexoffirstletter])
	//compare that to a slice of strings that represents the selected canon
	//get the index of each match for the first letter
	//for each index that was a match, look for the second letter
	//etc
	//if a book match is found, start looking for a reference within it

	//look up how many chapters and verses in each chapter there are for that book
	//parse request to see what reference numbers are being looked for
	//
	/*
		if booknames[bookNameIndex][0] != "" {
			result.Book = booknames[bookNameIndex][1]
		} else {
			result.Book = fmt.Sprintf("%s %s", booknames[bookNameIndex][0], booknames[bookNameIndex][1])
		}
		return result, nil
	*/
	result.Book = fmt.Sprintf("%s %s", booknames[bookNameIndex][0], booknames[bookNameIndex][1])
	return result, nil
}
