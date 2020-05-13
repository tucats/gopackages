package tokenizer

import (
	"strings"
	"text/scanner"
)

// Tokenizer is an instance of a tokenized string.
type Tokenizer struct {
	Source string
	Tokens []string
	TokenP int
}

// EndOfTokens is a reserved token that means end of the buffer was reached.
const EndOfTokens = "$$$$$EOF$$$$!"

// New creates a tokenizer instance and breaks the string
// up into an array of tokens
func New(src string) *Tokenizer {

	t := Tokenizer{Source: src, TokenP: 0}

	t.Tokens = make([]string, 0)

	var s scanner.Scanner
	s.Init(strings.NewReader(src))
	s.Filename = "string input"
	previousToken := ""

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {

		// See if this is one of the special cases where we patch up the previous token

		nextToken := s.TokenText()

		if nextToken == "=" {
			if InList(previousToken, []string{"!", "<", ">", ":"}) {
				t.Tokens[len(t.Tokens)-1] = previousToken + nextToken
				previousToken = ""
				continue
			}
		}

		previousToken = nextToken
		t.Tokens = append(t.Tokens, nextToken)

	}

	return &t
}

// Next gets the next token in the tokenizer
func (t *Tokenizer) Next() string {
	if t.TokenP >= len(t.Tokens) {
		return EndOfTokens
	}
	token := t.Tokens[t.TokenP]
	t.TokenP = t.TokenP + 1
	return token
}

// Peek looks ahead at the next token without advancing the pointer.
func (t *Tokenizer) Peek(offset int) string {
	if t.TokenP+(offset-1) >= len(t.Tokens) {
		return EndOfTokens
	}
	return t.Tokens[t.TokenP+(offset-1)]
}

// AtEnd indicates if we are at the end of the string
func (t *Tokenizer) AtEnd() bool {
	return t.TokenP >= len(t.Tokens)
}

// Advance moves the pointer
func (t *Tokenizer) Advance(p int) {
	t.TokenP = t.TokenP + p
	if t.TokenP < 0 {
		t.TokenP = 0
	}
	if t.TokenP >= len(t.Tokens) {
		t.TokenP = len(t.Tokens)
	}
}

// InList is an internal function to determine if a string is in a set of
// other possible strings.
func InList(s string, list []string) bool {
	for _, i := range list {
		if s == i {
			return true
		}
	}
	return false
}

// IsNext tests to see if the next token is the given token, and if so
// advances and returns true, else does not advance and returns false.
func (t *Tokenizer) IsNext(test string) bool {
	if t.Peek(1) == test {
		t.Advance(1)
		return true
	}
	return false
}

// IsAnyNext tests to see if the next token is in the given  list
// of tokens, and if so  advances and returns true, else does not
// advance and returns false.
func (t *Tokenizer) IsAnyNext(test []string) bool {

	n := t.Peek(1)
	for _, v := range test {
		if n == v {
			t.Advance(1)
			return true
		}
	}
	return false
}
