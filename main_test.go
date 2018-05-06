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
		subject  string
		src      string
		expected []string
	}{
		{"snake case", "test_string", []string{"test", "string"}},
		{"pascal case", "TestString", []string{"test", "string"}},
		{"camel case", "testString", []string{"test", "string"}},
		{"with numbers", "test_string_01", []string{"test", "string"}},
		{"with symbols", "test(string)", []string{"test", "string"}},
		{"with nonalphabets", "testがStringで", []string{"test", "string"}},
		// highly probable that these words are abbreviations
		{"with capitalized words", "TESTString", []string{"t", "e", "s", "t", "string"}},
		// fix me (though, how do we differentiate the abbreviations?)
		{"only capitalized words", "TEST_STRING", []string{"t", "e", "s", "t", "s", "t", "r", "i", "n", "g"}},
	}

	for _, tt := range tests {
		t.Run(tt.subject, func(t *testing.T) {
			actual := splitWords(tt.src)
			if len(actual) != len(tt.expected) {
				t.Errorf("Expected %v but got %v.", tt.expected, actual)
			}
			for i, e := range tt.expected {
				if actual[i] != e {
					t.Errorf("Expected %v but got %v.", e, actual[i])
				}
			}
		})
	}
}

func TestSplitByNonalphabets(t *testing.T) {
	var tests = []struct {
		subject  string
		src      string
		expected []string
	}{
		{"with underscore", "test_string", []string{"test", "string"}},
		{"with symbols", "test(string)", []string{"test", "string"}},
		{"with nonalphabets", "testがstringよ", []string{"test", "string"}},
		{"empty string", "", nil},
	}

	for _, tt := range tests {
		t.Run(tt.subject, func(t *testing.T) {
			actual := splitByNonalphabets(tt.src)
			if len(actual) != len(tt.expected) {
				t.Errorf("Expected %v but got %v.", tt.expected, actual)
			}
			for i, e := range tt.expected {
				if actual[i] != e {
					t.Errorf("Expected %v but got %v.", e, actual[i])
				}
			}
		})
	}
}

func TestSplitByUppercase(t *testing.T) {
	var tests = []struct {
		subject  string
		src      string
		expected []string
	}{
		{"single word", "test", []string{"test"}},
		{"camel case", "testString", []string{"test", "string"}},
		{"pascal case", "TestString", []string{"test", "string"}},
		{"capitals", "TEST", []string{"t", "e", "s", "t"}},
		{"partly capitals", "TESTString", []string{"t", "e", "s", "t", "string"}},
		{"partly capitals 2", "TESTStringGO", []string{"t", "e", "s", "t", "string", "g", "o"}},
		{"partly capitals 3", "testSTRINGGo", []string{"test", "s", "t", "r", "i", "n", "g", "go"}},
		{"empty string", "", nil},
	}

	for _, tt := range tests {
		t.Run(tt.subject, func(t *testing.T) {
			actual := splitByUppercase(tt.src)
			if len(actual) != len(tt.expected) {
				t.Errorf("Expected %v but got %v.", tt.expected, actual)
			}
			for i, e := range tt.expected {
				if actual[i] != e {
					t.Errorf("Expected %v but got %v.", e, actual[i])
				}
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
