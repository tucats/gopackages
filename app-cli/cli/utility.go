package cli

import (
	"strings"
)

// ValidKeyword does a case-insensitive compare of a string containing
// a keyword against a list of possible stirng values.
func ValidKeyword(test string, valid []string) bool {

	for _, v := range valid {
		if strings.ToLower(test) == strings.ToLower(v) {
			return true
		}
	}
	return false
}

// FindKeyword does a case-insensitive compare of a string containing
// a keyword against a list of possible string values. If the keyword
// is found, it's position in the list is returned. If it was not found,
// the value returned is -1
func FindKeyword(test string, valid []string) int {

	for n, v := range valid {
		if strings.ToLower(test) == strings.ToLower(v) {
			return n
		}
	}
	return -1
}
