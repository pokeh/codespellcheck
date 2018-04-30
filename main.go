package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func main() {
	path := os.Args[1]
	if !isFile(path) {
		log.Fatal("Given path is not a file.")
		return
	}
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		text := scanner.Text()
		for _, word := range splitWords(text) {
			word, err = removeNonAlphabets(word)
			if err != nil {
				log.Fatal(err)
				return
			}
			// skip shorter words to avoid abbreviations
			if len(word) < 5 {
				continue
			}
			_, err := exec.Command("look", word).Output()
			if err != nil {
				fmt.Printf("Check \"%v\".\n", word)
			}
		}
	}
}

func isFile(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func splitWords(src string) []string {
	if strings.Contains(src, "_") {
		return splitByUnderscore(src)
	} else {
		return splitByCapitals(src)
	}
}

// eg. something_in_snake_case -> [something, in, snake, case]
// eg. SOMETHING_IN_CAPITALS -> [something, in, capitals]
func splitByUnderscore(src string) []string {
	return strings.Split(src, "_")
}

// if more than three consecutive upper letters, use letters 0..n-1 as a word
// eg. ThisIsHTMLForYou -> [this, is, html, for, you]
func splitByCapitals(src string) []string {
	// TODO: IMPLEMENT
	return []string{src}
}

func removeNonAlphabets(src string) (string, error) {
	reg, err := regexp.Compile("[^a-zA-Z]+")
	if err != nil {
		return "", err
	}
	return reg.ReplaceAllString(src, ""), nil
}
