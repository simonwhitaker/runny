package main

//go:generate go run . ./runny.schema.json

import (
	"fmt"
	"os"

	"github.com/simonwhitaker/runny/runny"
)

func main() {
	schemaString, err := runny.GenerateSchema()
	if err != nil {
		panic(err)
	}

	if len(os.Args) > 1 {
		filename := os.Args[1]
		f, err := os.Create(filename)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		f.WriteString(schemaString)
	} else {
		fmt.Println(schemaString)
	}
}
