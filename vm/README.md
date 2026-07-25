# vm/

Bytecode encoding and the moonBASIC virtual machine.

| Package | Role |
|---------|------|
| `opcode/` | Instruction encoding and program layout |
| `moon/` | Interpreter / execution engine |
| Related | Shared handle dispatch lives in [`handlecall/`](../handlecall/) (used by the VM) |

Compile with the root module (`go test ./vm/...`). Graphical builtins need `-tags fullruntime` at the runtime layer, not here.
