package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/clarkmcc/go-typescript"
	"github.com/dop251/goja"
)

var (
	err      error
	tsCode   []byte
	result   goja.Value
	filePath string
)

func main() {
	flag.StringVar(&filePath, "file", "scripts.ts", "path to the TypeScript file")
	flag.Parse()

	if tsCode, err = os.ReadFile(filePath); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	log.Printf("TypeScript Code: \n%s\n", string(tsCode))

	if result, err = typescript.Evaluate(bytes.NewReader(tsCode)); err != nil {
		log.Fatalln("Typescript:", err)
	}

	jsCode := result.String()

	log.Printf("JavaScript Code: \n%s\n", jsCode)

	jsFilePath := filepath.Join(filepath.Dir(filePath), filepath.Base(filePath[:len(filePath)-len(filepath.Ext(filePath))])+".js")

	if err = os.WriteFile(jsFilePath, []byte(jsCode), 0644); err != nil {
		log.Fatalln("Error writing to file: ", err)
	}

	fmt.Println("JavaScript file created (or overwritten):", jsFilePath)
}
