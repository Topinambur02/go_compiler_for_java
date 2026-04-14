package main

import (
	"fmt"

	"github.com/topinambur02/compiler/compiler"
	"github.com/topinambur02/compiler/utils"
)

func main() {
	filename := "test.java"
	comp := compiler.Compiler{}
	cleanCode, err := comp.Preprocessing(filename)

	if err != nil {
		panic(err)
	}

	tokens, err := comp.LexicalAnalysis(cleanCode)

	if err != nil {
		panic(err)
	}

	result := utils.FormatOutput(tokens)
	fmt.Println(result)
}
