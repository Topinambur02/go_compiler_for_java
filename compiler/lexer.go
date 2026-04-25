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

var delimiters = map[string]bool{
	".": true, ";": true, "{": true, "}": true, "(": true, ")": true, "[": true, "]": true, ",": true, " ": true,
}

func (c *Compiler) LexicalAnalysis(code string) ([]model.Token, error) {
	var table []model.Token
	line := 1
	col := 1

	for i := 0; i < len(code); {
		ch := rune(code[i])

		if ch == ' ' || ch == '\t' || ch == '\r' {
			i++
			col++
			continue
		}

		if ch == '\n' {
			i++
			line++
			col = 1
			continue
		}

		if unicode.IsLetter(ch) || ch == '_' {
			word := ""

			for i < len(code) && (unicode.IsLetter(rune(code[i])) || unicode.IsDigit(rune(code[i])) || code[i] == '_') {
				word += string(code[i])
				i++
				col++
			}

			if keywords[word] {
				table = append(table, model.Token{Type: constants.KEYWORD, Value: word})
			} else {
				table = append(table, model.Token{Type: constants.IDENTIFIER, Value: word})
			}

			continue
		}

		if unicode.IsDigit(ch) {
			startCol := col
			val := ""
			dotCount := 0
			hasLetter := false

			for i < len(code) {
				curr := rune(code[i])

				if unicode.IsDigit(curr) {
					val += string(curr)
					i++
					col++
				} else if curr == '.' {
					dotCount++
					val += string(curr)
					i++
					col++
				} else if unicode.IsLetter(curr) {
					hasLetter = true
					val += string(curr)
					i++
					col++
				} else {
					break
				}

			}

			if hasLetter {
				return nil, fmt.Errorf("Лексическая ошибка [Строка %d, Колонка %d]: буквы в цифровых константах (или идентификатор начинается с цифры) '%s'", line, startCol, val)
			}

			if dotCount > 1 {
				return nil, fmt.Errorf("Лексическая ошибка [Строка %d, Колонка %d]: некорректно оформленное число (множественные точки) '%s'", line, startCol, val)
			}

			table = append(table, model.Token{Type: constants.CONSTANT_INT, Value: val})

			continue
		}

		if i+1 < len(code) {
			twoChars := code[i : i+2]

			if operators[twoChars] {
				table = append(table, model.Token{Type: constants.OPERATOR, Value: twoChars})
				i += 2
				col += 2
				continue
			}
		}

		charStr := string(ch)

		if operators[charStr] {
			table = append(table, model.Token{Type: constants.OPERATOR, Value: charStr})
			i++
			col++
			continue
		}

		if delimiters[charStr] {
			table = append(table, model.Token{Type: constants.DELIMITER, Value: charStr})
			i++
			col++
			continue
		}

		return nil, fmt.Errorf("Лексическая ошибка [Строка %d, Колонка %d]: неизвестный оператор или символ '%s'", line, col, charStr)
	}

	return table, nil
}
