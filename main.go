package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

func main() {
	targetFile := flag.String("file", "", "Target file")
	targetType := flag.String("type", "", "Target struct to be 'serializable'")
	serializer := flag.String("serializer", "", "Target serializer name (optional)")
	flag.Parse()

	if *targetType == "" || *targetFile == "" {
		flag.Usage()
		return ;
	}

	if *serializer != "" {
		fmt.Fprintf(os.Stderr, "-serializer not implemented yet\n")
		return ;
	}

	rootNode := createFileParser(*targetFile)
	typeNode, found := findTargetTypeNode(rootNode, *targetType)

	if !found {
		panic("Type declaration not found in file")
	}
	fmt.Printf("Found %s\n", typeNode.(*ast.TypeSpec).Name.Name)
	fields := parseStruct(typeNode)
	fmt.Printf("Fields: %v\n", fields)
}

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

type StructField struct {
	name string
	typeName string
}