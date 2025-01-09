package main

import (
	"bytes"
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/clarkmcc/go-typescript"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/js"
)

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

	if tsCode, err = os.ReadFile(filePath); err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	jsCodeStr, err := typescript.TranspileString(string(tsCode))
	if err != nil {
		log.Fatalf("Error transpiling TypeScript: %v", err)
	}

	jsCode = []byte(jsCodeStr)

	if minifyFlag {
		m := minify.New()
		m.AddFunc("application/javascript", js.Minify)
		var minified bytes.Buffer
		if err = m.Minify("application/javascript", &minified, bytes.NewReader(jsCode)); err != nil {
			log.Fatalf("Minify: %v", err)
		}
		jsCode = minified.Bytes()
	}

	ext := ".js"
	if minifyFlag {
		ext = ".js.min"
	}

	var jsFilePath string
	if filePath != "scripts.ts" {
		jsFilePath = filepath.Join(filepath.Dir(filePath), filepath.Base(filePath[:len(filePath)-len(filepath.Ext(filePath))])+ext)
	} else {
		root, err := os.Getwd()
		if err != nil {
			log.Fatalf("Error getting working directory: %v", err)
		}

		jsFilePath = filepath.Join(root, filepath.Base(filePath[:len(filePath)-len(filepath.Ext(filePath))])+ext)
	}

	if err = os.WriteFile(jsFilePath, jsCode, 0o644); err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}

	log.Println("JavaScript file created (or overwritten) at...\n", jsFilePath)
}
