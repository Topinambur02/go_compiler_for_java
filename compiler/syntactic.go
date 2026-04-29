package compiler

import (
	"github.com/topinambur02/compiler/model"
	"github.com/topinambur02/compiler/parser"
)

func (C *Compiler) SyntacticAnalyzer(tokens []model.Token) *model.Program {
	parser := &parser.Parser{Tokens: tokens}
	ast := parser.Parse()

	if len(parser.Errors) > 0 {
		parser.PrintErrors()
		return nil
	} else {
		return &ast
	}
}
