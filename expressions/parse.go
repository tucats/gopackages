package expressions

import (
	"strings"
	"text/scanner"
)

// Parse parses a text expression
func (e *Expression) Parse() error {

	e.Tokens = Tokenize(e.Source)
	e.TokenP = 0

	return nil
}

// Tokenize breaks the string up into an array of tokens
func Tokenize(src string) []string {

	var t = make([]string, 0)

	var s scanner.Scanner
	s.Init(strings.NewReader(src))
	s.Filename = "example"
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		t = append(t, s.TokenText())
	}

	return t
}
