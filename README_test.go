package main

import (
	"os"
	"testing"

	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
	"github.com/simonwhitaker/runny/runny"
	"gopkg.in/yaml.v3"
)

func TestREADMECodeBlocks(t *testing.T) {
	p := parser.New()
	md, err := os.ReadFile("README.md")
	if err != nil {
		t.Fatalf("Got an error reading README.md (which is ironic): %v", err)
	}
	root := p.Parse(md)
	ast.WalkFunc(root, func(node ast.Node, entering bool) ast.WalkStatus {
		if codeBlock, ok := node.(*ast.CodeBlock); ok && entering {
			lang := string(codeBlock.Info)
			if lang == "yaml" {
				var c runny.Config
				err := yaml.Unmarshal(codeBlock.Literal, &c)
				if err != nil {
					t.Fatalf("YAML example can't be parsed:\n%s", string(codeBlock.Literal))
				}
				if len(c.Commands) == 0 {
					t.Fatalf("YAML example doesn't contain any commands:\n%s", string(codeBlock.Literal))
				}
			}
		}
		return ast.GoToNext
	})
}
