# bytecode

The `bytecode` subpackage supports a simple bytecode intepreter. This allows operations (especially those that might be
repeated) to be compiled into an expression of the semantics of the operation, without having to have the string 
parsed and lexically analyized repeatedly.

## Example
Here is a trival example of generating bytecode and executing it.

    
    b := bytecode.New("sample program")
    b.Emit(I{bytecode.Push, "fruitcake"})
    b.Emit(I{bytecode.Push, 5})
    b.Emit(I{bytecode.Push, "left"})
    b.Emit(I{bytecode.Call, 2})
    b.Emit(I{bytecode.Stop, 0})

    err := b.Run()
    v, err := b.Pop()

    fmt.Printf("The result is %s\n", util.GetString(v))

This creates a new bytecode stream, and then adds instructions to it. These instructions would nominally
be added by a parser. The stream puts arguments to a function on a stack, and then calls the function. The
result is left on the stack, and can be popped off after execution completes. The result (which is always
an abstract interface{}) is then converted to a string and printed.