package functions

import (
	"encoding/hex"

	"github.com/tucats/gopackages/symbols"
	"github.com/tucats/gopackages/util"
)

// FunctionHash implements the _cipher.hash() function
func FunctionHash(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	return util.Hash(util.GetString(args[0])), nil
}

// FunctionEncrypt implements the _cipher.hash() function
func FunctionEncrypt(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	b, err := util.Encrypt(util.GetString(args[0]), util.GetString(args[1]))
	if err != nil {
		return b, err
	}
	return hex.EncodeToString([]byte(b)), nil

}

// FunctionDecrypt implements the _cipher.hash() function
func FunctionDecrypt(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {

	b, err := hex.DecodeString(util.GetString(args[0]))
	if err != nil {
		return nil, err
	}
	return util.Decrypt(string(b), util.GetString(args[1]))
}
