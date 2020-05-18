package tokenizer

import (
	"fmt"
	"strings"
	"text/scanner"
)

// Tokenizer is an instance of a tokenized string.
type Tokenizer struct {
	Source string
	Tokens []string
	TokenP int
	Line   []int
	Pos    []int
}

// EndOfTokens is a reserved token that means end of the buffer was reached.
const EndOfTokens = "$$$$$EOF$$$$!"

// New creates a tokenizer instance and breaks the string
// up into an array of tokens
func New(src string) *Tokenizer {

	// Strip '#' comments from input file
	src = stripComments(src)
	t := Tokenizer{Source: src, TokenP: 0}
	t.Tokens = make([]string, 0)

	var s scanner.Scanner
	s.Init(strings.NewReader(src))
	s.Filename = "Input"

	previousToken := ""

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {

		// See if this is one of the special cases where we patch up the previous token

		nextToken := s.TokenText()

		if nextToken == "[" {
			if previousToken == "[" {
				t.Tokens[len(t.Tokens)-1] = previousToken + nextToken
				previousToken = ""
				continue
			}
		}

		if nextToken == "]" {
			if previousToken == "]" {
				t.Tokens[len(t.Tokens)-1] = previousToken + nextToken
				previousToken = ""
				continue
			}
		}

		if nextToken == "=" {
			if InList(previousToken, []string{"!", "<", ">", ":"}) {
				t.Tokens[len(t.Tokens)-1] = previousToken + nextToken
				previousToken = ""
				continue
			}
		}

		previousToken = nextToken
		t.Tokens = append(t.Tokens, nextToken)
		t.Line = append(t.Line, s.Line)
		t.Pos = append(t.Pos, s.Column)
	}

	return &t
}

// PositionString reports the position of the current
// token in terms of line and column information.
func (t *Tokenizer) PositionString() string {

	p := t.TokenP
	if p >= len(t.Line) {
		p = len(t.Line) - 1
	}
	return fmt.Sprintf("at line %d, column %d,", t.Line[p], t.Pos[p])
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

func stripComments(source string) string {

	var result strings.Builder

	ignore := false
	startOfLine := true
	for _, c := range source {

		// Is this a # on the start of a line? If so, start
		// ignoring characters. If it's the end of line, then
		// reset to end-of-line and resume processing characters.
		// Finally, if nothing else, copy character if not ignoring.
		if c == '#' && startOfLine {
			ignore = true
		} else if c == '\n' {
			ignore = false
			startOfLine = true
			result.WriteRune(c)
		} else if !ignore {
			result.WriteRune(c)
		}
	}

	return result.String()
}

// IsSymbol is a utility function to determine if a token is a symbol name.
func IsSymbol(s string) bool {

	for n, c := range s {
		if isLetter(c) {
			continue
		}

		if isDigit(c) && n > 0 {
			continue
		}
		return false
	}
	return true
}

func isLetter(c rune) bool {
	for _, d := range "_abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		if c == d {
			return true
		}
	}
	return false
}

func isDigit(c rune) bool {

	for _, d := range "0123456789" {
		if c == d {
			return true
		}
	}
	return false
}

// GetLine returns a given line of text from the token stream.
func (t *Tokenizer) GetLine(line int) string {

	var b strings.Builder

	for n, text := range t.Tokens {
		if t.Line[n] == line {
			if b.Len() > 0 && !InList(text, []string{",", ";"}) {
				b.WriteRune(' ')
			}
			b.WriteString(text)
		}
		if t.Line[n] > line {
			break
		}
	}
	return b.String()
}
