# Expressions

This part of the package handles interpreting arbitrary expressions written as text in a string value. The
expression is evaluated using the common rules for precedence and syntax of most common languages.

## Overview

The expression handler supports values of type int,
float64, bool, and string. The evaluator provides automatic type conversion where ever possible. The expression evaluator 
can also access values stored in symbols, which are passed in as a map[string]interface{} object to the evaluator.
Additionally, built-in and caller-supplied functions can be declared in the symbol table as well. Functions accept arbitrary
numbers of arguments of any type, and then operate on them, performing type coercions as needed.

Here is a simple example of using the expression handler:

    symbols := map[string]interface{}{ "name": "Tom", "age": 35}
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

