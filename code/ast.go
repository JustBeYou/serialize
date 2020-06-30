package code

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func CreateFileParser(targetFile string) *ast.File {
	fileSet := token.NewFileSet()
	rootNode, err := parser.ParseFile(fileSet, targetFile, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	return rootNode
}

func FindTargetTypeNode(rootNode ast.Node, targetType string) (ast.Node, bool) {
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

func ParseStruct(typeNode ast.Node) []StructField {
	var fields []StructField
	ast.Inspect(typeNode, func(x ast.Node) bool {
		s, ok := x.(*ast.StructType)
		if !ok {
			return true
		}

		for _, field := range s.Fields.List {
			asArray, isArray := field.Type.(*ast.ArrayType)
			var typeName string
			if isArray {
				fmt.Printf("Array of %s's\n", asArray.Elt)
				typeName = fmt.Sprintf("%s", asArray.Elt)
			} else {
				typeName = field.Type.(*ast.Ident).Name
			}

			fields = append(fields, StructField{
				field.Names[0].Name,
				typeName,
				isArray,
			})
		}
		return false
	})
	return fields
}

func FindPackageName(rootNode ast.Node) string {
	packageName := ""
	ast.Inspect(rootNode, func(n ast.Node) bool {
		i, ok := n.(*ast.Ident)
		if ok && packageName == "" && i.Name != "" {
			packageName = i.Name
		}
		return true
	})
	return packageName
}