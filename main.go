package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/clarkmcc/go-typescript"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/js"
)

const streamThreshold = 10 * 1024 * 1024 // 10 MB

var (
	tsCode     []byte
	jsCode     []byte
	filePath   string
	minifyFlag bool
)

func main() {
	flag.StringVar(&filePath, "filePath", "scripts.ts", "path to the TypeScript file")
	flag.BoolVar(&minifyFlag, "minify", false, "minify the JavaScript output")
	flag.Parse()

	fileInfo, err := os.Stat(filePath)
	checkError(err, "Error getting file info")

	if fileInfo.Size() > streamThreshold {
		tsCode, err = readFileStream(filePath)
	} else {
		tsCode, err = os.ReadFile(filePath)
	}
	checkError(err, "Error reading file")

	jsCodeStr, err := typescript.TranspileString(string(tsCode))
	checkError(err, "Error transpiling TypeScript")

	jsCode = []byte(jsCodeStr)

	if minifyFlag {
		jsCode, err = minifyJavaScript(jsCode)
		checkError(err, "Minify")
	}

	jsFilePath, err := generateOutputFilePath(filePath, minifyFlag)
	checkError(err, "Error generating output file path")

	err = os.WriteFile(jsFilePath, jsCode, 0o644)
	checkError(err, "Error writing to file")

	log.Println("JavaScript file created (or overwritten) at...\n", jsFilePath)
}

func readFileStream(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, file); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func minifyJavaScript(jsCode []byte) ([]byte, error) {
	m := minify.New()
	m.AddFunc("application/javascript", js.Minify)
	var minified bytes.Buffer
	if err := m.Minify("application/javascript", &minified, bytes.NewReader(jsCode)); err != nil {
		return nil, err
	}
	return minified.Bytes(), nil
}

func generateOutputFilePath(filePath string, minifyFlag bool) (string, error) {
	ext := ".js"
	if minifyFlag {
		ext = ".js.min"
	}

	baseName := filepath.Base(filePath[:len(filePath)-len(filepath.Ext(filePath))])
	if filePath != "scripts.ts" {
		return filepath.Join(filepath.Dir(filePath), baseName+ext), nil
	}

	root, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Join(root, baseName+ext), nil
}

func checkError(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
}
