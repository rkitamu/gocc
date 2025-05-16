package main

import (
	"flag"
	"fmt"
	"os"

	"rkitamu/gocc/generator"
	"rkitamu/gocc/lexer"
	"rkitamu/gocc/parser"
)

type Args struct {
	Input  string
	Output string
	Debug  bool
}

func parseArgs() (*Args, error) {
	input := flag.String("i", "", "Input file name")
	output := flag.String("o", "out.s", "Output file name")
	debug := flag.Bool("d", false, "Enable debug mode")

	flag.Parse()

	if *input == "" {
		return nil, fmt.Errorf("input file name is required")
	}

	args := &Args{
		Input:  *input,
		Output: *output,
		Debug:  *debug,
	}

	return args, nil
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	// コマンドライン引数を解析
	cliArgs, err := parseArgs()
	if err != nil {
		return err
	}

	// read input file
	inputByte, err := os.ReadFile(cliArgs.Input)
	input := string(inputByte)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	// lex input
	tokens, err := lexer.Lex(input)
	if err != nil {
		return err
	}

	// parse tokens
	parser := parser.NewParser(tokens, input)
	node, err := parser.Parse()
	if err != nil {
		return err
	}

	// optionally print AST
	if cliArgs.Debug {
		fmt.Println("=== AST ===")
		parser.PrintTree(node)
	}

	// generate assembly code
	gen := generator.NewGenerator()
	asm := gen.Generate(node)

	// write to output file
	if err := os.WriteFile(cliArgs.Output, []byte(asm), 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}
