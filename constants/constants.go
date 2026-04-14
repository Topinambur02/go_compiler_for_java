package constants

type TokenType string

const (
	KEYWORD       TokenType = "KEYWORD"
	IDENTIFIER    TokenType = "IDENTIFIER"
	CONSTANT_INT  TokenType = "CONSTANT_INT"
	CONSTANT_REAL TokenType = "CONSTANT_REAL"
	CONSTANT_STR  TokenType = "CONSTANT_STR"
	OPERATOR      TokenType = "OPERATOR"
	DELIMITER     TokenType = "DELIMITER"
)
