package main

import (
	"fmt"
	"os"

	"rkitamu/gocc/args"
)

func main() {
	// コマンドライン引数を解析
	args, err := args.ParseArgs()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Printf("Input file name: %s\n", args.Input)
}
