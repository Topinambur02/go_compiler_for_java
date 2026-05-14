package main

import (
	"fmt"

	"github.com/topinambur02/compiler/compiler"
	"github.com/topinambur02/compiler/utils"
)

func main() {
	filename := "test.java"
	comp := compiler.Compiler{}
	fmt.Println("######## 1. Preprocessing ########")
	fmt.Println()
	cleanCode, err := comp.Preprocessing(filename)

	if err != nil {
		panic(err)
	}

	fmt.Println(cleanCode)
	fmt.Println()
	fmt.Println("######## 2. Lexical Analysis ########")
	fmt.Println()

	tokens, err := comp.LexicalAnalysis(cleanCode)

	if err != nil {
		panic(err)
	}

	fmt.Println(tokens)
	fmt.Println()
	fmt.Println("######## 3. Syntactic Analysis ########")
	fmt.Println()

	ast := comp.SyntacticAnalyzer(tokens)
	utils.PrintAST(*ast)

	fmt.Println()
	fmt.Println("######## 4. Semantic Analysis ########")
	fmt.Println()
	
	comp.SemanticAnalysis(ast)
}
