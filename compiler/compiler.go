package compiler

import (
	"fmt"
	"os"

	"github.com/topinambur02/compiler/utils"
)

type Compiler struct {}

func (c *Compiler) Preprocessing(filename string) (string, error) {
	data, err := os.ReadFile(filename)

    if err != nil {
        return "", fmt.Errorf("cannot read file: %w", err)
    }
	
    content := string(data)

    if err := utils.CheckInvalidChars(content); err != nil {
        return "", err
    }

    noComments, err := utils.RemoveCommentsAndSpaces(content)

    if err != nil {
        return "", err
    }

    return noComments, nil
}
