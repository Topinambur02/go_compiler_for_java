package main

import (
	"github.com/topinambur02/compiler/compiler"
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
	comp.SemanticAnalysis(ast)
}
