package main

import (
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

	ast := comp.SyntacticAnalyzer(tokens)
	utils.PrintAST(*ast)
	comp.SemanticAnalysis(ast)
}
