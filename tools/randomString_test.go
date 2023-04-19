package tools_test

import (
	"sokker-org-auto-bidder/tools"
	"testing"
)

func TestString(t *testing.T) {
	for i := 0; i < 10; i++ {
		if strlen := len(tools.String(i)); strlen != i {
			t.Errorf("Generated string length %d not equal to asked %d", strlen, i)
		}
	}
}

type stringWithCharsetTest struct {
	charset string
	length int
}

var stringWithCharsetTests = []stringWithCharsetTest{
	stringWithCharsetTest{"AB", 10},
	stringWithCharsetTest{"12", 20},
	stringWithCharsetTest{"AB12ab", 30},
}

func sliceIncludes(slice []rune, char rune) bool {
	for _, c := range slice {
		if c == char {
			return true
		}
	}
	return false
}

func getUniqueChars(s string) []rune {
	chars := []rune(s)
	uniqueChars := []rune{}
	for	_, c := range chars {
		if !sliceIncludes(uniqueChars, c) {
			uniqueChars = append(uniqueChars, c)
		}
	}
	return uniqueChars
}

func compareSlices(s1, s2 []rune) bool {
	// TODO it can be done better way
	s1len := len(s1)
	s2len := len(s2)
	if s1len != s2len {
		return false
	}

	same := []rune{}
	for _, c := range s1 {
		if sliceIncludes(s2, c) {
			same = append(same, c)
		}
	}
	if s1len != len(same) {
		return false
	}

	same = []rune{}
	for _, c := range s2 {
		if sliceIncludes(s1, c) {
			same = append(same, c)
		}
	}
	return s2len == len(same)
}

func TestStringWithCharset(t *testing.T) {
	for _, test := range stringWithCharsetTests {
		s := tools.StringWithCharset(test.length, test.charset)
		if strlen := len(s); strlen != test.length {
			t.Errorf("Generated string length %d not equal to asked %d", strlen, test.length)
		}
		
		usedChars := getUniqueChars(s)
		askedChars := getUniqueChars(test.charset)
		if !compareSlices(usedChars, askedChars) {
			t.Errorf("Generated string consists of %v but asked for %v", usedChars, askedChars)
		}
	}
}
