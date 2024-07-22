package main

//go:generate go run . ./runny.schema.json

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/invopop/jsonschema"
	"github.com/simonwhitaker/runny/runny"
)

func main() {
	schema := jsonschema.Reflect(&runny.Config{})
	bytes, err := schema.MarshalJSON()
	if err != nil {
		panic(err)
	}

	// jsonschema doesn't support indenting, so we need to unmarshal/marshal with the json package
	var tempJsonObj map[string]interface{}
	err = json.Unmarshal(bytes, &tempJsonObj)
	if err != nil {
		panic(err)
	}

	indentedSchema, err := json.MarshalIndent(tempJsonObj, "", "    ")
	schemaString := string(indentedSchema[:])

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
