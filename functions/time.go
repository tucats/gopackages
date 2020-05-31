package functions

import (
	"time"

	"github.com/tucats/gopackages/symbols"
	"github.com/tucats/gopackages/util"
)

const basicLayout = "Mon Jan 2 15:04:05 MST 2006"

// FunctionTimeNow implements _time.now()
func FunctionTimeNow(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	return time.Now().Format(basicLayout), nil
}

// FunctionTimeAdd implements _time.duration()
func FunctionTimeAdd(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	t, err := time.Parse(basicLayout, util.GetString(args[0]))
	if err != nil {
		return nil, err
	}
	d, err := time.ParseDuration(util.GetString(args[1]))
	if err != nil {
		return nil, err
	}

	t2 := t.Add(d)
	return t2.Format(basicLayout), nil
}

// FunctionTimeSub implements _time.duration()
func FunctionTimeSub(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	t, err := time.Parse(basicLayout, util.GetString(args[0]))
	if err != nil {
		return nil, err
	}
	d, err := time.Parse(basicLayout, util.GetString(args[1]))
	if err != nil {
		return nil, err
	}

	t2 := t.Sub(d)
	return t2.String(), nil
}
