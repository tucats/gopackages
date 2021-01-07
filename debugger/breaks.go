package debugger

import (
	"errors"
	"strconv"

	"github.com/tucats/gopackages/tokenizer"
)

type breakPointType int

const (
	BreakDisabled breakPointType = 0
	BreakAlways   breakPointType = iota
	BreakValue
)

type breakPoint struct {
	kind breakPointType
	line int
	hit  int
}

var breakPoints = []breakPoint{}

func Break(t *tokenizer.Tokenizer) error {
	var err error
	var line int
	t.Advance(1)

	for t.Peek(1) != tokenizer.EndOfTokens {
		switch t.Next() {
		case "at":
			line, err = strconv.Atoi(t.Next())
			if err == nil {
				err = breakAtLine(line)
			}
		default:
			err = errors.New(InvalidBreakClauseError)
		}

		if err != nil {
			break
		}
	}
	return err
}

func breakAtLine(line int) error {
	b := breakPoint{line: line, hit: 0, kind: BreakAlways}
	breakPoints = append(breakPoints, b)
	return nil
}
