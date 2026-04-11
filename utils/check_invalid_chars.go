package utils

import "fmt"

func CheckInvalidChars(text string) error {
    for _, r := range text {
        if r < 32 && r != '\n' && r != '\r' && r != '\t' {
            return fmt.Errorf("invalid character: U+%04X (code %d)", r, r)
        }
    }

    return nil
}