package model

import "github.com/topinambur02/compiler/constants"

type Token struct {
	Type  constants.TokenType
	Value string
}
