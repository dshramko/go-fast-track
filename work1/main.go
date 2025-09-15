package main

import (
	"bufio"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"
)

type WordCounter struct {
	text  string
	count int
}

type ReadFileError struct {
	msg string
	err error
}

func (e ReadFileError) Error() string {
	return fmt.Sprintf("%s: %v", e.msg, e.err)
}

const minWordLength = 4

func main() {
	filename := flag.String("file", "", "file to read")
	flag.Parse()

	if *filename == "" {
		showUsage()
		return
	}

	words, err := getWordsFromFile(*filename)
	if err != nil {
		slog.Error("Error reading file", "error", err)
		return
	}

	topWords := getTopWords(words, 10)

	printTopWords(topWords)
}

func showUsage() {
	fmt.Println("Counts the number of occurrences of each word in a file.")
	fmt.Println("Usage: go run main.go -file=words.txt")
}

func printTopWords(topWords []WordCounter) {
	fmt.Println("Top words are:")
	for _, wc := range topWords {
		fmt.Printf("%s: %d\n", wc.text, wc.count)
	}
}

func getWordsFromFile(filename string) (map[string]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, ReadFileError{msg: "could not open file", err: err}
	}
	defer file.Close()

	wordCount := make(map[string]int)
	scanner := bufio.NewScanner(file)

	wordRegex := regexp.MustCompile(`[a-zA-Zа-яА-ЯёЁ0-9]+`)

	for scanner.Scan() {
		line := scanner.Text()

		words := wordRegex.FindAllString(line, -1)
		for _, w := range words {
			if utf8.RuneCountInString(w) < minWordLength {
				continue
			}
			wordCount[strings.ToLower(w)]++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, ReadFileError{msg: "error reading file", err: err}
	}

	return wordCount, nil
}

func getTopWords(wordCount map[string]int, topN int) []WordCounter {
	words := make([]WordCounter, 0, len(wordCount))

	for w, c := range wordCount {
		words = append(words, WordCounter{text: w, count: c})
	}

	sort.Slice(words, func(i, j int) bool {
		return words[i].count > words[j].count
	})

	if topN > len(words) {
		topN = len(words)
	}
	return words[:topN]
}
