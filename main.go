package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"slices"
	"strings"
	mysqlconnector "tf-idf/cmd/mysql"
)

const dockPath string = "./test/sampleFileTrim.txt"

// const dockPath string = "./test/sampleFile.txt"
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
	// TFcount(string(uncleanContentLowerCase))
	TFcount(string(cleanData))
	IDFcount(cleanData)
	return cleanData
}

func TFcount(fileText string) {
	bagOfWordMapTF := make(map[string]float64)
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
		bagOfWordMapTF[word] = ratio
		mysqlconnector.Insert(word, count, ratio, 0)

		// For debug make uncomment below line.
		// fmt.Printf("Word '%v': repeat with %v, and with ratio %v \n", word, count, ratio*100)
	}
	fmt.Println("All TF Output: ", bagOfWordMapTF)
}

func IDFcount(cleanContent string) {
	var documentNumber int
	var uncleanBagOfWord []string
	documentLines := strings.Split(string(cleanContent), "\n")
	bagOfWordMap := make(map[string]int)

	// To Make unclean bag of words.
	wordCleanContent := strings.Fields(cleanContent)
	for _, word := range wordCleanContent {
		uncleanBagOfWord = append(uncleanBagOfWord, word)
	}

	// To remove dulplication words in bag of word:
	slices.Sort(uncleanBagOfWord)
	cleanBagOfWord := slices.Compact(uncleanBagOfWord)

	// make init map of clean bags of word.
	for _, word := range cleanBagOfWord {
		bagOfWordMap[word] = 0
	}

	// Find total document numbers:
	for index, _ := range documentLines {
		documentNumber = index
	}

	// Iterate over document:
	for _, documentLine := range documentLines {
		documentLineNew := strings.Fields(documentLine)
		var wordSlices []string

		// Make document to slice format.
		for _, word := range documentLineNew {
			wordSlices = append(wordSlices, word)
		}

		var wordInDocument []string

		// Find bag of word per document:
		for index := range len(wordSlices) {
			for _, word := range cleanBagOfWord {
				if word == wordSlices[index] {
					wordInDocument = append(wordInDocument, word)
				}
			}
		}

		// Make clean word of documents :
		slices.Sort(wordInDocument)
		cleanWordOfDocumnet := slices.Compact(wordInDocument)

		// Compare all word bag with word current sentence.
		for _, wordInSlice := range cleanWordOfDocumnet {
			for wordInBag, value := range bagOfWordMap {
				if wordInSlice == wordInBag {
					counter := 0
					counter++
					bagOfWordMap[wordInBag] = counter + value
				}
			}
		}
	}

	// Make IDF calculation with map elements:
	var ratio float64
	for word, count := range bagOfWordMap {
		ratio = 0
		ratio = float64(documentNumber) / float64(count)
		mysqlconnector.Update(word, math.Log10(ratio))
		fmt.Println("IDF Word: ", word, "Is :", math.Log10(ratio))
	}
}
