package compiler

import (
	"strconv"
	"sync"
)

var index = 0
var indexMux sync.Mutex

// MakeSymbol creates a unique symbol name for use
// as a temporary variable, etc. during compilation.
func MakeSymbol() string {

	var i int

	indexMux.Lock()
	index = index + 1
	i = index
	indexMux.Unlock()

	return "COMPILER_TEMP__" + strconv.Itoa(i)
}
