package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/clarkmcc/go-typescript"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/js"
)

const bytesPerMegabyte = 1024 * 1024

var (
	filePath           *string  = flag.String("file", "scripts.ts", "Path to the TypeScript file.")
	minifyFlag         *bool    = flag.Bool("minify", false, "Minify the JavaScript output.")
	streamMinThreshold *float64 = flag.Float64("stream", 2.5, "TypeScript file minimum streaming threshold in megabytes.")
)

func main() {
	flag.Parse()

	// Convert streamMinThreshold from megabytes to bytes
	*streamMinThreshold *= bytesPerMegabyte

	if err := ProcessFile(); err != nil {
		log.Fatalf("Error processing file: %v", err)
	}
}

func ProcessFile() error {
	tsCode, err := readFile()
	if err != nil {
		return err
	}

	jsCode, err := transpileTypeScript(tsCode)
	if err != nil {
		return err
	}

	if *minifyFlag {
		jsCode, err = minifyJavaScript(jsCode)
		if err != nil {
			return err
		}
	}

	jsFilePath, err := generateOutputFilePath()
	if err != nil {
		return err
	}

	if err := os.WriteFile(jsFilePath, jsCode, 0o644); err != nil {
		return fmt.Errorf("writing JavaScript file: %w", err)
	}

	log.Println("JavaScript file written to:", jsFilePath)
	return nil
}

func readFileStream() ([]byte, error) {
	var buf bytes.Buffer

	file, err := os.Open(*filePath)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(&buf, file); err != nil {
		return nil, fmt.Errorf("reading file stream: %w", err)
	}
	return buf.Bytes(), nil
}

func readFile() ([]byte, error) {
	fileInfo, err := os.Stat(*filePath)
	if err != nil {
		return nil, fmt.Errorf("getting file info: %w", err)
	}

	// If file size exceeds the streamMinThreshold, read the file in a streaming manner.
	// This approach is more memory-efficient for large files as it avoids loading the entire file into memory at once.
	if fileInfo.Size() > int64(*streamMinThreshold) {
		return readFileStream()
	}

	// Otherwise, read the entire file into memory.
	return os.ReadFile(*filePath)
}

func transpileTypeScript(tsCode []byte) ([]byte, error) {
	jsCodeStr, err := typescript.TranspileString(string(tsCode))
	if err != nil {
		return nil, fmt.Errorf("transpiling TypeScript to JavaScript: %w", err)
	}
	return []byte(jsCodeStr), nil
}

func minifyJavaScript(jsCode []byte) ([]byte, error) {
	var minified bytes.Buffer

	m := minify.New()
	m.AddFunc("application/javascript", js.Minify)

	if err := m.Minify("application/javascript", &minified, bytes.NewReader(jsCode)); err != nil {
		return nil, fmt.Errorf("minifying JavaScript: %w", err)
	}
	return minified.Bytes(), nil
}

func generateOutputFilePath() (string, error) {
	ext := ".js"
	if *minifyFlag {
		ext = ".js.min"
	}

	fileExt := filepath.Ext(*filePath)
	baseName := filepath.Base((*filePath)[:len(*filePath)-len(fileExt)])
	if *filePath != "scripts.ts" {
		return filepath.Join(filepath.Dir(*filePath), baseName+ext), nil
	}

	root, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("getting current working directory: %w", err)
	}

	return filepath.Join(root, baseName+ext), nil
}
