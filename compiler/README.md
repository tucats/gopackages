# compiler

The `compiler` package is used to compile text in the _Solve_ language into
bytecode that can be executed using the `bytecode` package. This allows for
compiled scripts to be integrated into the application, and run repeatedly
without incurring the overhead of parsing and semantic analysis each time.

The _Solve_ language is loosely based on _C_ and it's derivative languages.
Some important attributes of _Solve_ programs are:

* There are currently no pointer types, and no dynamic memory allocation.
* All objects are passed by value in function calls.
* Variables are untyped, but can be cast explicitly or will be type converted
automatically when possible.

The program stream executes at the topmost scope. You can define one or more
functions in that topmost scope, or execute commands directly. Each function
runs in its own scope; it can access variables from outer scopes but cannot
set them. Functions defined within another function only exist as long as
that function is running.

## array
The `array` statement is used to allocate an array. An array can also be
created as an array constant and stored in a variable. The array statement
identifies the name of the array and the size, and optionally an initial
value for each member of the array.

    array x[5]
    array y[2] := 10

The first example creates an array of 5 elements, but the elements are
`<nil>` which means they do not have a usable value yet. The array elements
must have a value stored in them before they can be used in an expression.
The second example assigns an initial value to each element of the array,
so the second statement is really identical to `y := [10,10]`.

## if
The `if` statement provides conditional execution. The statement must start
with a expression which can be cast as a boolean value. That value is
tested; if it is true then the following statement (or statement block)
is execued. By convention, even if the conditional code is a single
statement, it is enclosed in a statement block. For example,

    if age > 50 {
        call aarp(name)
    }

This tests the variable age to determine if it is greater than or
equal to the integer value 50, and if so, it calls the function 
named `aarp` with the value of the `name` symbol.

You can optionally include an "else" clause to execute if the
condition is false, as in 

    if flag = "-d" {
        call debug()
    } else {
        call regular()
    }

If the value of `flag` does not equal the string "-d" then the 
code will call the function `regular()` instead of `debug()`.

## print
The `print` statement accepts a list of expressions, and displays them on
the current stdout stream. There is no formatting built into the `print`
statement at this time; each term in the list is printed sequentially,
and a newline is added after all the items are printed.

    print "My name is ", name


Using `print` without any arguments just prints a newline character.

## call
The `call` statement is used to invoke a function that does not return
a value. It is followed by a function name and arguments, and these are
used to call the named function. However, even if the function uses a
`return` statement to return a value, it is ignored by the `call` 
statement. 

    call profile("secret", "cookies")

This calls the `profile()` function. When that function gets two
paramters, it sets the profile value named in the first argument to
the string value of the second argument. The function returns true
because all functions must return a value, but the `call` statement
discards the result.  This is the same as using the statement:

    
    _ := profile("secret", "cookies")

Where the "_" is the name of the ignored value.


## function
The `function` statement defines a function. This must have a name
which is a valid symbol, and an argument list. The argument list is
a list of names which become local variables in the running function
containing the arguments from the caller. This is then followed by
a statement or block defining the code to execute when the function
is used in an expression or in a `call` statement. For example,

    function double(x) {
        return x * 2
    }

This accepts a single value, named `x` when the function is running.
The function returns that value multiplied by 2. The function can
then be used in an expression, such as:

    fun := 2
    moreFun := double(fun)

After this code executes, `moreFun` will contain the value 4.

## return
The `return` statement contains an expression that is identified as
the result of the function value. The generated code adds the value
to the runtime stack, and then exits the function. The caller can
then retrieve the value from the stack to use in an expression or
statement.

    
    return salary/12.0

This statement returns the value of the expression `salary/12.0` as
the result of the function.

If you use the `return` statement with no value, then the function
simply stops without leaving a value on the arithmetic stack. This is
the appropriate behavior for a function that is meant to be invoked
with a `call` statement.