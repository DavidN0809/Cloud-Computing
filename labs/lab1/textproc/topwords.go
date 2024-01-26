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

// topWords takes a file path and an integer K, then returns the top K most frequent words in the file.
func topWords(path string, K int) []WordCount {
	// Open the file located at the given path.
	// If there is an error opening the file, it is handled by the checkError function.
	file, err := os.Open(path)
	checkError(err)
	// Defer the closing of the file until the end of the function execution.
	defer file.Close()

	// Create a map to store word counts. The keys are words, and the values are their counts.
	wordCounts := make(map[string]int)

	// Use a scanner to read the file.
	// The scanner is set to split the input into words (using whitespace as the delimiter).
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	// Iterate over all words in the file.
	for scanner.Scan() {
		// Convert each word to lowercase to ensure case-insensitive counting.
		word := strings.ToLower(scanner.Text())
		// Increment the count of the word in the map.
		wordCounts[word]++
	}

	// Check for any errors that might have occurred during the scanning process.
	checkError(scanner.Err())

	// Create a slice to store word counts in a format that can be sorted.
	wordCountSlice := make([]WordCount, 0, len(wordCounts))
	// Populate the slice with word counts from the map.
	for word, count := range wordCounts {
		wordCountSlice = append(wordCountSlice, WordCount{Word: word, Count: count})
	}

	// Sort the word counts in descending order of count and in alphabetical order for ties.
	sortWordCounts(wordCountSlice)

	// If there are fewer words than K, adjust K to the number of unique words.
	if len(wordCountSlice) < K {
		K = len(wordCountSlice)
	}

	// Return the top K words.
	return wordCountSlice[:K]
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
