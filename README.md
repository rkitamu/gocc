# gocc

**gocc** is a simple C compiler written in Go.
It supports a small subset of the C language and is intended for educational purposes.

## Getting Started

```bash
go run . -i examples/input.c
```

* CLI interface with flags:
  * `-i` input file
  * `-o` output file (not yet used)
  * `-d` debug mode

## Project Structure

```
.
├── main.go         # CLI entry point
├── lexer/          # Tokenizer for input source
├── parser/         # Parser that builds AST from tokens
├── generator/      # (WIP) Code generation backend
└── examples/       # Sample C source files
```

## TODO

* Code generation to x86-64 assembly
* AST pretty printer
* Type checking
* Function definition and call support
