package functions

import (
	"encoding/hex"

	"github.com/tucats/gopackages/symbols"
	"github.com/tucats/gopackages/util"
)

// Hash implements the _cipher.hash() function
func Hash(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	return util.Hash(util.GetString(args[0])), nil
}

// Encrypt implements the _cipher.hash() function
func Encrypt(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	b, err := util.Encrypt(util.GetString(args[0]), util.GetString(args[1]))
	if err != nil {
		return b, err
	}
	return hex.EncodeToString([]byte(b)), nil

}

// Decrypt implements the _cipher.hash() function
func Decrypt(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {

	b, err := hex.DecodeString(util.GetString(args[0]))
	if err != nil {
		return nil, err
	}
	return util.Decrypt(string(b), util.GetString(args[1]))
}
