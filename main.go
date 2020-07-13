package main

import (
	"flag"
	"fmt"
	"github.com/JustBeYou/serialize/code"
	"go/ast"
	"os"
	"os/exec"
	"strings"
)

const defaultSerializer = ""

func main() {
	targetFile := flag.String("file", "", "Target file")
	targetType := flag.String("type", "", "Target struct to be 'serializable'")
	targetTable := flag.Bool("table", false, "True if a type table should be generated")
	serializer := flag.String("serializer", "", "Target serializer name (optional)")
	flag.Parse()

	if *targetFile == "" {
		flag.Usage()
		return
	}

	if *targetType != "" {
		var serializersList []string
		if *serializer != "" {
			serializersList = strings.Split(*serializer, ",")
		}
		serializersList = append(serializersList, defaultSerializer)
		generateStructSerializers(serializersList, *targetFile, *targetType)
	} else if *targetTable {
		// TODO: implement generateTypeTable(*targetFile)
		generateTypeTable(*targetFile)
		return
	} else {
		flag.Usage()
		return
	}
}

func generateTypeTable(targetFile string) {
	rootNode := code.CreateFileParser(targetFile)
	output := code.GenPackageHeaderAndImports(code.FindPackageName(rootNode), "")

	tableNode, found := code.FindTypeTableNode(rootNode)
	if !found {
		panic("Type table not found in file")
	}

	values := tableNode.Values[0].(*ast.CompositeLit)
	var typeTable []string
	for _, v := range values.Elts {
		key := v.(*ast.KeyValueExpr).Key.(*ast.BasicLit).Value
		typeTable = append(typeTable, key[1:len(key)-1])
	}

	output += code.GenTypeTable(typeTable)

	outputFilePath := strings.Replace(targetFile, ".go", ".type.ser.go", 1)
	outputFile, _ := os.Create(outputFilePath)
	writeAndFormatFile(outputFile, outputFilePath, output)
}

func generateStructSerializers(serializersList []string, targetFile, targetType string) {

	rootNode := code.CreateFileParser(targetFile)
	typeNode, found := code.FindTargetTypeNode(rootNode, targetType)

	if !found {
		panic("Type declaration not found in file")
	}
	fields := code.ParseStruct(typeNode, serializersList)

	output := code.GenPackageHeaderAndImports(code.FindPackageName(rootNode), targetType)
	
	capabilities := code.UsedCapabilities{
		ArraySerialize:      false,
		InterfaceSerialize:  false,
		CustomTypeSerialize: false,
	}
	for _, serializerName := range serializersList {
		tempOutput := ""
		for _, i := range fields {
			tempOutput += code.GenFieldSerialization(serializerName, i, &capabilities)
		}
		tempOutput = code.GenSerializationHeader(serializerName, targetType, capabilities) + tempOutput
		tempOutput += code.GenSerializationFooter()

		output += tempOutput
		tempOutput = ""

		for _, i := range fields {
			tempOutput += code.GenFieldUnserialization(serializerName, i)
		}
		tempOutput = code.GenUnserializationHeader(serializerName, targetType, capabilities) + tempOutput
		tempOutput += code.GenUnserializationFooter()

		output += tempOutput
	}

	outputFilePath := strings.Replace(targetFile, ".go", fmt.Sprintf(".%s.ser.go", strings.ToLower(targetType)), 1)
	outputFile, _ := os.Create(outputFilePath)
	writeAndFormatFile(outputFile, outputFilePath, output)
}

func writeAndFormatFile(outputFile *os.File, outputFilePath, output string) {
	_, _ = outputFile.WriteString(output)
	_ = outputFile.Close()
	c := exec.Command("gofmt", "-w", outputFilePath)
	_, _ = c.Output()

	c = exec.Command("goimports", "-w", outputFilePath)
	_, _ = c.Output()
}