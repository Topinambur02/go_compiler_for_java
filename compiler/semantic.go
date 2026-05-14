package compiler

import (
	"fmt"
	"os"
	"strings"

	"github.com/topinambur02/compiler/model"
)

type Symbol struct {
	Name        string
	Type        string
	Declared    bool
	Initialized bool
	ParamCount  int
}

type Triad struct {
	Op   string
	Arg1 string
	Arg2 string
}

type Analyzer struct {
	symbols     []Symbol
	triads      []Triad
	errors      []string
	symbolIndex map[string]int
}

func (c *Compiler) SemanticAnalysis(ast *model.Program) {
	a := &Analyzer{
		symbols:     []Symbol{},
		triads:      []Triad{},
		errors:      []string{},
		symbolIndex: make(map[string]int),
	}

	for _, class := range ast.Classes {
		for _, method := range class.Methods {
			if _, exists := a.symbolIndex[method.Name]; exists {
				a.errors = append(a.errors, fmt.Sprintf("Повторное объявление функции: %s", method.Name))
			} else {
				a.symbolIndex[method.Name] = len(a.symbols)
				a.symbols = append(a.symbols, Symbol{
					Name:        method.Name,
					Type:        "function",
					Declared:    true,
					Initialized: true,
					ParamCount:  len(method.Params),
				})
			}
			for _, p := range method.Params {
				a.addSymbol(p.Name, p.Type, true)
			}
			a.analyzeBlock(method.Body)
		}
	}

	a.printResults()

	if len(a.errors) > 0 {
		fmt.Println("Компиляция прервана из-за ошибок.")
		os.Exit(1)
	}
}

func (a *Analyzer) addSymbol(name, typeName string, init bool) {
	if _, exists := a.symbolIndex[name]; exists {
		a.errors = append(a.errors, fmt.Sprintf("Повторное объявление переменной: %s", name)) // [cite: 20]
		return
	}
	a.symbolIndex[name] = len(a.symbols)
	a.symbols = append(a.symbols, Symbol{
		Name:        name,
		Type:        typeName,
		Declared:    true,
		Initialized: init,
	})
}

func (a *Analyzer) analyzeBlock(block *model.Block) {
	if block == nil {
		return
	}
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
			a.triads = append(a.triads, Triad{Op: ":=", Arg1: s.Name, Arg2: valRef}) // [cite: 15]
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
			if idx, exists := a.symbolIndex[ident.Name]; exists {
				a.symbols[idx].Initialized = true
			}
		}

	case model.ReturnStmt:
		valRef := a.analyzeExpr(s.Value)
		a.triads = append(a.triads, Triad{Op: "return", Arg1: valRef})

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
				if idx, exists := a.symbolIndex[ident.Name]; exists {
					a.symbols[idx].Initialized = true
				}
			}
		} else {
			a.errors = append(a.errors, fmt.Sprintf("Неожиданный оператор: %s", s.Op))
		}

	case model.Ident:
		if _, exists := a.symbolIndex[s.Name]; !exists {
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
		idx, exists := a.symbolIndex[e.Name]

		if !exists {
			a.errors = append(a.errors, fmt.Sprintf("Использование необъявленной переменной: %s", e.Name))
			return e.Name
		}

		if !a.symbols[idx].Initialized {
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
		_, exists := a.symbolIndex[e.Name]
		if !exists {
			a.errors = append(a.errors, fmt.Sprintf("Использование необъявленной переменной: %s", e.Name))
			return e.Name
		}
		return e.Name
	default:
		return a.analyzeExpr(expr)
	}
}

func (a *Analyzer) printResults() {
	if len(a.errors) > 0 {
		fmt.Println("Семантический анализ завершен с ошибками:")
		for _, err := range a.errors {
			fmt.Println("-", err)
		}
	} else {
		fmt.Printf("%-10s | %-10s | %-10s | %-10s\n", "Name", "Type", "Declared", "Initialized")
		fmt.Println("-----------+------------+------------+-------------")
		for _, s := range a.symbols {
			fmt.Printf("%-10s | %-10s | %-10v | %-10v\n", s.Name, s.Type, s.Declared, s.Initialized)
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
		if idx, exists := a.symbolIndex[funcName]; exists {
			sym := a.symbols[idx]
			if sym.Type == "function" && argsCount != sym.ParamCount {
				a.errors = append(a.errors, fmt.Sprintf("Неверное количество аргументов при вызове функции '%s': ожидается %d, передано %d", funcName, sym.ParamCount, argsCount))
			}
		}
	}
}
