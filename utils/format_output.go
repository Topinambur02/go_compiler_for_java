package utils

import (
	"fmt"
	"strings"

	"github.com/topinambur02/compiler/model"
)

func FormatOutput(tokens []model.Token) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("%-15s | %s\n", "Лексема", "Тип"))
	builder.WriteString("----------------+----------------------\n")

	var rawList []string

	for _, t := range tokens {
		builder.WriteString(fmt.Sprintf("%-15s | %s\n", t.Value, t.Type))
		rawList = append(rawList, fmt.Sprintf("(%s, %s)", t.Type, t.Value))
	}

	builder.WriteString("\n[")
	builder.WriteString(strings.Join(rawList, ", "))
	builder.WriteString("]\n\n")

	builder.WriteString(fmt.Sprintf("Лексический анализ завершён успешно. Обнаружено %d токенов. Ошибок не найдено.\n", len(tokens)))

	return builder.String()
}