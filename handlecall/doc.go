// Package handlecall is the shared handle-method dispatch table for moonBASIC.
//
// The VM uses it at runtime (via thin wrappers in package vm). The compiler
// semantic analyzer uses the same table so --check / LSP reject unknown or
// wrong-arity recv.METHOD(...) calls when the receiver handle type is known.
package handlecall
