package main

import (
	"go/ast"
	"go/parser"
	"go/token"
)

func createFileParser(targetFile string) *ast.File {
	fileSet := token.NewFileSet()
	rootNode, err := parser.ParseFile(fileSet, targetFile, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	return rootNode
}

func findTargetTypeNode(rootNode ast.Node, targetType string) (ast.Node, bool) {
	var typeNode ast.Node
	found := false
	ast.Inspect(rootNode, func(n ast.Node) bool {
		t, ok := n.(*ast.TypeSpec)
		if ok && t.Name.Name == targetType {
			typeNode = n
			found = true
			return false
		}
		return true
	})
	return typeNode, found
}

func parseStruct(typeNode ast.Node) []StructField {
	var fields []StructField
	ast.Inspect(typeNode, func(x ast.Node) bool {
		s, ok := x.(*ast.StructType)
		if !ok {
			return true
		}

		for _, field := range s.Fields.List {
			fields = append(fields, StructField{
				field.Names[0].Name,
				field.Type.(*ast.Ident).Name,
			})
		}
		return false
	})
	return fields
}