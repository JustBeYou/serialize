package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
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

	rootNode := createFileParser(*targetFile)
	typeNode, found := findTargetTypeNode(rootNode, *targetType)

	if !found {
		panic("Type declaration not found in file")
	}
	fields := parseStruct(typeNode)

	output := GenPackageHeaderAndImports("main") + GenSerializationHeader(*targetType)
	for _, i := range fields {
		output += GenFieldSerialization(i)
	}
	output += GenSerializationFooter()

	outputFilePath := strings.Replace(*targetFile, ".go", ".ser.go", 1)
	outputFile, _ := os.Create(outputFilePath)

	outputFile.WriteString(output)
	outputFile.Close()
	c := exec.Command("gofmt", "-w", outputFilePath)
	c.Output()
}