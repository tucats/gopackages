package builtins

import (
	"sync"

	"github.com/tucats/gopackages/errors"
	"github.com/tucats/gopackages/expressions/data"
	"github.com/tucats/gopackages/expressions/symbols"
)

// Mutex functions.

// sync.Mutex.Lock() function.
func mutexLock(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	if len(args) != 0 {
		return nil, errors.ErrArgumentCount.In("Lock")
	}

	this := getNativeThis(s)
	if m, ok := this.(*sync.Mutex); ok {
		m.Lock()

		return nil, nil
	}

	return nil, errors.ErrInvalidThis
}

// sync.Mutex.Unlock() function.
func mutexUnlock(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	if len(args) != 0 {
		return nil, errors.ErrArgumentCount.In("Unock")
	}

	this := getNativeThis(s)
	if m, ok := this.(*sync.Mutex); ok {
		m.Unlock()

		return nil, nil
	}

	return nil, errors.ErrInvalidThis
}

// Waitgroup functions.

// sync.WaitGroup Add() function.
func waitGroupAdd(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, errors.ErrArgumentCount.In("Add")
	}

	this := getNativeThis(s)
	if wg, ok := this.(*sync.WaitGroup); ok {
		count := data.Int(args[0])
		wg.Add(count)

		return nil, nil
	}

	return nil, errors.ErrInvalidThis
}

// sync.WaitGroup Done() function.
func waitGroupDone(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	if len(args) != 0 {
		return nil, errors.ErrArgumentCount.In("Done")
	}

	this := getNativeThis(s)
	if wg, ok := this.(*sync.WaitGroup); ok {
		wg.Done()

		return nil, nil
	}

	return nil, errors.ErrInvalidThis
}

// sync.WaitGroup Wait() function.
func waitGroupWait(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	if len(args) != 0 {
		return nil, errors.ErrArgumentCount.In("Wait")
	}

	this := getNativeThis(s)
	if wg, ok := this.(*sync.WaitGroup); ok {
		wg.Wait()

		return nil, nil
	}

	return nil, errors.ErrInvalidThis
}
