package utils

import (
	"fmt"
	"strings"

	"github.com/topinambur02/compiler/model"
)

func PrintAST(node any) {
	printNode(node, "", true)
}

func printNode(node any, indent string, isLast bool) {
	marker := "├── "
	if isLast {
		marker = "└── "
	}
	fmt.Print(indent + marker)

	nextIndent := indent
	if isLast {
		nextIndent += "    "
	} else {
		nextIndent += "│   "
	}

	switch n := node.(type) {
	case model.Program:
		fmt.Println("Program Root")

		numImports := len(n.Imports)
        numClasses := len(n.Classes)

		for i, imp := range n.Imports {
            isLast := (i == numImports-1) && (numClasses == 0)
            printNode(imp, nextIndent, isLast)
        }

		for i, cls := range n.Classes {
			printNode(cls, nextIndent, i == len(n.Classes)-1)
		}
	case model.ClassDecl:
		fmt.Printf("Class: %s\n", n.Name)
		for i, m := range n.Methods {
			printNode(m, nextIndent, i == len(n.Methods)-1)
		}
	case model.MethodDecl:
		params := []string{}
		for _, p := range n.Params {
			params = append(params, p.Type+" "+p.Name)
		}
		fmt.Printf("Method: %s(%s) returns %s\n", n.Name, strings.Join(params, ", "), n.ReturnType)
		printNode(n.Body, nextIndent, true)
	case *model.Block:
		fmt.Println("{block}")
		for i, stmt := range n.Statements {
			printNode(stmt, nextIndent, i == len(n.Statements)-1)
		}
	case model.ForStmt:
		fmt.Println("FOR loop")
		printNode(n.Condition, nextIndent, false)
		printNode(n.Body, nextIndent, true)
	case model.IfStmt:
		fmt.Println("IF condition")
		printNode(n.Condition, nextIndent, false)
		printNode(n.Then, nextIndent, true)
	case model.BinaryExpr:
		fmt.Printf("BinaryExpr: %s\n", n.Op)
		printNode(n.Left, nextIndent, false)
		printNode(n.Right, nextIndent, true)
	case model.VarDecl:
		fmt.Printf("VarDecl: %s %s\n", n.Type, n.Name)
	case model.Ident:
		fmt.Printf("Identifier: %s\n", n.Name)
	case model.Literal:
		fmt.Printf("Literal: %s\n", n.Value)
	case model.ReturnStmt:
		fmt.Println("Return")
		printNode(n.Value, nextIndent, true)
	case model.SelectorExpr:
		fmt.Printf("Selector: .%s\n", n.Sel)
		printNode(n.X, nextIndent, true)
	case model.PostfixExpr:
		fmt.Printf("Postfix: %s\n", n.Op)
		printNode(n.X, nextIndent, true)
	case model.ArrayAccess:
		fmt.Println("ArrayAccess")
		printNode(n.Array, nextIndent, false)
		printNode(n.Index, nextIndent, true)
	default:
		fmt.Printf("%T\n", n)
	}
}
