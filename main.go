package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"unicode"
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
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		text := scanner.Text()
		text, err = removeNonAlphabets(text)
		if err != nil {
			log.Fatal(err)
		}
		words, err := splitWords(text)
		if err != nil {
			log.Fatal(err)
		}
		for _, word := range words {
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
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !fileInfo.IsDir()
}

func splitWords(src string) ([]string, error) {
	if strings.Contains(src, "_") {
		return splitByUnderscore(src), nil
	} else {
		text, err := splitByCapitals(src)
		if err != nil {
			return nil, err
		}
		return text, nil
	}
}

func splitByUnderscore(src string) []string {
	return strings.Split(strings.ToLower(src), "_")
}

func splitByCapitals(src string) ([]string, error) {
	var res []string
	buf := bytes.NewBuffer(make([]byte, 0, 100))
	for _, rune := range src {
		switch true {
		case unicode.IsLower(rune):
			buf.Write([]byte(string(rune)))
		case unicode.IsUpper(rune):
			if len(buf.String()) > 0 {
				res = append(res, buf.String())
				buf.Reset()
			}
			buf.Write([]byte(string(unicode.ToLower(rune))))
		default:
			return nil, fmt.Errorf("Unexpected letter: %v", string(rune))
		}
	}
	res = append(res, buf.String())
	return res, nil
}

// CAVEAT: also leaves in underscores for later parsing
func removeNonAlphabets(src string) (string, error) {
	reg, err := regexp.Compile("[^a-zA-Z_]+")
	if err != nil {
		return "", err
	}
	return reg.ReplaceAllString(src, ""), nil
}
