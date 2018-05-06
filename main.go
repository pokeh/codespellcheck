package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"unicode"
)

func main() {
	path := os.Args[1]
	if !isFile(path) {
		log.Fatal("Given path is not a file.")
	}
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		text := scanner.Text()
		for _, word := range splitWords(text) {
			// skip shorter words to avoid abbreviations
			if len(word) < 5 {
				continue
			}
			if !isInDictionary(word) {
				fmt.Printf("Check \"%v\".\n", word)
			}
		}
	}
}

func isFile(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !fileInfo.IsDir()
}

func splitWords(src string) []string {
	alphabets := splitByNonalphabets(src)

	res := make([]string, 0, len(src))
	for _, a := range alphabets {
		for _, word := range splitByUppercase(a) {
			res = append(res, word)
		}
	}
	return res
}

func splitByNonalphabets(src string) []string {
	res := make([]string, 0, len(src))
	word := make([]rune, 0, len(src))
	for _, r := range src {
		if unicode.IsLower(r) || unicode.IsUpper(r) {
			word = append(word, r)
		} else if len(word) > 0 {
			res = append(res, string(word))
			word = make([]rune, 0, len(src))
		}
	}
	if len(word) > 0 {
		res = append(res, string(word))
	}
	return res
}

func splitByUppercase(src string) []string {
	res := make([]string, 0, len(src))
	word := make([]rune, 0, len(src))
	uppers := make([]string, 0, len(src))
	for _, r := range src {
		if unicode.IsUpper(r) {
			if len(uppers) == 0 && len(word) > 0 {
				res = append(res, string(word))
				word = make([]rune, 0, len(src))
			}
			uppers = append(uppers, strings.ToLower(string(r)))
		} else if unicode.IsLower(r) {
			if len(uppers) > 0 {
				for _, u := range uppers[:len(uppers)-1] {
					res = append(res, string(u))
				}
				lastU := []rune(uppers[len(uppers)-1])
				word = append(word, []rune(lastU)[0])
				uppers = make([]string, 0, len(src))
			}
			word = append(word, r)
		}
	}
	if len(uppers) > 0 {
		for _, u := range uppers[:len(uppers)-1] {
			res = append(res, string(u))
		}
		lastU := uppers[len(uppers)-1]
		word = append(word, []rune(lastU)[0])
		uppers = make([]string, 0, len(src))
	}
	if len(word) > 0 {
		res = append(res, string(word))
	}
	return res
}

// note: we assume irregular plurals are in the dictionary
func isInDictionary(word string) bool {
	// check singular
	if checkByLook(word) {
		return true
	}
	// check plurals that end with s
	if word[len(word)-1:] == "s" {
		try := strings.TrimRight(word, "s")
		if checkByLook(try) {
			return true
		}
	}
	// check plurals that end with es
	if word[len(word)-2:] == "es" {
		try := strings.TrimRight(word, "es")
		if checkByLook(try) {
			return true
		}
	}
	// check past tense verbs
	if word[len(word)-2:] == "ed" {
		try := strings.TrimRight(word, "ed")
		if checkByLook(try) {
			return true
		}
		try = strings.TrimRight(word, "d")
		if checkByLook(try) {
			return true
		}
	}
	return false
}

func checkByLook(word string) bool {
	_, err := exec.Command("look", word).Output()
	return err == nil
}
