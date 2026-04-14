package compiler

import (
	"fmt"
	"unicode"

	"github.com/topinambur02/compiler/constants"
	"github.com/topinambur02/compiler/model"
)

var keywords = map[string]bool{
	"import": true, "public": true, "class": true, "static": true, "void": true, "int": true, "for": true, "if": true, "return": true,
}

var operators = map[string]bool{
	"=": true, "<": true, ">": true, "-": true, "+": true, "++": true,
}

var delimiters = map[byte]bool{
	'.': true, ';': true, '{': true, '}': true, '(': true, ')': true, '[': true, ']': true, ',': true,
}

func (c *Compiler) LexicalAnalysis(code string) ([]model.Token, error) {
	var tokens []model.Token
	length := len(code)
	i := 0

	for i < length {
		char := code[i]

		if unicode.IsSpace(rune(char)) {
			i++
			continue
		}

		if unicode.IsLetter(rune(char)) || char == '_' {
			start := i

			for i < length && (unicode.IsLetter(rune(code[i])) || unicode.IsDigit(rune(code[i])) || code[i] == '_') {
				i++
			}

			val := code[start:i]

			if keywords[val] {
				tokens = append(tokens, model.Token{Type: constants.KEYWORD, Value: val})
			} else {
				tokens = append(tokens, model.Token{Type: constants.IDENTIFIER, Value: val})
			}

			continue
		}

		if unicode.IsDigit(rune(char)) {
			start := i
			dotCount := 0
			hasLetter := false

			for i < length && !unicode.IsSpace(rune(code[i])) && !delimiters[code[i]] && !operators[string(code[i])] {
				if code[i] == '.' {
					dotCount++
				} else if unicode.IsLetter(rune(code[i])) {
					hasLetter = true
				}

				i++
			}

			val := code[start:i]

			if hasLetter {
				return nil, fmt.Errorf("Лексическая ошибка: буква в цифровой константе '%s'", val)
			}

			if dotCount == 1 {
				tokens = append(tokens, model.Token{Type: constants.CONSTANT_REAL, Value: val})
			} else if dotCount > 1 {
				return nil, fmt.Errorf("Лексическая ошибка: некорректно оформленное число (множественные точки) '%s'", val)
			} else {
				tokens = append(tokens, model.Token{Type: constants.CONSTANT_INT, Value: val})
			}

			continue
		}

		if char == '"' {
			start := i
			i++
			closed := false

			for i < length {
				if code[i] == '"' {
					closed = true
					i++
					break
				}

				if code[i] == '\n' {
					break
				}

				i++
			}
			val := code[start:i]

			if !closed {
				return nil, fmt.Errorf("Лексическая ошибка: незакрытый строковый литерал '%s'", val)
			}

			tokens = append(tokens, model.Token{Type: constants.CONSTANT_STR, Value: val})

			continue
		}

		if i+1 < length && operators[code[i:i+2]] {
			tokens = append(tokens, model.Token{Type: constants.OPERATOR, Value: code[i : i+2]})
			i += 2
			continue
		}

		if operators[string(char)] {
			tokens = append(tokens, model.Token{Type: constants.OPERATOR, Value: string(char)})
			i++
			continue
		}

		if delimiters[char] {
			tokens = append(tokens, model.Token{Type: constants.DELIMITER, Value: string(char)})
			i++
			continue
		}

		return nil, fmt.Errorf("Лексическая ошибка: недопустимый символ '%c'", char)
	}

	return tokens, nil
}
