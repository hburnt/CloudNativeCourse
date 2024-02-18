// Find the top K most common words in a text document.
// Input path: location of the document, K top words
// Output: Slice of top K words
// For this excercise, word is defined as characters separated by a whitespace

// Note: You should use `checkError` to handle potential errors.

package textproc

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

func topWords(path string, K int) []WordCount {
	//1.Read the text file (path specifies the txt file name)
	//2.Create a map to store each unique word & occurrence
	filePath := path
	wordMap := make(map[string]int)

	file, err := os.Open(filePath)

	checkError(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var words []string
	//Stores each word in the text file in a string slice
	for scanner.Scan() {
		line := scanner.Text()
		lineWords := strings.Fields(line)
		words = append(words, lineWords...) //"..." means we are allowing a variable number of arguments
	}

	for _, value := range words {
		_, prs := wordMap[value]

		if prs {
			wordMap[value] = wordMap[value] + 1
			continue
		}
		wordMap[value] = 1
	}

	wordCountSlice := []WordCount{}
	var temp WordCount
	for key, value := range wordMap {
		temp.Word = key
		temp.Count = value
		wordCountSlice = append(wordCountSlice, temp)
	}
	//println(wordCountSlice)
	sortWordCounts(wordCountSlice)
	topKWords := []WordCount{}
	for i := 0; i < K; i++ {
		topKWords = append(topKWords, wordCountSlice[i])
	}

	return topKWords
}

//--------------- DO NOT MODIFY----------------!

// A struct that represents how many times a word is observed in a document
type WordCount struct {
	Word  string
	Count int
}

// Method to convert struct to string format
func (wc WordCount) String() string {
	return fmt.Sprintf("%v: %v", wc.Word, wc.Count)
}

// Helper function to sort a list of word counts in place.
// This sorts by the count in decreasing order, breaking ties using the word.

func sortWordCounts(wordCounts []WordCount) {
	sort.Slice(wordCounts, func(i, j int) bool {
		wc1 := wordCounts[i]
		wc2 := wordCounts[j]
		if wc1.Count == wc2.Count {
			return wc1.Word < wc2.Word
		}
		return wc1.Count > wc2.Count
	})
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
