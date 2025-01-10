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

var (
	tsCode             []byte
	jsCode             []byte
	filePath           string
	minifyFlag         bool
	streamMinThreshold int64
)

func main() {
	flag.StringVar(&filePath, "file", "scripts.ts", "Path to the TypeScript file.")
	flag.BoolVar(&minifyFlag, "minify", false, "Minify the JavaScript output.")
	flag.Int64Var(&streamMinThreshold, "stream", 10, "File streaming minimum threshold in megabytes.")
	flag.Parse()

	streamMinThreshold *= 1024 * 1024 // Convert streamMinThreshold from megabytes to bytes

	err := ProcessFile(filePath, minifyFlag)
	checkError(err, "processing file")
}

func ProcessFile(filePath string, minifyFlag bool) error {
	fileInfo, err := os.Stat(filePath)
	checkError(err, "getting file info")

	if fileInfo.Size() > streamMinThreshold {
		tsCode, err = readFileStream(filePath)
	} else {
		tsCode, err = os.ReadFile(filePath)
	}
	checkError(err, "reading TypeScript file")

	jsCodeStr, err := typescript.TranspileString(string(tsCode))
	checkError(err, "transpiling TypeScript to JavaScript")

	jsCode = []byte(jsCodeStr)

	if minifyFlag {
		jsCode, err = minifyJavaScript(jsCode)
		checkError(err, "minifying JavaScript")
	}

	jsFilePath, err := generateOutputFilePath(filePath, minifyFlag)
	checkError(err, "generating output file path")

	err = os.WriteFile(jsFilePath, jsCode, 0o644)
	checkError(err, "writing JavaScript file")

	log.Println("JavaScript file written to:\n", jsFilePath)
	return nil
}

func readFileStream(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	checkError(err, "opening file")
	defer file.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	checkError(err, "reading file stream")
	return buf.Bytes(), nil
}

func minifyJavaScript(jsCode []byte) ([]byte, error) {
	m := minify.New()
	m.AddFunc("application/javascript", js.Minify)
	var minified bytes.Buffer
	err := m.Minify("application/javascript", &minified, bytes.NewReader(jsCode))
	checkError(err, "minifying JavaScript")
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
	checkError(err, "getting current working directory")

	return filepath.Join(root, baseName+ext), nil
}

func checkError(err error, message string) {
	if err != nil {
		log.Fatalf("Error %s: %v", message, err)
	}
}
