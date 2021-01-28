package fileCollector

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"packyou/pku/compiler"
	"packyou/pku/errorPkg"
	"packyou/pku/pathResolver"
	"path/filepath"
	"strings"
)

type fileCollector struct {
	entry        string
	root         string
	output       string
	rewrites     map[string]string
	pathResolver pathResolver.PathResolver
	compiler     compiler.Compiler
}

func New(entry string, root string, output string) *fileCollector {
	pr := pathResolver.New(root, entry, output)
	return &fileCollector{
		entry:        entry,
		root:         root,
		output:       output,
		rewrites:     make(map[string]string, 0),
		pathResolver: pr,
		compiler: compiler.New(pr),
	}
}

func (f *fileCollector) Collect() {
	originFilePath := filepath.Join(f.root, f.entry)
	f.collect(originFilePath)
}

func (f fileCollector) collect(originFilePath string) {
	file, err := f.getFile(originFilePath)
	if err != nil {
		fmt.Println(errorPkg.New(err, "collect"))
	}

	lines := strings.Split(string(file), "\n")
	for _, line := range lines {
		if f.pathResolver.IsES6Module(line) {
			f.parseES6Module(line, originFilePath)
		}

		if f.pathResolver.IsCommonJs(line) {
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
	if f.pathResolver.IsUnnamedES6Import(line) {
		// TODO: handle un-named imports "import babel/regenerator"
		// no-op
		return
	}

	importPath = f.pathResolver.ExtractImportPathFromLine(line)
	if f.pathResolver.IsNodeModule(importPath) {
		// TODO: for now just copy node_modules into the dest folder
		// no-op
		return
	}

	originFileLocation := f.pathResolver.GetOriginFileLocation(currentOriginFilePath, importPath)
	if err := f.moveFileToDest(originFileLocation); err != nil {
		fmt.Println(err)
		return
	}

	f.addES6Rewrites(line, importPath)

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
		return nil, errorPkg.New(err, "getFile")
	}

	return file, err
}

// Change all the imports of external references into the new path
func (f fileCollector) rewriteExternalReferencesToFile(file string) string {
	for oldPath, newPath := range f.rewrites {
		file = strings.ReplaceAll(file, oldPath, newPath)
	}

	return file
}

// This function store reference that are outside ../
func (f *fileCollector) addES6Rewrites(line string, importPath string) {
	newLine, newImportPath := f.pathResolver.ChangeMovedFileImportPath(line, importPath)
	newLine = f.compiler.TransformImport(newLine, newImportPath)
	f.rewrites[line] = newLine
}
