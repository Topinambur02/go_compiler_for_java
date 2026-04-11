package utils

import (
	"fmt"
	"regexp"
)

func RemoveComments(src string) (string, error) {
    singleLine := regexp.MustCompile(`//.*`)
    multiLine := regexp.MustCompile(`/\*.*?\*/`)
    openComment := regexp.MustCompile(`/\*`)
    closeComment := regexp.MustCompile(`\*/`)
	
    openIndices := openComment.FindAllStringIndex(src, -1)
    closeIndices := closeComment.FindAllStringIndex(src, -1)

    if len(openIndices) != len(closeIndices) {
        return "", fmt.Errorf("unclosed multi-line comment")
    }

    cleaned := multiLine.ReplaceAllString(src, "")
    cleaned = singleLine.ReplaceAllString(cleaned, "")

    return cleaned, nil
}