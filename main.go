package main

import (
	"bytes"
	"flag"
	"fmt"
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
	strictFlag bool
)

func main() {
	flag.StringVar(&filePath, "file", "scripts.ts", "path to the TypeScript file")
	flag.BoolVar(&minifyFlag, "minify", false, "minify the JavaScript output")
	flag.BoolVar(&strictFlag, "strict", false, `add "use strict"; to the JavaScript output`)
	flag.Parse()

	if tsCode, err = os.ReadFile(filePath); err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	// log.Printf("TypeScript Code: \n%s\n----------\n", string(tsCode))

	jsCodeStr, err := typescript.TranspileString(string(tsCode))
	if err != nil {
		log.Fatalf("Error transpiling TypeScript: %v", err)
	}

	if strictFlag {
		if minifyFlag {
			jsCodeStr = `"use strict";` + jsCodeStr
		} else {
			jsCodeStr = `"use strict";\n` + jsCodeStr
		}
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

	jsFilePath := filepath.Join(filepath.Dir(filePath), filepath.Base(filePath[:len(filePath)-len(filepath.Ext(filePath))])+ext)

	if err = os.WriteFile(jsFilePath, jsCode, 0o644); err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}

	fmt.Println("JavaScript file created (or overwritten):", jsFilePath)
}
