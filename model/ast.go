package model

type Node any
type Expr any
type Statement any

type ImportDecl struct {
	Path string
}

type ClassDecl struct {
	Name    string
	Methods []MethodDecl
}

type Program struct {
	Imports []ImportDecl
	Classes []ClassDecl
}

type MethodDecl struct {
	Name       string
	ReturnType string
	Params     []Param
	Body       *Block
	IsStatic   bool
}

type Param struct {
	Type string
	Name string
}

type Block struct {
	Statements []Statement
}

type IfStmt struct {
	Condition Expr
	Then      *Block
}

type ForStmt struct {
	Init      Statement
	Condition Expr
	Post      Statement
	Body      *Block
}

type ReturnStmt struct {
	Value Expr
}

type VarDecl struct {
	Type  string
	Name  string
	Value Expr
}

type AssignStmt struct {
	Left  Expr
	Right Expr
}

type BinaryExpr struct {
	Left  Expr
	Op    string
	Right Expr
}

type Ident struct {
	Name string
}

type Literal struct {
	Value string
}

type ArrayAccess struct {
	Array Expr
	Index Expr
}

type CallExpr struct {
	Name string
	Args []Expr
}

type SelectorExpr struct {
	X   Expr
	Sel string
}

type PostfixExpr struct {
	X  Expr
	Op string
}

type ArrayLiteral struct {
    Elements []Expr
}