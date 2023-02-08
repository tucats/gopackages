// Compiler is a package that compiles expression testinto pseudo-code.
// This pseudo-code can then be executed using the bytecode package.
//
// The compiler functions by reading a string value that is entirely contained in
// memmory (there is no Reader interface). It generates a bytecode stream that is
// also stored in memory.
//
// The compiler is a top-down, recursive-descent compiler that works on a stream
// of tokens. Each token contains it's spelling and class (identifier, reserved,
// integer, string, etc). In this way, the tokenizer owns a part of the parsing
// of the code, to establish token meaning. The tokenizer is also reponsible for
// creating composite tokens. For example "<" followed by "=" is converted to a
// single token "<=" by the tokenizer. Thus, the compiler can assume semantically
// correct individual tokens.
//
// There is no linkage phase to compilation. All identifiers are preserved by
// text name in the generated code, and are resolved at runtime. This is used to
// support the flexibility of untyped or relaxed type operations, or to enforce
// strict typing as requested when Ego is run.
//
// The compiler fails on the first error found. There is no "find the next valid
// source and keep trying". Compiler errors are reported using standard Ego Error
// objects, which include the module name and line number (retrieved from the
// tokenizer data) where the error occurred, as well as any context information.
package compiler
