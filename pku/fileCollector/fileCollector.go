package fileCollector

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type fileCollector struct {
	entry    string
	root     string
	output   string
	rewrites map[string]string
	pathResolver pathResolver
}

func New(entry string, root string, output string) *fileCollector {
	return &fileCollector{
		entry:    entry,
		root:     root,
		output:   output,
		rewrites: make(map[string]string, 0),
		pathResolver: pathResolver{
			projectRoot:   root,
			entryFilePath: entry,
			outputPath: output,
		},
	}
}

func (f *fileCollector) Collect() {
	originFilePath := filepath.Join(f.root, f.entry)
	f.collect(originFilePath)
}

func (f fileCollector) collect(originFilePath string) {
	file, err := f.getFile(originFilePath)
	if err != nil {
		fmt.Println(errors.Wrap(err, "collectRoot"))
	}

	lines := strings.Split(string(file), "\n")
	for _, line := range lines {
		if f.isES6Module(line) {
			f.parseES6Module(line, originFilePath)
		}

		if f.isCommonJs(line) {
			// Implement commonJs
		}
	}

	toFile := f.rewriteExternalReferencesToFile(string(file))
	path, err := f.pathResolver.GetDestFileLocation(originFilePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := f.saveFile(path, []byte(toFile)); err != nil {
		fmt.Println(err)
		return
	}
}

func (f *fileCollector) parseES6Module(line, currentOriginFilePath string) {
	var importPath string
	if !strings.Contains(line, "from") {
		// TODO: handle un-named imports "import babel/regenerator"
		// no-op
		return
	}

	importPath = f.pathResolver.ExtractImportPathFromLine(line)
	if f.isNodeModule(importPath) {
		// TODO: for now just copy node_modules into the dest folder
		// no-op
		return
	}

	originFileLocation := f.pathResolver.GetOriginFileLocation(currentOriginFilePath, importPath)
	if err := f.moveFileToDest(originFileLocation); err != nil {
		fmt.Println(err)
		return
	}

	if f.pathResolver.IsExternalReference(importPath) {
		f.addRewrites(line, originFileLocation, importPath)
	}

	f.collect(originFileLocation)
}

func (f *fileCollector) moveFileToDest(originFileLocation string) error {
	file, err := f.getFile(originFileLocation)
	if err != nil {
		return err
	}

	outputPath, err := f.pathResolver.GetDestFileLocation(originFileLocation)
	if err != nil {
		return err
	}

	if err := f.saveFile(outputPath, file); err != nil {
		return errors.Wrap(err, "saveFile\n")
	}

	return nil
}

func (f fileCollector) saveFile(outputPath string, file []byte) error {
	if err := os.MkdirAll(filepath.Dir(outputPath), 0770); err != nil {
		return errors.Wrap(err, "MkdirAll\n")
	}

	if err := ioutil.WriteFile(outputPath, file, os.ModePerm); err != nil {
		return errors.Wrap(err, "writeFile\n")
	}

	return nil
}

func (f fileCollector) getFile(importPath string) ([]byte, error) {
	file, err := ioutil.ReadFile(importPath)
	if err != nil {
		return nil, errors.Wrap(err, "getFile\n")
	}

	return file, err
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

// Change all the imports of external references into the new path
func (f fileCollector) rewriteExternalReferencesToFile(file string) string {
	for oldPath, newPath := range f.rewrites {
		file = strings.ReplaceAll(file, oldPath, newPath)
	}

	return file
}

// This function store reference that are outside ../
func (f *fileCollector) addRewrites(path string, currentOriginPath string, importPath string) {
	newPath := f.pathResolver.ChangeMovedFileImportPath(path, currentOriginPath, importPath)
	f.rewrites[path] = newPath
}
