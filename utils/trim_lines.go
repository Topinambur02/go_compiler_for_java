package utils

import (
	"bufio"
	"strings"
)

func TrimLines(text string) string {
    var result []string
    scanner := bufio.NewScanner(strings.NewReader(text))

    for scanner.Scan() {
        line := scanner.Text()
        trimmed := strings.TrimSpace(line)
		
        if trimmed != "" {
            result = append(result, trimmed)
        }
    }

    return strings.Join(result, "\n")
}