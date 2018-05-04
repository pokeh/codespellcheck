package main

import (
	"bufio"
	"bytes"
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
	cap := len([]byte(src))
	alphabets := splitByNonalphabets(cap, src)

	res := make([]string, 0, len(src))
	for _, a := range alphabets {
		for _, word := range splitByCapitals(cap, a) {
			res = append(res, word)
		}
	}
	return res
}

func splitByNonalphabets(cap int, src string) []string {
	res := make([]string, 0, len(src))
	buf := bytes.NewBuffer(make([]byte, 0, cap))
	for _, rune := range src {
		if unicode.IsLower(rune) || unicode.IsUpper(rune) {
			buf.Write([]byte(string(rune)))
		} else if len(buf.String()) > 0 {
			res = append(res, buf.String())
			buf.Reset()
		}
	}
	if len(buf.String()) > 0 {
		res = append(res, buf.String())
		buf.Reset()
	}
	return res
}

// TODO: refactor
func splitByCapitals(cap int, src string) []string {
	res := make([]string, 0, len(src))
	buf := bytes.NewBuffer(make([]byte, 0, cap))
	cs := make([]rune, 0, len(src))
	for _, r := range src {
		if unicode.IsUpper(r) {
			if len(cs) == 0 && len(buf.String()) > 0 {
				res = append(res, buf.String())
				buf.Reset()
			}
			lower := []rune(strings.ToLower(string(r)))[0]
			cs = append(cs, lower)
		} else if unicode.IsLower(r) {
			if len(cs) > 0 {
				for _, c := range cs[:len(cs)-1] {
					res = append(res, string(c))
				}
				lastC := string(cs[len(cs)-1])
				buf.Write([]byte(lastC))
				cs = cs[:0]
			}
			buf.Write([]byte(string(r)))
		}
	}
	if len(cs) > 0 {
		for _, c := range cs[:len(cs)-1] {
			res = append(res, string(c))
		}
		lastC := string(cs[len(cs)-1])
		buf.Write([]byte(lastC))
		cs = cs[:0]
	}
	if len(buf.String()) > 0 {
		res = append(res, buf.String())
		buf.Reset()
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
