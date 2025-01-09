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
	err        error
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
	if err != nil {
		log.Fatalf("Error getting file info: %v", err)
	}

	if fileInfo.Size() > streamThreshold {
		tsCode, err = readFileStream(filePath)
	} else {
		tsCode, err = os.ReadFile(filePath)
	}
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	jsCodeStr, err := typescript.TranspileString(string(tsCode))
	if err != nil {
		log.Fatalf("Error transpiling TypeScript: %v", err)
	}

	jsCode = []byte(jsCodeStr)

	if minifyFlag {
		jsCode, err = minifyJavaScript(jsCode)
		if err != nil {
			log.Fatalf("Minify: %v", err)
		}
	}

	jsFilePath, err := generateOutputFilePath(filePath, minifyFlag)
	if err != nil {
		log.Fatalf("Error generating output file path: %v", err)
	}

	if err = os.WriteFile(jsFilePath, jsCode, 0o644); err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}

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
