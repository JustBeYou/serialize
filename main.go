package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"serialize/code"
	"strings"
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

	rootNode := code.CreateFileParser(*targetFile)
	typeNode, found := code.FindTargetTypeNode(rootNode, *targetType)

	if !found {
		panic("Type declaration not found in file")
	}
	fields := code.ParseStruct(typeNode)

	output := code.GenPackageHeaderAndImports(code.FindPackageName(rootNode)) + code.GenSerializationHeader(*targetType)
	for _, i := range fields {
		output += code.GenFieldSerialization(i)
	}
	output += code.GenSerializationFooter()

	outputFilePath := strings.Replace(*targetFile, ".go", ".ser.go", 1)
	outputFile, _ := os.Create(outputFilePath)

	outputFile.WriteString(output)
	outputFile.Close()
	c := exec.Command("gofmt", "-w", outputFilePath)
	c.Output()
}