package main

import (
	"flag"
	"fmt"
	"go/ast"
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