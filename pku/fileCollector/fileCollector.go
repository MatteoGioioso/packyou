package fileCollector

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type fileCollector struct {
	entry    string
	root     string
	output   string
	rewrites map[string]string
}

func New(entry string, root string, output string) *fileCollector {
	return &fileCollector{
		entry:    entry,
		root:     root,
		output:   output,
		rewrites: make(map[string]string, 0),
	}
}

func (f *fileCollector) Collect() {
	f.collect(f.entry)
}

func (f fileCollector) collect(importPath string) {
	entryPointAbsPath := filepath.Join(f.root, importPath)
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

	toFile := f.rewriteExternalReferencesToFile(string(file))
	path, err := f.getOutputPath(importPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := f.saveFile(path, []byte(toFile)); err != nil {
		fmt.Println(err)
		return
	}
}

func (f *fileCollector) parseES6Module(line string) {
	var importPath string
	if !strings.Contains(line, "from") {
		// TODO: handle un-named imports "import babel/regenerator"
		// no-op
		return
	}

	importPath = strings.Split(line, "from")[1]
	importPath = strings.ReplaceAll(importPath, " ", "")
	importPath = strings.ReplaceAll(importPath, ";", "")
	importPath = strings.ReplaceAll(importPath, "'", "")
	importPath = strings.ReplaceAll(importPath, "\"", "")
	importPath = fmt.Sprintf("%v.js", importPath)

	if f.isExternalReference(importPath) {
		f.addRewrites(line)
	}

	if f.isNodeModule(importPath) {
		// TODO: for now just copy node_modules into the dest folder
		// no-op
		return
	}

	if err := f.moveFileToDest(importPath); err != nil {
		fmt.Println(err)
		return
	}

	f.collect(importPath)
}

func (f *fileCollector) moveFileToDest(importPath string) error {
	file, err := f.getFile(importPath)
	if err != nil {
		return err
	}

	outputPath, err := f.getOutputPath(importPath)
	if err != nil {
		return err
	}

	if err := f.saveFile(outputPath, file); err != nil {
		return err
	}

	return nil
}

func (f fileCollector) saveFile(outputPath string, file []byte) error {
	if err := os.MkdirAll(filepath.Dir(outputPath), 0770); err != nil {
		return err
	}

	if err := ioutil.WriteFile(outputPath, file, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func (f fileCollector) getFile(importPath string) ([]byte, error) {
	originalPath := filepath.Join(f.root, importPath)
	file, err := ioutil.ReadFile(originalPath)
	if err != nil {
		return nil, err
	}

	return file, err
}

func (f fileCollector) getOutputPath(importPath string) (string, error) {
	outputPath, err := filepath.Abs(f.output)
	if err != nil {
		return "", err
	}
	if f.isExternalReference(importPath) {
		importPath = f.rewritePath(importPath)
	}

	outputPath = filepath.Join(outputPath, importPath)
	return outputPath, err
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

// does file exists in the parent directories?
func (f fileCollector) isExternalReference(path string) bool {
	if strings.Contains(path, "../") {
		return true
	}

	return false
}

// Change all the imports of external references into the new path
func (f fileCollector) rewriteExternalReferencesToFile(file string) string {
	for oldPath, newPath := range f.rewrites {
		file = strings.ReplaceAll(file, oldPath, newPath)
	}

	return file
}

func (f fileCollector) rewritePath(path string) string {
	aORb := regexp.MustCompile("../")
	matches := aORb.FindAllStringIndex(path, -1)

	var newPath string
	if len(matches) > 1 {
		newPath = strings.Replace(path, "../", "", 1)
	} else {
		newPath = strings.ReplaceAll(path, "../", "./")
	}

	return newPath
}

// This function store reference that are outside ../
func (f *fileCollector) addRewrites(path string) {
	newPath := f.rewritePath(path)
	f.rewrites[path] = newPath
}
