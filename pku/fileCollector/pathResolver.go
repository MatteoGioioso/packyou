package fileCollector

import (
	"fmt"
	"path/filepath"
	"strings"
)

type pathResolver struct {
	projectRoot   string
	entryFilePath string
	outputPath    string
}

func (r pathResolver) GetRawImportPathForES6Module(line string) string {
	return strings.Split(line, "from")[1]
}

func (r pathResolver) ExtractImportPathFromLine(line string) string {
	importPath := r.GetRawImportPathForES6Module(line)
	importPath = r.CleanRawImportPath(importPath)
	importPath = fmt.Sprintf("%v.js", importPath)
	return importPath
}

func (r pathResolver) CleanRawImportPath(rawImportPath string) string {
	rawImportPath = strings.ReplaceAll(rawImportPath, " ", "")
	rawImportPath = strings.ReplaceAll(rawImportPath, ";", "")
	rawImportPath = strings.ReplaceAll(rawImportPath, "'", "")
	rawImportPath = strings.ReplaceAll(rawImportPath, "\"", "")
	return rawImportPath
}

func (r pathResolver) GetOriginFileLocation(currentOriginFilePath, importPath string) string {
	currentOriginDir := filepath.Dir(currentOriginFilePath)
	originalPath := filepath.Join(currentOriginDir, importPath)
	return originalPath
}

func (r pathResolver) GetDestFileLocation(currentOriginFilePath string) (string, error) {
	outputPath, err := filepath.Abs(r.outputPath)
	if err != nil {
		return "", err
	}

	rootPath := filepath.Join(r.projectRoot, r.entryFilePath)
	if len(currentOriginFilePath) < len(rootPath) {
		outputPath = filepath.Join(
			outputPath,
			filepath.Dir(r.entryFilePath),
			filepath.Base(currentOriginFilePath),
		)
		outputPath = strings.ReplaceAll(outputPath, filepath.Dir("/"+r.entryFilePath), "")
	} else {
		outputPath = filepath.Join(outputPath, currentOriginFilePath)
		outputPath = strings.ReplaceAll(outputPath, r.projectRoot, "")
		outputPath = strings.ReplaceAll(outputPath, filepath.Dir("/"+r.entryFilePath), "")
	}

	fmt.Println(outputPath)
	return outputPath, err
}

func (r pathResolver) IsExternalReference(path string) bool {
	if strings.Contains(path, "../") {
		return true
	}

	return false
}
