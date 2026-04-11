package main

import (
	"fmt"

	"github.com/topinambur02/compiler/compiler"
)

func main() {
	filename := "test.java"
	comp := compiler.Compiler{}
	result, err := comp.Preprocessing(filename)

	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}
