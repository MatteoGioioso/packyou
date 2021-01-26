package fileCollector

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type fileCollector struct {
	entry string
	root string
	output string
}

func New(entry string, root string, output string) *fileCollector {
	return &fileCollector{entry: entry, root: root, output: output}
}

func (f *fileCollector) isCommonJs(line string) bool {
	if strings.Contains(line, "require(") {
		return true
	}

	return false
}

func (f *fileCollector) isES6Module(line string) bool {
	if strings.Contains(line, "import") {
		return true
	}

	return false
}

func (f *fileCollector) isNodeModule(importPath string) bool {
	if strings.HasPrefix(importPath, "./") || strings.HasPrefix(importPath, "../") {
		return false
	}

	return true
}

func (f *fileCollector) Collect() {
	entryPointAbsPath := filepath.Join(f.root, f.entry)
	f.moveFileToDest(f.entry)
	file, err := ioutil.ReadFile(entryPointAbsPath)
	if err != nil {
		fmt.Println(err)
	}

	lines := strings.Split(string(file), "\n")
	for _, line := range lines {
		if f.isCommonJs(line) {

		}

		if f.isES6Module(line) {
			f.parseES6Module(line)
		}
	}
}

func (f *fileCollector) parseES6Module(line string) {
	importPath := strings.Split(line, "from")[1]
	importPath = strings.ReplaceAll(importPath, " ", "")
	importPath = strings.ReplaceAll(importPath, ";", "")
	importPath = strings.ReplaceAll(importPath, "'", "")
	importPath = strings.ReplaceAll(importPath, "\"", "")
	importPath = fmt.Sprintf("%v.js", importPath)
	if f.isNodeModule(importPath) {
		return
	}

	f.moveFileToDest(importPath)
}

func (f *fileCollector) moveFileToDest(importPath string) error {
	filePath := filepath.Join(f.root, importPath)
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Opening: ", err)
		return err
	}

	outputPath, err := filepath.Abs(f.output)
	if err != nil {
		return err
	}

	outputPath = filepath.Join(outputPath, importPath)
	if err := os.MkdirAll(filepath.Dir(outputPath), 0770); err != nil {
		return err
	}

	if err := ioutil.WriteFile(outputPath, file, os.ModePerm); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
