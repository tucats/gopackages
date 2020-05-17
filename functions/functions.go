package functions

// FunctionDefinition is an element in the function dictionary
type FunctionDefinition struct {
	Pkg string
	Min int
	Max int
	F   interface{}
}

// FunctionDictionary is the dictionary of functions
var FunctionDictionary = map[string]FunctionDefinition{
	"int":       FunctionDefinition{Min: 1, Max: 1, F: FunctionInt},
	"bool":      FunctionDefinition{Min: 1, Max: 1, F: FunctionBool},
	"float":     FunctionDefinition{Min: 1, Max: 1, F: FunctionFloat},
	"string":    FunctionDefinition{Min: 1, Max: 1, F: FunctionString},
	"len":       FunctionDefinition{Min: 1, Max: 1, F: FunctionLen},
	"left":      FunctionDefinition{Min: 2, Max: 2, F: FunctionLeft, Pkg: "strings"},
	"right":     FunctionDefinition{Min: 2, Max: 2, F: FunctionRight, Pkg: "strings"},
	"substring": FunctionDefinition{Min: 3, Max: 3, F: FunctionSubstring, Pkg: "strings"},
	"index":     FunctionDefinition{Min: 2, Max: 2, F: FunctionIndex},
	"upper":     FunctionDefinition{Min: 1, Max: 1, F: FunctionUpper, Pkg: "strings"},
	"lower":     FunctionDefinition{Min: 1, Max: 1, F: FunctionLower, Pkg: "strings"},
	"min":       FunctionDefinition{Min: 1, Max: 99999, F: FunctionMin},
	"max":       FunctionDefinition{Min: 1, Max: 99999, F: FunctionMax},
	"sum":       FunctionDefinition{Min: 1, Max: 99999, F: FunctionSum},
	"uuid":      FunctionDefinition{Min: 0, Max: 0, F: FunctionUUID, Pkg: "util"},
	"profile":   FunctionDefinition{Min: 1, Max: 2, F: FunctionProfile, Pkg: "util"},
	"array":     FunctionDefinition{Min: 1, Max: 2, F: FunctionArray},
}
