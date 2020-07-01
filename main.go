package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"serialize/code"
	"strings"
)

const defaultSerializer = ""

func main() {
	targetFile := flag.String("file", "", "Target file")
	targetType := flag.String("type", "", "Target struct to be 'serializable'")
	serializer := flag.String("serializer", "", "Target serializer name (optional)")
	flag.Parse()

	if *targetType == "" || *targetFile == "" {
		flag.Usage()
		return ;
	}

	var serializersList []string
	if *serializer != "" {
		serializersList = strings.Split(*serializer, ",")
	}
	serializersList = append(serializersList, defaultSerializer)

	rootNode := code.CreateFileParser(*targetFile)
	typeNode, found := code.FindTargetTypeNode(rootNode, *targetType)

	if !found {
		panic("Type declaration not found in file")
	}
	fields := code.ParseStruct(typeNode, serializersList)

	output := code.GenPackageHeaderAndImports(code.FindPackageName(rootNode))

	for _, serializerName := range serializersList {
		output += code.GenSerializationHeader(serializerName, *targetType)
		for _, i := range fields {
			output += code.GenFieldSerialization(serializerName, i)
		}
		output += code.GenSerializationFooter()

		output += code.GenUnserializationHeader(serializerName, *targetType)
		for _, i := range fields {
			output += code.GenFieldUnserialization(serializerName, i)
		}
		output += code.GenUnserializationFooter()
	}

	outputFilePath := strings.Replace(*targetFile, ".go", fmt.Sprintf(".%s.ser.go", strings.ToLower(*targetType)), 1)
	outputFile, _ := os.Create(outputFilePath)

	outputFile.WriteString(output)
	outputFile.Close()
	c := exec.Command("gofmt", "-w", outputFilePath)
	c.Output()
}