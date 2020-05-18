# bytecode

The `bytecode` subpackage supports a simple bytecode intepreter. This allows operations (especially those that might be
repeated) to be compiled into an expression of the semantics of the operation, without having to have the string 
parsed and lexically analyized repeatedly.

Bytecode can be generated explicitly (as in the first example below) or by using the compiler package which accepts
text in a Go-like language called _Solve_ and generates bytecode. Once the bytecode is generated, a runtime `Context`
object is created which is used to manage the execution of a bytecode stream. This includes it's active symbol table,
program counter, stack, etc. A `Context` is separate from the bytecode because the same bytecode could be executed 
on multiple threads, each with it's own `Context`.

The bytecode also supports a symbol table. This can be used to store named values and retrieve them as part of the
execution of the bytecode. The symbol table also contains function pointers for each of the built-in function and
function packages. Calling functions is managed by the bytecode, but can be used to call a function provided
by the caller as native Go code.


## Example
Here is a trival example of generating bytecode and executing it.
    
    b := bytecode.New("sample program")
    b.Emit2(I{bytecode.Load, "left"})
    b.Emit2(I{bytecode.Push, "fruitcake"})
    b.Emit2(I{bytecode.Push, 5})
    b.Emit2(I{bytecode.Call, 2})
    b.Emit1(I{bytecode.Stop})

    // Make a runtime context for this bytecode, and then run it.
    // The context contains the stack symbol table (if any), etc.
    c := bytecode.NewContext(nil, b)
    err := c.Run()

    // Retrieve the last value
    v, err := b.Pop()

    fmt.Printf("The result is %s\n", util.GetString(v))

This creates a new bytecode stream, and then adds instructions to it. These instructions would nominally
be added by a parser. The `Emit1()` function emits an instruction with only one value, the opcode. The
`Emit2()` method emits an instruction with two values, the opcode and an arbitrary operand value.

The stream puts arguments to a function on a stack, and then calls the function. The
result is left on the stack, and can be popped off after execution completes. The result (which is always
an abstract interface{}) is then converted to a string and printed.
