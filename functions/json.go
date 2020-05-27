package functions

import (
	"encoding/json"
	"strings"

	"github.com/tucats/gopackages/symbols"
	"github.com/tucats/gopackages/util"
)

// FunctionDecode reads a string as JSON data
func FunctionDecode(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {

	var v interface{}

	jsonBuffer := util.GetString(args[0])
	err := json.Unmarshal([]byte(jsonBuffer), &v)

	return v, err
}

// FunctionEncode writes a  JSON string from arbitrary data
func FunctionEncode(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {

	if len(args) == 1 {
		jsonBuffer, err := json.Marshal(args[0])
		return string(jsonBuffer), err
	}

	var b strings.Builder
	b.WriteString("[")

	for n, v := range args {
		if n > 0 {
			b.WriteString(", ")
		}
		jsonBuffer, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		b.WriteString(string(jsonBuffer))
	}
	b.WriteString("]")
	return b.String(), nil
}

// FunctionEncodeFormatted writes a  JSON string from arbitrary data
func FunctionEncodeFormatted(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {

	if len(args) == 1 {
		jsonBuffer, err := json.MarshalIndent(args[0], "", "  ")
		return string(jsonBuffer), err
	}

	var b strings.Builder
	b.WriteString("[")

	for n, v := range args {
		if n > 0 {
			b.WriteString(", ")
		}
		jsonBuffer, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			return "", err
		}
		b.WriteString(string(jsonBuffer))
	}
	b.WriteString("]")
	return b.String(), nil
}
