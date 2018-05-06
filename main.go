package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode"
)

func main() {
	paths, err := getFilePaths(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	for _, path := range paths {
		if check(path) != nil {
			log.Fatal(err)
		}
	}
}

func getFilePaths(path string) ([]string, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if fileInfo.IsDir() {
		paths := make([]string, 0, 10)
		err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				paths = append(paths, p)
			}
			return nil
		})
		return paths, err
	} else {
		return []string{path}, nil
	}
}

func check(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		t := scanner.Text()
		for _, word := range splitWords(t) {
			// skip shorter words to avoid abbreviations
			if len(word) < 5 {
				continue
			}
			if !isInDictionary(word) {
				fmt.Printf("Check \"%v\".\n", word)
			}
		}
	}

	return nil
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
	if checkByLook(word) {
		return true
	}
	if strings.HasSuffix(word, "s") {
		try := strings.TrimRight(word, "s")
		if checkByLook(try) {
			return true
		}
	}
	if strings.HasSuffix(word, "es") {
		try := strings.TrimRight(word, "es")
		if checkByLook(try) {
			return true
		}
	}
	if strings.HasSuffix(word, "ies") {
		try := strings.TrimRight(word, "ies")
		try = try + "y"
		if checkByLook(try) {
			return true
		}
	}
	// check past tense verbs
	if strings.HasSuffix(word, "ed") {
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
	o, err := exec.Command("look", word).Output()
	if err != nil {
		return false
	}
	s := strings.ToLower(strings.Split(string(o), "\n")[0])
	return s == word
}
