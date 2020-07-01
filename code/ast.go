package code

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
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

func ParseStruct(typeNode ast.Node, serializersList []string) []StructField {
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
				typeName = fmt.Sprintf("%s", asArray.Elt)
			} else {
				typeName = field.Type.(*ast.Ident).Name
			}

			tagOptions := make(FieldOptions)
			if field.Tag != nil {
				tag := reflect.StructTag(field.Tag.Value[1:len(field.Tag.Value)-1])
				for _, serializerName := range serializersList {
					value, ok := tag.Lookup(serializerName)
					if ok {
						if tagOptions[serializerName] == nil {
							tagOptions[serializerName] = make(SerializerOptions)
						}
						tagOptions[serializerName][value] = true
					}
				}
			}

			fields = append(fields, StructField{
				field.Names[0].Name,
				typeName,
				isArray,
				tagOptions,
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