package utils

import (
	"fmt"
	"strings"

	"github.com/topinambur02/compiler/model"
)

func FormatOutput(tokens []model.Token) string {
	var builder strings.Builder

	fmt.Fprintf(&builder, "%-15s | %s\n", "Лексема", "Тип")
	builder.WriteString("----------------+----------------------\n")

	for _, t := range tokens {
		fmt.Fprintf(&builder, "%-15s | %s\n", t.Value, t.Type)
	}

	fmt.Fprintf(&builder, "\nЛексический анализ завершён успешно. Обнаружено %d токенов. Ошибок не найдено.\n", len(tokens))
	return builder.String()
}
