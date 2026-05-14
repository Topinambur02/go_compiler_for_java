package compiler

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/topinambur02/compiler/model"
)

type Symbol struct {
	Name        string
	Type        string
	Declared    bool
	Initialized bool
	ParamCount  int
	ReturnType  string
}

type Scope struct {
	symbols map[string]*Symbol
	outer   *Scope
}

type Triad struct {
	Op   string
	Arg1 string
	Arg2 string
}

type Analyzer struct {
	scopes      *Scope
	globalScope *Scope
	triads      []Triad
	errors      []string
	currentFunc *Symbol
	allScopes   []*Scope
}

func (c *Compiler) SemanticAnalysis(ast *model.Program) {
	a := &Analyzer{
		triads: []Triad{},
		errors: []string{},
		// globalScope: newScope(nil),
	}
	a.globalScope = a.newScope(nil)
	a.scopes = a.globalScope

	for _, class := range ast.Classes {
		for _, method := range class.Methods {
			if a.lookup(method.Name) != nil {
				a.errors = append(a.errors, fmt.Sprintf("Повторное объявление функции: %s", method.Name))
				continue
			}
			sym := &Symbol{
				Name:        method.Name,
				Type:        "function",
				Declared:    true,
				Initialized: true,
				ParamCount:  len(method.Params),
				ReturnType:  method.ReturnType,
			}
			a.globalScope.symbols[method.Name] = sym

			a.currentFunc = sym

			a.pushScope()
			for _, p := range method.Params {
				a.addSymbol(p.Name, p.Type, true)
			}
			a.analyzeBlock(method.Body)
			a.popScope()

			a.currentFunc = nil
		}
	}
	a.printResults()

	if len(a.errors) > 0 {
		fmt.Println("Компиляция прервана из-за ошибок.")
		os.Exit(1)
	}
}

func (a *Analyzer) newScope(outer *Scope) *Scope {
	s := &Scope{symbols: make(map[string]*Symbol), outer: outer}
	a.allScopes = append(a.allScopes, s)
	return s
}

func (a *Analyzer) pushScope() {
	a.scopes = a.newScope(a.scopes)
}

func (a *Analyzer) popScope() {
	if a.scopes != nil {
		a.scopes = a.scopes.outer
	}
}

func (a *Analyzer) lookup(name string) *Symbol {
	for scope := a.scopes; scope != nil; scope = scope.outer {
		if sym, ok := scope.symbols[name]; ok {
			return sym
		}
	}
	return nil
}

func (a *Analyzer) getType(expr model.Expr) string {
	switch e := expr.(type) {
	case model.Literal:
		if _, err := strconv.Atoi(e.Value); err == nil {
			return "int"
		}
		if strings.HasPrefix(e.Value, `"`) {
			return "String"
		}
		return "unknown"
	case model.Ident:
		if sym := a.lookup(e.Name); sym != nil {
			return sym.Type
		}
		return "unknown"
	case model.CallExpr:
		if sym := a.globalScope.symbols[e.Name]; sym != nil && sym.Type == "function" {
			return sym.ReturnType
		}
		return "unknown"
	default:
		return "unknown"
	}
}

func (a *Analyzer) addSymbol(name, typeName string, init bool) *Symbol {
	if a.lookup(name) != nil {
		a.errors = append(a.errors, fmt.Sprintf("Повторное объявление переменной: %s", name))
		return nil
	}
	sym := &Symbol{
		Name:        name,
		Type:        typeName,
		Declared:    true,
		Initialized: init,
	}
	a.scopes.symbols[name] = sym
	return sym
}

func (a *Analyzer) analyzeBlock(block *model.Block) {
	if block == nil {
		return
	}
	a.pushScope()
	defer a.popScope()
	for _, stmt := range block.Statements {
		a.analyzeStatement(stmt)
	}
}

func (a *Analyzer) analyzeStatement(stmt model.Statement) {
	switch s := stmt.(type) {
	case model.VarDecl:
		valRef := a.analyzeExpr(s.Value)
		a.addSymbol(s.Name, s.Type, s.Value != nil)
		if s.Value != nil {
			a.triads = append(a.triads, Triad{Op: ":=", Arg1: s.Name, Arg2: valRef})
		}

	case model.ForStmt:
		a.analyzeStatement(s.Init)
		a.analyzeExpr(s.Condition)
		a.analyzeExpr(s.Post)
		a.analyzeBlock(s.Body)

	case model.IfStmt:
		a.analyzeExpr(s.Condition)
		a.analyzeBlock(s.Then)

	case model.AssignStmt:
		leftRef := a.analyzeLValue(s.Left)
		rightRef := a.analyzeExpr(s.Right)
		a.triads = append(a.triads, Triad{Op: ":=", Arg1: leftRef, Arg2: rightRef})
		if ident, ok := s.Left.(model.Ident); ok {
			if sym := a.lookup(ident.Name); sym != nil {
				sym.Initialized = true
			}
		}

	case model.ReturnStmt:
		valRef := a.analyzeExpr(s.Value)
		a.triads = append(a.triads, Triad{Op: "return", Arg1: valRef})
		if a.currentFunc != nil {
			retType := a.getType(s.Value)
			if retType != "unknown" && retType != a.currentFunc.ReturnType && a.currentFunc.ReturnType != "void" {
				a.errors = append(a.errors, fmt.Sprintf(
					"Несоответствие типов: функция '%s' должна возвращать '%s', а возвращает '%s'",
					a.currentFunc.Name, a.currentFunc.ReturnType, retType))
			}
		}

	case model.CallExpr:
		a.checkFunctionCall(s.Name, len(s.Args))
		for _, arg := range s.Args {
			a.analyzeExpr(arg)
		}

	case model.BinaryExpr:
		if s.Op == "=" {
			leftRef := a.analyzeLValue(s.Left)
			rightRef := a.analyzeExpr(s.Right)
			a.triads = append(a.triads, Triad{Op: ":=", Arg1: leftRef, Arg2: rightRef})
			if ident, ok := s.Left.(model.Ident); ok {
				if sym := a.lookup(ident.Name); sym != nil {
					sym.Initialized = true
				}
			}
		} else {
			a.errors = append(a.errors, fmt.Sprintf("Неожиданный оператор: %s", s.Op))
		}

	case model.Ident:
		if sym := a.lookup(s.Name); sym != nil {
			a.errors = append(a.errors, fmt.Sprintf("Использование необъявленной переменной: %s", s.Name))
		}
	default:
		a.errors = append(a.errors, fmt.Sprintf("Неизвестный тип оператора: %T", s))
	}
}

func (a *Analyzer) analyzeExpr(expr model.Expr) string {
	if expr == nil {
		return ""
	}

	if s, ok := expr.(string); ok {
		return s
	}
	switch e := expr.(type) {
	case model.BinaryExpr:
		left := a.analyzeExpr(e.Left)
		right := a.analyzeExpr(e.Right)
		a.triads = append(a.triads, Triad{Op: e.Op, Arg1: left, Arg2: right})
		return fmt.Sprintf("^%d", len(a.triads))

	case model.ArrayAccess:
		arrName := a.analyzeExpr(e.Array)
		a.analyzeExpr(e.Index)
		return fmt.Sprintf("%s[]", arrName)

	case model.Ident:
		sym := a.lookup(e.Name)
		if sym == nil {
			a.errors = append(a.errors, fmt.Sprintf("Использование необъявленной переменной: %s", e.Name))
			return e.Name
		}
		if !sym.Initialized {
			a.errors = append(a.errors, fmt.Sprintf("Использование неинициализированной переменной: %s", e.Name))
		}
		return e.Name

	case model.Literal:
		return e.Value

	case model.SelectorExpr:
		obj := a.analyzeExpr(e.X)
		a.triads = append(a.triads, Triad{Op: ".", Arg1: obj, Arg2: e.Sel})
		return fmt.Sprintf("^%d", len(a.triads))

	case model.PostfixExpr:
		a.analyzeExpr(e.X)
		a.triads = append(a.triads, Triad{Op: e.Op, Arg1: fmt.Sprintf("^%d", len(a.triads))})
		return fmt.Sprintf("^%d", len(a.triads))

	case model.CallExpr:
		a.checkFunctionCall(e.Name, len(e.Args))
		funcRef := a.analyzeExpr(e.Name)
		for _, arg := range e.Args {
			a.analyzeExpr(arg)
		}
		a.triads = append(a.triads, Triad{Op: "call", Arg1: funcRef, Arg2: fmt.Sprintf("%d", len(e.Args))})
		return fmt.Sprintf("^%d", len(a.triads))

	default:
		a.errors = append(a.errors, fmt.Sprintf("Неизвестный тип выражения: %T", e))
		return ""
	}
}

func (a *Analyzer) analyzeLValue(expr model.Expr) string {
	switch e := expr.(type) {
	case model.Ident:
		if sym := a.lookup(e.Name); sym == nil {
			a.errors = append(a.errors, fmt.Sprintf("Использование необъявленной переменной: %s", e.Name))
			return e.Name
		}
		return e.Name
	default:
		return a.analyzeExpr(expr)
	}
}

func (a *Analyzer) collectAllSymbols() []Symbol {
	var result []Symbol

	for _, sym := range a.globalScope.symbols {
		result = append(result, *sym)
	}

	for _, s := range a.scopes.symbols {
		result = append(result, *s)
	}

	return result
}

func (a *Analyzer) printResults() {
	if len(a.errors) > 0 {
		fmt.Println("Семантический анализ завершен с ошибками:")
		for _, err := range a.errors {
			fmt.Println("-", err)
		}
	} else {
		fmt.Printf("%-15s | %-10s | %-10s | %-10s\n", "Name", "Type", "Declared", "Init")
		fmt.Println(strings.Repeat("-", 55))
		for _, scope := range a.allScopes {
			for _, sym := range scope.symbols {
				fmt.Printf("%-15s | %-10s | %-10v | %-10v\n",
					sym.Name, sym.Type, sym.Declared, sym.Initialized)
			}
		}
		fmt.Println()

		fmt.Println("Семантический анализ завершён успешно. Ошибок не найдено.")

		fmt.Println()

		for i, t := range a.triads {
			fmt.Printf("%d) (%s, %s, %s)\n", i+1, t.Op, t.Arg1, t.Arg2)
		}
	}
}

func (a *Analyzer) checkFunctionCall(nameExpr model.Expr, argsCount int) {
	var funcName string

	switch v := nameExpr.(type) {
	case model.Ident:
		funcName = v.Name
	case string:
		funcName = strings.Trim(v, "{}")
		parts := strings.Fields(funcName)
		if len(parts) > 0 {
			funcName = parts[len(parts)-1]
		}
	}

	if funcName != "" {
		if sym := a.lookup(funcName); sym != nil {
			if sym.Type == "function" && argsCount != sym.ParamCount {
				a.errors = append(a.errors, fmt.Sprintf("Неверное количество аргументов при вызове функции '%s': ожидается %d, передано %d", funcName, sym.ParamCount, argsCount))
			}
		}
	}
}
