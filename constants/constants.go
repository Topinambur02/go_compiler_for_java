package constants

type TokenType string

const (
	KEYWORD       TokenType = "KEYWORD"
	IDENTIFIER    TokenType = "IDENTIFIER"
	CONSTANT_INT  TokenType = "CONSTANT_INT"
	OPERATOR      TokenType = "OPERATOR"
	DELIMITER     TokenType = "DELIMITER"
)
