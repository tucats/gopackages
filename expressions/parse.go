package expressions

import (
	"strings"
	"text/scanner"
)

// Parse parses a text expression
func (e *Expression) Parse(s string) error {

	e.Source = s
	e.Tokens = Tokenize(e.Source)
	e.TokenP = 0

	return nil
}

// Tokenize breaks the string up into an array of tokens
func Tokenize(src string) []string {

	var tokens = make([]string, 0)

	var s scanner.Scanner
	s.Init(strings.NewReader(src))
	s.Filename = "example"
	previousToken := ""

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {

		// See if this is one of the special cases where we patch up the previous token

		nextToken := s.TokenText()

		if nextToken == "=" {
			if inList(previousToken, []string{"!", "<", ">", ":"}) {
				tokens[len(tokens)-1] = previousToken + nextToken
				previousToken = ""
				continue
			}
		}

		previousToken = nextToken
		tokens = append(tokens, nextToken)

	}

	return tokens
}

// inList is an internal function to determine if a string is in a set of
// other possible strings.
func inList(s string, list []string) bool {
	for _, i := range list {
		if s == i {
			return true
		}
	}
	return false
}
