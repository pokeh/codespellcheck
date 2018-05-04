package main

import (
	"testing"
)

func TestIsFile(t *testing.T) {
	var tests = []struct {
		subject  string
		path     string
		expected bool
	}{
		{"file exists", "test/test_file.go", true},
		{"file doesn't exist", "test/test_file_2.go", false},
		{"path is a directory", "test", false},
	}

	for _, tt := range tests {
		t.Run(tt.subject, func(t *testing.T) {
			actual := isFile(tt.path)
			if actual != tt.expected {
				t.Errorf("Expected %v but got %v.", tt.expected, actual)
			}
		})
	}
}

func TestSplitWords(t *testing.T) {
	var tests = []struct {
		subject       string
		src           string
		expected      []string
		errorExpected bool
	}{
		{"snake case", "test_string", []string{"test", "string"}, false},
		{"constants", "TEST_STRING", []string{"test", "string"}, false},
		{"pascal case", "TestString", []string{"test", "string"}, false},
		{"camel case", "testString", []string{"test", "string"}, false},
		{"error with numbers", "testString1", nil, true},
		{"error with symbols", "testString%#$~", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.subject, func(t *testing.T) {
			actual, err := splitWords(tt.src)
			if tt.errorExpected && err == nil {
				t.Errorf("Expected error but did not receive one.")
			}
			if !tt.errorExpected && err != nil {
				t.Errorf("Expected no errors but received: %v", err.Error())
			}
			for i, e := range tt.expected {
				if actual[i] != e {
					t.Errorf("Expected %v but got %v.", e, actual[i])
				}
			}
		})
	}
}

func TestSplitByUnderscore(t *testing.T) {
	var tests = []struct {
		subject  string
		src      string
		expected []string
	}{
		{"multiple words", "test_string", []string{"test", "string"}},
		{"capitals", "TEST_STRING", []string{"test", "string"}},
		{"empty string", "", []string{""}},
	}

	for _, tt := range tests {
		t.Run(tt.subject, func(t *testing.T) {
			actual := splitByUnderscore(tt.src)
			for i, e := range tt.expected {
				if actual[i] != e {
					t.Errorf("Expected %v but got %v.", e, actual[i])
				}
			}
		})
	}
}

func TestSplitByCapitals(t *testing.T) {
	var tests = []struct {
		subject       string
		src           string
		expected      []string
		errorExpected bool
	}{
		{"camel case", "testString", []string{"test", "string"}, false},
		{"pascal case", "TestString", []string{"test", "string"}, false},
		{"single word", "test", []string{"test"}, false},
		{"capitals", "TEST", []string{"t", "e", "s", "t"}, false},
		{"empty string", "", []string{""}, false},
		{"error with symbols", "testString%#$~", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.subject, func(t *testing.T) {
			actual, err := splitByCapitals(tt.src)
			if tt.errorExpected && err == nil {
				t.Errorf("Expected error but did not receive one.")
			}
			if !tt.errorExpected && err != nil {
				t.Errorf("Expected no errors but received: %v", err.Error())
			}
			for i, e := range tt.expected {
				if actual[i] != e {
					t.Errorf("Expected %v but got %v.", e, actual[i])
				}
			}
		})
	}
}

func TestRemoveNonAlphabets(t *testing.T) {
	var tests = []struct {
		subject  string
		src      string
		expected string
	}{
		{"only alphabets", "testString", "testString"},
		{"with underscore", "test_string", "test_string"},
		{"with numbers", "testString01", "testString"},
		{"with symbols", "testString!", "testString"},
		{"empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.subject, func(t *testing.T) {
			actual, err := removeNonAlphabets(tt.src)
			if err != nil {
				t.Errorf("Expected no errors but received: %v", err.Error())
			}
			if actual != tt.expected {
				t.Errorf("Expected %v but got %v.", tt.expected, actual)
			}
		})
	}
}

func TestIsInDictionary(t *testing.T) {
	var tests = []struct {
		subject  string
		word     string
		expected bool
	}{
		{"singular", "rain", true},
		{"plural with s", "grapefruits", true},
		{"plural with es", "passes", true},
		{"past tense verbs", "twirled", true},
	}

	for _, tt := range tests {
		t.Run(tt.subject, func(t *testing.T) {
			actual := isInDictionary(tt.word)
			if actual != tt.expected {
				t.Errorf("Expected %v but got %v: %v", tt.expected, actual, tt.word)
			}
		})
	}
}

func TestcheckByLook(t *testing.T) {
	var tests = []struct {
		subject  string
		word     string
		expected bool
	}{
		{"existing word", "erinaceous", false},
		{"non-existing word", "hedgehogious", true},
	}

	for _, tt := range tests {
		t.Run(tt.subject, func(t *testing.T) {
			actual := checkByLook(tt.word)
			if actual != tt.expected {
				t.Errorf("Expected %v but got %v: %v", tt.expected, actual, tt.word)
			}
		})
	}
}
