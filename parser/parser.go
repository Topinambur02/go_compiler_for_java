package parser

import (
	"fmt"

	"github.com/topinambur02/compiler/constants"
	"github.com/topinambur02/compiler/model"
)

type Parser struct {
	Tokens []model.Token
	Pos    int
	Errors []string
}

func (p *Parser) peek() model.Token {
	if p.Pos >= len(p.Tokens) {
		return model.Token{Type: "EOF", Value: "EOF"}
	}
	return p.Tokens[p.Pos]
}

func (p *Parser) consume() model.Token {
	t := p.peek()
	if p.Pos < len(p.Tokens) {
		p.Pos++
	}
	return t
}

func (p *Parser) reportError(errType, expected, actual string) {
	msg := fmt.Sprintf("Ошибка: %s. Позиция: %d. Ожидалось: %s, получено: '%s'", errType, p.Pos, expected, actual)
	p.Errors = append(p.Errors, msg)
}

func (p *Parser) PrintErrors() {
	if len(p.Errors) == 0 {
		fmt.Println("Синтаксических ошибок не найдено.")
		return
	}
	fmt.Printf("Найдено %d синтаксических ошибок:\n", len(p.Errors))
	for _, err := range p.Errors {
		fmt.Println(" -", err)
	}
}

func (p *Parser) expect(expectedType constants.TokenType) model.Token {
	t := p.peek()
	if t.Type != expectedType {
		p.reportError("Неожиданный тип токена", string(expectedType), fmt.Sprintf("%s (%s)", t.Type, t.Value))
		return model.Token{Type: expectedType, Value: "<missing>"}
	}
	return p.consume()
}

func (p *Parser) expectValue(val string) {
	t := p.peek()
	if t.Value != val {
		errType := "Неожиданный токен в текущей позиции"
		if val == ";" || val == "," {
			errType = "Отсутствие обязательного разделителя"
		} else if val == "{" || val == "}" || val == "(" || val == ")" || val == "[" || val == "]" {
			errType = "Непарные операторные скобки"
		}
		p.reportError(errType, val, t.Value)
		return
	}
	p.consume()
}

func (p *Parser) Parse() model.Program {
	program := model.Program{}
	for p.Pos < len(p.Tokens) {
		token := p.peek()
		if token.Type == "EOF" {
			break
		}

		if token.Value == "import" {
			program.Imports = append(program.Imports, p.parseImport())
		} else if token.Value == "public" || token.Value == "class" {
			program.Classes = append(program.Classes, p.parseClass())
		} else {
			p.reportError("Нарушение структуры программы (токен вне класса или импорта)", "import или class", token.Value)
			p.consume()
		}
	}
	return program
}

func (p *Parser) parseImport() model.ImportDecl {
	p.expectValue("import")
	path := ""
	for p.peek().Value != ";" && p.peek().Type != "EOF" {
		path += p.consume().Value
	}
	p.expectValue(";")
	return model.ImportDecl{Path: path}
}

func (p *Parser) parseClass() model.ClassDecl {
	if p.peek().Value == "public" {
		p.consume()
	}
	p.expectValue("class")

	classNameToken := p.expect(constants.IDENTIFIER)
	className := classNameToken.Value
	if className == "<missing>" {
		className = "UnknownClass"
	}

	p.expectValue("{")
	class := model.ClassDecl{Name: className}

	for p.peek().Value != "}" && p.peek().Type != "EOF" {
		if p.isMethodStart() {
			class.Methods = append(class.Methods, p.parseMethod())
		} else {
			p.reportError("Нарушение структуры программы (ожидался метод)", "модификатор метода или тип", p.peek().Value)
			p.consume()
		}
	}
	p.expectValue("}")
	return class
}

func (p *Parser) isMethodStart() bool {
	return p.peek().Value == "public" || p.peek().Value == "static" || p.peek().Value == "private"
}

func (p *Parser) parseType() string {
	t := p.consume()
	typeName := t.Value

	for p.peek().Value == "[" {
		if p.Pos+1 < len(p.Tokens) && p.Tokens[p.Pos+1].Value == "]" {
			p.consume()
			p.consume()
			typeName += "[]"
		} else {
			p.reportError("Непарные скобки (ошибка объявления массива)", "]", p.peek().Value)
			p.consume()
			break
		}
	}

	return typeName
}

func (p *Parser) parseMethod() model.MethodDecl {
	method := model.MethodDecl{}

	for p.peek().Value == "public" || p.peek().Value == "static" || p.peek().Value == "private" {
		t := p.consume()
		if t.Value == "static" {
			method.IsStatic = true
		}
	}

	method.ReturnType = p.parseType()
	method.Name = p.expect(constants.IDENTIFIER).Value

	p.expectValue("(")
	for p.peek().Value != ")" && p.peek().Type != "EOF" {
		pType := p.parseType()
		pName := p.expect(constants.IDENTIFIER).Value
		method.Params = append(method.Params, model.Param{Type: pType, Name: pName})

		if p.peek().Value == "," {
			p.consume()
		} else if p.peek().Value != ")" {
			p.reportError("Отсутствие обязательного разделителя", ",", p.peek().Value)
			break
		}
	}
	p.expectValue(")")

	method.Body = p.parseBlock()

	return method
}

func (p *Parser) parseBlock() *model.Block {
	p.expectValue("{")
	block := &model.Block{}
	for p.peek().Value != "}" && p.peek().Type != "EOF" {
		startPos := p.Pos
		block.Statements = append(block.Statements, p.parseStatement())
		if p.Pos == startPos {
			p.consume()
		}
	}
	p.expectValue("}")
	return block
}

func (p *Parser) parseStatement() model.Statement {
	token := p.peek()

	if token.Type == constants.KEYWORD {
		switch token.Value {
		case "for":
			return p.parseFor()
		case "if":
			return p.parseIf()
		case "return":
			p.consume()
			val := p.parseExpression()
			p.expectValue(";")
			return model.ReturnStmt{Value: val}
		}
	}

	if p.isTypeStart() {
		typ := p.parseType()
		if p.peek().Type == constants.IDENTIFIER {
			name := p.consume().Value
			var init model.Expr
			if p.peek().Value == "=" {
				p.consume()
				init = p.parseExpression()
			}
			p.expectValue(";")
			return model.VarDecl{Type: typ, Name: name, Value: init}
		} else {
			p.reportError("Нарушение структуры (ожидалось имя переменной)", "IDENTIFIER", p.peek().Value)
		}
	}

	expr := p.parseExpression()
	if p.peek().Value == "=" {
		p.consume()
		right := p.parseExpression()
		p.expectValue(";")
		return model.AssignStmt{Left: expr, Right: right}
	}
	p.expectValue(";")
	return expr
}

func (p *Parser) isTypeStart() bool {
	t := p.peek()

	if t.Value == "int" || t.Value == "void" || t.Value == "String" || t.Value == "double" {
		return true
	}

	if t.Type == constants.IDENTIFIER {
		if p.Pos+1 < len(p.Tokens) {
			next := p.Tokens[p.Pos+1]
			if next.Type == constants.IDENTIFIER {
				return true
			}
			if next.Value == "[" && p.Pos+2 < len(p.Tokens) {
				if p.Tokens[p.Pos+2].Value == "]" {
					return true
				}
			}
		}
	}
	return false
}

func (p *Parser) parseExpression() model.Expr {
	left := p.parsePrimary()

	next := p.peek()
	if next.Type == constants.OPERATOR && next.Value != "++" && next.Value != "--" {
		op := p.consume().Value
		right := p.parseExpression()
		return model.BinaryExpr{Left: left, Op: op, Right: right}
	}

	return left
}

func (p *Parser) parsePrimary() model.Expr {
	t := p.peek()

	if t.Value == "{" {
		p.consume()
		for p.peek().Value != "}" && p.peek().Type != "EOF" {
			p.consume()
		}
		p.expectValue("}")
		return model.Literal{Value: "array_literal"}
	}

	t = p.consume()
	var node model.Expr

	switch t.Type {
	case constants.CONSTANT_INT:
		node = model.Literal{Value: t.Value}
	case constants.IDENTIFIER:
		node = model.Ident{Name: t.Value}
	case constants.KEYWORD:
		node = model.Ident{Name: t.Value}
	default:
		p.reportError("Неожиданный токен в выражении", "выражение (идентификатор или константа)", t.Value)
		return model.Ident{Name: "<error>"}
	}

	for {
		next := p.peek()
		if next.Value == "." {
			p.consume()
			selToken := p.expect(constants.IDENTIFIER)
			node = model.SelectorExpr{X: node, Sel: selToken.Value}
		} else if next.Value == "[" {
			p.consume()
			index := p.parseExpression()
			p.expectValue("]")
			node = model.ArrayAccess{Array: node, Index: index}
		} else if next.Value == "(" {
			p.consume()
			p.skipUntil(")")
			p.expectValue(")")
			node = model.CallExpr{Name: fmt.Sprintf("%v", node)}
		} else if next.Value == "++" || next.Value == "--" {
			op := p.consume().Value
			node = model.PostfixExpr{X: node, Op: op}
		} else {
			break
		}
	}
	return node
}

func (p *Parser) parseFor() model.Statement {
	p.expectValue("for")
	p.expectValue("(")
	init := p.parseStatement()
	cond := p.parseExpression()
	p.expectValue(";")
	post := p.parseExpression()
	p.expectValue(")")
	body := p.parseBlock()
	return model.ForStmt{
		Init:      init,
		Condition: cond,
		Post:      post,
		Body:      body,
	}
}

func (p *Parser) parseIf() model.Statement {
	p.expectValue("if")
	p.expectValue("(")
	cond := p.parseExpression()
	p.expectValue(")")
	then := p.parseBlock()
	return model.IfStmt{Condition: cond, Then: then}
}

func (p *Parser) skipUntil(val string) {
	for p.peek().Value != val && p.peek().Type != "EOF" {
		p.Pos++
	}
}
