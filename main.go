package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/clarkmcc/go-typescript"
	"github.com/dop251/goja"
)

var (
	err    error
	tsCode []byte
	result goja.Value
	filePath string
)

func main() {
	flag.StringVar(&filePath, "file", "/scripts.ts", "path to the TypeScript file")
	flag.Parse()

	if tsCode, err = os.ReadFile(filePath); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	if result, err = typescript.Evaluate(string(tsCode)); err != nil {
		fmt.Println("Error evaluating TypeScript:", err)
		return
	}

	jsCode := result.String()

	jsFilePath := filepath.Join(filepath.Dir(filePath), filepath.Base(filePath[:len(filePath)-len(filepath.Ext(filePath))])+".js")

	if err = os.WriteFile(jsFilePath, []byte(jsCode), 0644); err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("JavaScript file created (or overwritten):", jsFilePath)
}