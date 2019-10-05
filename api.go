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

func handleAPIBeta(w http.ResponseWriter, r *http.Request) {
	// get the query
	url := strings.Split(r.URL.String(), "?")
	urlparams := strings.Split(url[1], "&")
	// read the url to see what action to take
	params := make(map[string]string)
	for _, param := range urlparams {
		thisParam := strings.Split(param, "=")
		params[thisParam[0]] = thisParam[1]
	}
	// get reference
	ref, _ := getRef(params["query"], "protestant")
	// get verse
	verse, _ := getVerse(ref)
	jsonVerse, _ := json.Marshal(verse)
	fmt.Fprintf(w, "%s", jsonVerse)
	return
}

func narrowBook(canon string, query []byte, startingNarrow []int) (int, []int) {

	var selectedBook int
	var bookIsSelected bool

	//get the *index* of the last character of the query using len
	lenIndex := len(query) - 1
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

// Narrows the possible books a reference could be by the prefix
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

// verifies the remaining bits of a string match a selected book name
func verifyBook(canon string, query []byte, selectedBook int) bool {
	//query is expected to not include the prefix or the verse number

	var verified bool

	//get selected book
	bookname := booknames[selectedBook][1]

	//look for a match
	for i, character := range query {
		if unicode.ToLower(rune(character)) == unicode.ToLower(rune(bookname[i])) {
			verified = true
		} else {
			verified = false
			break
		}
	}
	return verified
}

// generates a reference object based on a query
func getRef(request string, canon string) (Reference, error) {

	var result Reference
	//get the right canon
	request, _ = url.QueryUnescape(request)

	//clean up the string a bit by removing anything other than recognized characters
	var firstLetterIndex int
	var firstNumberIndex int
	var lastLetterIndex int
	var queryBooknameStart int

	var firstLetterIdentified bool
	var firstNumberIdentified bool

	// determine the parts of the query. first letter, last letter, etc.
	for i, character := range []rune(request) {
		if unicode.IsLetter(character) {
			if firstLetterIdentified {
				lastLetterIndex = i
			} else {
				firstLetterIndex = i
				queryBooknameStart = i
				firstLetterIdentified = true
				lastLetterIndex = i
			}
		}
		if unicode.IsNumber(character) && !firstNumberIdentified {
			firstNumberIndex = i
			firstNumberIdentified = true
		}
	}

	var narrow []int
	//var hasPrefix bool

	//determine if we have a prefix (1, I, First, etc)
	//do we have a number before a letter?
	if firstNumberIndex < firstLetterIndex {
		firstNumberString := string(request[firstNumberIndex])
		firstNumber, _ := strconv.Atoi(firstNumberString)

		if firstNumber < 4 {
			if !(unicode.IsNumber(rune(request[firstNumberIndex+1]))) {
				_, narrow = narrowPrefix(canon, firstNumber)
			}
		}
		//if not, is the first letter an i?
	} else if unicode.ToLower(rune(request[firstLetterIndex])) == []rune("i")[0] {
		//if so, what is the i followed by?
		//loop through subsequent characters to determine a course of action
		for i, char := range request {
			character := unicode.ToLower(rune(char))
			if i > 0 {
				if character == []rune(" ")[0] {
					//hasPrefix = true
					if strings.ToLower(request[:i]) == "i" {
						_, narrow = narrowPrefix(canon, 1)
						queryBooknameStart = i + 1
						break
					} else if strings.ToLower(request[:i]) == "ii" {
						_, narrow = narrowPrefix(canon, 2)
						queryBooknameStart = i + 1
						break
					} else if strings.ToLower(request[:i]) == "iii" {
						_, narrow = narrowPrefix(canon, 3)
						queryBooknameStart = i + 1
						break
					}
				} else if character == []rune("i")[0] {
					if i < 3 {
						continue
					}
				} else {
					break
				}
			}
		}
		//TODO try narrowing by name, if the next character doesn't match a book name, treat this as a prefix
	}
	//TODO if not, is the first letter an f?
	//TODO if so, is it the word first?
	//TODO if not, is the first letter an s?
	//TODO if so, is it the word second?
	//TODO if not, is the first letter a t?
	//TODO if so, is it the word third?

	//now we need to narrow this down to a specific book
	var reqWord []byte
	reqWord = []byte(fmt.Sprintf("%s", request[queryBooknameStart:]))

	var bookNameIndex int
	// narrow the query down to a book
	for i, _ := range reqWord {
		curBook, newNarrow := narrowBook(canon, reqWord[:i+1], narrow)
		if len(newNarrow) == 1 {
			bookNameIndex = curBook
			break
		}
		if len(newNarrow) < 1 {
			break
		}
		//we only get here if newNarrow is greater than 1
		bookNameIndex = curBook
		narrow = newNarrow
		i = i + 1
	}
	//req := unicode.ToLower(request[indexoffirstletter])
	//compare that to a slice of strings that represents the selected canon
	//get the index of each match for the first letter
	//for each index that was a match, look for the second letter
	//etc

	//verify the bookname matches the rest of the query
	//clean up the end of the string by removing anything after the last character
	queryBookname := request[queryBooknameStart : lastLetterIndex+1]
	verified := verifyBook(canon, []byte(queryBookname), bookNameIndex)
	if verified {
		if booknames[bookNameIndex][0] != "" {
			result.Book = fmt.Sprintf("%s %s", booknames[bookNameIndex][0], booknames[bookNameIndex][1])
		} else {
			result.Book = fmt.Sprintf("%s", booknames[bookNameIndex][1])
		}
		//if a book match is found, start looking for a reference within it

		var chapterNumberGroup []rune
		var chapterNumber int
		var foundChapterNumber bool
		var verseNumberGroup []rune
		var verseNumber int
		var foundVerseNumber bool
		var verseRangeGroup []rune
		var verseRange int
		var foundVerseRange bool
		var lastNumberIndex int

		for i, character := range request[lastLetterIndex:] {
			if unicode.IsNumber(character) {
				if !foundChapterNumber {
					foundChapterNumber = true
					chapterNumberGroup = append(chapterNumberGroup, character)
					lastNumberIndex = i
				} else if !foundVerseNumber && lastNumberIndex == i-1 {
					chapterNumberGroup = append(chapterNumberGroup, character)
					lastNumberIndex = i
				} else if !foundVerseNumber && lastNumberIndex < i-1 {
					foundVerseNumber = true
					verseNumberGroup = append(verseNumberGroup, character)
					lastNumberIndex = i
				} else if !foundVerseRange && lastNumberIndex == i-1 {
					verseNumberGroup = append(verseNumberGroup, character)
					lastNumberIndex = i
				} else if !foundVerseRange && lastNumberIndex < i-1 {
					foundVerseRange = true
					verseRangeGroup = append(verseRangeGroup, character)
					lastNumberIndex = i
				} else if lastNumberIndex == i-1 {
					verseRangeGroup = append(verseRangeGroup, character)
					lastNumberIndex = i
				}
			}
		}

		//convert rune slice groups to ints
		chapterNumber, _ = strconv.Atoi(string(chapterNumberGroup))
		verseNumber, _ = strconv.Atoi(string(verseNumberGroup))
		verseRange, _ = strconv.Atoi(string(verseRangeGroup))

		if !foundChapterNumber {
			chapterNumber = 1
		}
		if !foundVerseNumber {
			verseNumber = 1
		}
		result.Chapter = chapterNumber
		result.VerseNumber = verseNumber
		result.VerseRange = verseRange
	}

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
	return result, nil
}

// return a verse object based on a reference
func getVerse(ref Reference) (Verse, error) {

	var verse *Verse
	for _, bible := range bibleData {
		if bible.Version == "kjv" {
			for _, book := range bible.Books {
				if book.Name == ref.Book {
					for _, chapter := range book.Chapters {
						if chapter.Number == ref.Chapter {
							for _, v := range chapter.Verses {
								if v.Reference.VerseNumber == ref.VerseNumber {
									verse = &v
									break
								}
							}
						}
					}
				}
			}
		}
	}

	return *verse, nil
}
