package utils

import (
	"fmt"
	"regexp"
	"strings"
)

func RemoveCommentsAndSpaces(src string) (string, error) {
    singleLine := regexp.MustCompile(`//.*`)
    multiLine := regexp.MustCompile(`(?s)/\*.*?\*/`)
    openComment := regexp.MustCompile(`/\*`)
    closeComment := regexp.MustCompile(`\*/`)
    deleteSpaces := regexp.MustCompile(`[ \t]+`)
    emptyLine := regexp.MustCompile(`\n\s*\n`)
	
    openIndices := openComment.FindAllStringIndex(src, -1)
    closeIndices := closeComment.FindAllStringIndex(src, -1)

    if len(openIndices) != len(closeIndices) {
        return "", fmt.Errorf("unclosed multi-line comment")
    }

    cleaned := multiLine.ReplaceAllString(src, "")
    cleaned = singleLine.ReplaceAllString(cleaned, "")
    cleaned = deleteSpaces.ReplaceAllString(cleaned, " ")
    cleaned = emptyLine.ReplaceAllString(cleaned, "\n")

    cleaned = strings.TrimSpace(cleaned)

    return cleaned, nil
}