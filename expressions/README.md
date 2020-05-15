# Expressions

This part of the package handles compilation and evaluation of arbitrary expressions written as 
text in a string value. The expression is evaluated using the rules for precedence and syntax of
most common languages.

## Overview

The expression handler supports values of type int, float64, bool, and string. The evaluator 
provides automatic type conversion where ever possible. The expression evaluator 
can also access values stored in symbols, which are passed in as a SymbolTable 
object to the evaluator. Additionally, built-in and caller-supplied functions can be declared 
in the symbol table as well. Functions accept arbitrary numbers of arguments of any type, 
and then operate on them, performing type coercions as needed.

Here is a simple example of using the expression handler:

    // Create a symbol table for use during expression
    // evaluation. This is optional, but must be provided
    // if your expression uses variables.
    symbols := bytecode.NewSymbolTable()
    symbols.Set("name", "Tom")
    symbols.Set("age", 54)

    // Compile a string as an expression and then evaluate
    // the resulting expression to get its value
    e := expressions.New("age + 10")
    v, err := e.eval(symbols)
  
The value of the expression is returned as an opaque interface, along with an error object. If the
error object is nil, no errors occurred during expression evaluation. If the err object is not nil,
it contains an error description, and the value returned will be set to nil.

It is up to the caller to handle the type of the return value. A number of functions exist to fetch
specific types from the opaque value, performing type conversions as needed:

* GetInt()
* GetFloat()
* GetBool()
* GetString()

## Symbols
As shown in the simple example above, you can provide a map of symbols available to the
expression evaluator. All symbol names should be stored as lower-case names as the symbol
table is case-insensitive and all symbol references in the expression are converted to
lower case.

The value of the symbol is any value of the supported types. When that symbol is referenced
in an expression, the map is used to locate the value to insert into the expression evaluator 
at that point. 

The symbol table is passed to the Eval() method to allow a single expression to be used with
many possible values. After the symbol table is created or updated, a new evaluation can be
run using the previoiusly parsed expression.

This is used in github.com/tucats/gopackages/app-cli/tables for example, to support a filter
expression in the tables.SetWhere() method. When an expression is set on a table, the table
printing operation will only show rows for which the where-clause results in a boolean "true"
value. The symbol table is refreshed with each row of the table (the symbol names are taken
from the table column names) so the expression can be re-evaluated with each row.

## Functions
The expression evaluator can process function calls as part of it's processing. The fucntion call
is directed back to code that is either built-in to the expression package, or is supplied by the
user.

The following provides a short summary of each of the built-in functions, broken into
categories based on the general data types or functionality used.

### Type Casting Functions

These functions are used to explicity specify the type of a value to be used in the
evaluation of the expression. They take an arbitrary value and return that value
coerced to a function-specific type.

#### int(any)
Return the argument coerced to an int data type. For a boolean, this
will result in 0 or 1. For a float, it returns the integer component.
A string must contain a valid representation of an integer to convert
without error.

    int(33.5)

This returns the value 33.

#### bool(any)
Return the argument coerced to a bool data type. For numeric values,
this means zero for false, or non-zero for true. For a string, it must
contain the strings "true" or "false" to be converted without error.

    bool("true")

This returns the value true.

#### float(any)
Return the argument coerced to an float64 data type. For a boolean, this
will result in 0.0 or 1.0 values. For an, it returns the floating point
equivalent of the integer value.
A string must contain a valid representation of an floating point value to convert
without error.
    
    float("3.1415")

Thsi returns the float64 value 3.1415.

### String Functions

These functions act on string values, and usually return a string values as the
result.

#### len(string)
Return the length of a string argument in characters as an int value.

    len("fortitude")

This returns the value 9.

#### left(string, count)
This returns `count` characters from the left side of the string.
    
    left("abraham", 3)

This returns the value "abr".


#### right(string, count)
This returns `count` characters from the right side of the string.
    
    right("abraham", 4)

This returns the value "aham".

#### substring(string, start, count)
This extracts a substring from the string argument. The substring
starts at the `start` character position, and includes `count` characters
in the result.
    
    substring("Thomas Jefferson", 8, 4)

This returns the string "Jeff".

#### index(string, substring)
This searches the `string` parameter for the first instance of the
`substring` parameter. If it is found, the function returns the
character position where it starts. If it was not found, it returns
an integer zero.
    
    index("Scores of fun", "ore")

This returns the value 3, indicating that the string "ore" starts
at the third character of the string.

#### lower(string)
This converts the string value given to lower-case letters. If the
value is not a string, it is first coerced to a string type.
    
    lower("Tom")

This results in the string value "tom".

#### upper(string)
This converts the string value given to uooer-case letters. If the
value is not a string, it is first coerced to a string type.
    
    upper("Jeffrey")

This results in the string value "JEFFREY".

### General Functions

These functions work generally with any type of value, and perform coercsions
as needed. The first value in the argument list determines the type that all
the remaining items will be coerced to.

#### min(v1, v2...)
This gets the minimum (smallest numeric or alphabetic) value from the list.
If the first item is a string, then all values are converted to a string for
comparison and the result will be the lexigraphically first element. IF the
values are int or float values, then a numeric comparison is done and the
result is the numerically smallest value.
    
    min(33.5, 22.76, 9, 55)
    
This returns the float value 9.0


#### max(v1, v2...)
This gets the maximum (largest numeric or alphabetic) value from the list.
If the first item is a string, then all values are converted to a string for
comparison and the result will be the lexigraphically lsat element. IF the
values are int or float values, then a numeric comparison is done and the
result is the numerically largest value.
    
    min("shoe", "mouse", "cake", "whistle")
    
This returns the string value "whistle".

### User Suppied Functions
The caller of the expressions package can supply additional functions to
supplement the built-in functions.  The function must be declared as a
function of type func([]interface{})(interface{}, error).  For example,
this is a simplified function that creates a floating point sum of all
the supplied values (which will be type-coerced to be floats):
    
    func sum( args []interface{})(interface{}, error) {
        result := 0
        for _, v := range args {
            result = result + util.GetInt(v)
        }
        return result, nil
    }

The body of the function operates on the argument list, processing values
as appropriate for the function. The service functions GetInt, GetFloat,
GetString, and GetBool can be used to get the value of an opaque argument
and coerce it to the desired type. The function can also implement type
switch statements to handle other argument types.

The result is always returned as an int, float64, string, or bool. In
addition, if the function encounters an error (for example, an incorrect
number of arguments passed in) then an error should be returned. When an
error is returned, the function result is ignored and the expression
evaluation reports an error.

To declare the function to the expression evaluator, just add it to the
symbol table. The name must be the name of the function as it would be
specified in an expression, followed by "()". The value of the item
in the symbol table map is the function pointer or value itself.
    
    symbols.Set("sum()", sum)

This will add the sum function described above to the available symbols
for processing an expression. Note that user-supplied functions must be
put in any symbol table where the function is expected to be used; 
user-defined function definitions do not persist past the Eval() call.
