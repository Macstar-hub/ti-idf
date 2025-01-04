package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

const dockPath string = "./test/sampleFileTrim.txt"
const stopWordFilePath string = "./configs/stopWords.txt"

func main() {
	fmt.Println("Hello from TF-IDF project.")
	DocReader(dockPath)
	DocCleanUp(stopWordFilePath)
}

func DocReader(fileName string) string {
	file, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Cannot read file name: %v with error: %v \n", fileName, err)
	}

	fmt.Println("File size: ", len(file))
	return string(file)
}

func DocCleanUp(stopWordFilePath string) string {
	var cleanContent string
	stopWordFile2, _ := os.Open(stopWordFilePath)
	uncleanContent := DocReader(dockPath)
	uncleanContentLowerCase := strings.ToLower(uncleanContent)

	Lines := bufio.NewScanner(stopWordFile2)
	for Lines.Scan() {
		fmt.Println(Lines.Text())
		stopWordRegex := regexp.MustCompile(Lines.Text())
		fmt.Println(stopWordRegex)
		cleanContent = stopWordRegex.ReplaceAllString(uncleanContentLowerCase, " ")
		uncleanContentLowerCase = cleanContent
	}
	if err := Lines.Err(); err != nil {
		log.Fatal(err)
	}

	cleanData := regexp.MustCompile(`\.\s*`).ReplaceAllString(uncleanContentLowerCase, "\n")
	fmt.Println(uncleanContentLowerCase)
	TFcount(string(uncleanContentLowerCase))
	IDFcount(cleanData)
	fmt.Println("all sentences: \n", cleanData)
	return cleanData
}

func TFcount(fileText string) {
	words := strings.Fields(fileText)
	fmt.Println("All words count: ", len(words))
	for _, word := range words {
		count := 0
		var ratio float64
		for _, targetWord := range words {
			if word == targetWord {
				count = count + 1
			}
		}
		ratio = float64(count) / float64(len(words))
		fmt.Printf("Word '%v': repeat with %v, and with ratio %v \n", word, count, ratio*100)
	}
}

func IDFcount(cleanContent string) {
	var documentNumber int
	var uncleanBagOfWord []string
	documentLines := strings.Split(string(cleanContent), "\n")
	wordRepeatedInDocumnt := make(map[string]int)
	// counter := 0

	// To Make unclean bag of words.
	wordCleanContent := strings.Fields(cleanContent)
	for _, word := range wordCleanContent {
		uncleanBagOfWord = append(uncleanBagOfWord, word)
	}

	// To remove dulplication words in bag of word:
	slices.Sort(uncleanBagOfWord)
	cleanBagOfWord := slices.Compact(uncleanBagOfWord)
	// fmt.Println("=========", cleanBagOfWord)

	// Find total document numbers:
	for index, _ := range documentLines {
		documentNumber = index
	}

	// Iterate over document:
	for wordCounter, documentLine := range documentLines {
		documentLineNew := strings.Fields(documentLine)
		var wordSlices []string

		// Make document to slice format.
		for _, word := range documentLineNew {
			// fmt.Println("Debug each word in document line: ", word)
			wordSlices = append(wordSlices, word)
		}

		// counter = 0
		fmt.Println("Documnt: +++++++++++++++")
		var wordInDocument []string
		// Find bag of word per document:
		for index := range len(wordSlices) {
			for _, word := range cleanBagOfWord {
				if word == wordSlices[index] {
					// counter = 1
					wordInDocument = append(wordInDocument, word)
					fmt.Println("Word: ", word, "Find in document"+strconv.Itoa(wordCounter))
					// wordRepeatedInDocumnt[wordSlices[index]] = counter
				}
			}
		}
		// Make clean word of documents :
		slices.Sort(wordInDocument)
		cleanWordOfDocumnet := slices.Compact(wordInDocument)
		fmt.Println("-------------", cleanWordOfDocumnet)
	}

	// fmt.Println("-------------", wordRepeatedInDocumnt)

	// Make IDF calculation with map elements:
	var ratio float64
	for word, count := range wordRepeatedInDocumnt {
		ratio = 0
		fmt.Println("==========", word, count)
		ratio = float64(documentNumber) / float64(count)

		fmt.Println("=========", math.Log10(ratio))
	}
	// This section for debug purpose.
	// fmt.Println("***************", wordRepeatedInDocumnt)
	fmt.Println("Debug from IDF count: ", documentNumber)
}
