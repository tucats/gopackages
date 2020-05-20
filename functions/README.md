# functions

The `functions` package contains the builtin functions (some of which are global and
some of which are arranged into packages). The following provides a short summary of 
each of the built-in functions, broken into categories based on the general data types 
or functionality used.

## Type Casting Functions

These functions are used to explicity specify the type of a value to be used in the
evaluation of the expression. They take an arbitrary value and return that value
coerced to a function-specific type.

### int(any)
Return the argument coerced to an int data type. For a boolean, this
will result in 0 or 1. For a float, it returns the integer component.
A string must contain a valid representation of an integer to convert
without error.

    int(33.5)

This returns the value 33.

### bool(any)
Return the argument coerced to a bool data type. For numeric values,
this means zero for false, or non-zero for true. For a string, it must
contain the strings "true" or "false" to be converted without error.

    bool("true")

This returns the value true.

### float(any)
Return the argument coerced to an float64 data type. For a boolean, this
will result in 0.0 or 1.0 values. For an, it returns the floating point
equivalent of the integer value.
A string must contain a valid representation of an floating point value to convert
without error.
    
    float("3.1415")

Thsi returns the float64 value 3.1415.

## String Functions

These functions act on string values, and usually return a string values as the
result.


### _strings.left(string, count)
This returns `count` characters from the left side of the string.
    
    _strings.left("abraham", 3)

This returns the value "abr".


### _strings.right(string, count)
This returns `count` characters from the right side of the string.
    
    _strings.right("abraham", 4)

This returns the value "aham".

### _strings.substring(string, start, count)
This extracts a substring from the string argument. The substring
starts at the `start` character position, and includes `count` characters
in the result.
    
    _strings.substring("Thomas Jefferson", 8, 4)

This returns the string "Jeff".

### _strings.index(string, substring)
This searches the `string` parameter for the first instance of the
`substring` parameter. If it is found, the function returns the
character position where it starts. If it was not found, it returns
an integer zero.
    
    _strings.index("Scores of fun", "ore")

This returns the value 3, indicating that the string "ore" starts
at the third character of the string.

### _strings.lower(string)
This converts the string value given to lower-case letters. If the
value is not a string, it is first coerced to a string type.
    
    _strings.lower("Tom")

This results in the string value "tom".

### _strings.upper(string)
This converts the string value given to uooer-case letters. If the
value is not a string, it is first coerced to a string type.
    
    _strings.upper("Jeffrey")

This results in the string value "JEFFREY".

## General Functions

These functions work generally with any type of value, and perform coercsions
as needed. The first value in the argument list determines the type that all
the remaining items will be coerced to.


### len(string)
Return the length of the argument. The meaning of _length_ depends on the 
type of the argument. For a string, this returns the number of characters
in the string. For an int, float, or bool value, it returns the number of
characters when the value is formatted for output.

Some examples:

| Example | Result |
|:-|:-|
| len("fortitude")   | 9, the number of characters in the string. |
| len(135)           | 3, the number of characters when 135 is converted to string "135" |
| len(false)         | 5, the number of characters in "false" |
| len(3.1415)        | 6, the number of characters in "3.1415" |
| len([5,3,1])       | 3, the number of elements in the array | 
| len({a:1, b:true}) | 2, the number of fields in the array |

### min(v1, v2...)
This gets the minimum (smallest numeric or alphabetic) value from the list.
If the first item is a string, then all values are converted to a string for
comparison and the result will be the lexigraphically first element. IF the
values are int or float values, then a numeric comparison is done and the
result is the numerically smallest value.
    
    min(33.5, 22.76, 9, 55)
    
This returns the float value 9.0


### max(v1, v2...)
This gets the maximum (largest numeric or alphabetic) value from the list.
If the first item is a string, then all values are converted to a string for
comparison and the result will be the lexigraphically lsat element. IF the
values are int or float values, then a numeric comparison is done and the
result is the numerically largest value.
    
    max("shoe", "mouse", "cake", "whistle")
    
This returns the string value "whistle".

### sum(v1, v2...)
This function returns the sum of the arguments. The meaning of _sum_ depends
on the arguments. The values must not be arrays or structures.

For a numeric value (int or float), the function returns the mathematical
sum of all the numeric values.

    x := sum(3.5, 15, .5)

This results in `x` having the value 19.  For a boolean value, this is the
same as a boolean "and" operation being performed on all values.

For a string, it concatenates all the string values together into a single
long string.


## Utility Functions

These are miscellaneous funcctions to support writing programs in _Solve_.

### _util.sort(array)
This sorts an array into ascending order. The type of the first element in the
array determines the type used to sort all the data; the second and following
array elements are cast to the same type as the first element for the purposes
of sorting the data.

It is an error to call this function with an array that contains elements that
are arrays or structures. It is also an error to call this function with a data
type other than an array.

### _util.uuid()
This generates a UUID (universal unique identifier) and returns it formatted
as a string value. Every call to this function will result in a new unique
value.

### _util.members(st)

Returns an array of strings containing the names of each member of the 
structure passed as an argument. If the value passed is not a structure
it causes an error. Note that the resulting array elements can be used
to reference fields in a structure using array index notation.

    e := { name: "Dave", age: 33 }
    m := _utils.members(e)

    e[m[1]] := 55

The `_util.members()` function returns an array [ "age", "name" ]. These are
the fields of the structure, and they are always returned in alphabetical
order. The assignment statement uses the first array element ("age") to access
the value of e.age.

### _util.symbols()
Returns a string containing a formatted expression of the symbol table at
the moment the function is called, including all nested levels of scope.
The typical use is to simply print the string:

    x := 55
    {
        x = 42
        y := "test"
        print _util.symbols()
    }

This will print the symbols for the nested basic block as well as the
symbols for the main program.