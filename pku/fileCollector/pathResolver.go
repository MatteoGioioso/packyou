package fileCollector

import (
	"fmt"
	"packyou/pku/errorPkg"
	"path/filepath"
	"regexp"
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
		return "", errorPkg.New(err, "GetDestFileLocation")
	}

	entryFolderName := r.getEntryFolderName()
	fileName := filepath.Base(currentOriginFilePath)
	outputPath = filepath.Join(outputPath, entryFolderName, fileName)

	return outputPath, err
}

func (r pathResolver) IsExternalReference(path string) bool {
	if strings.Contains(path, "../") {
		return true
	}

	return false
}

// ChangeMovedFileImportPath since we are moving all the references to the root
// of the function we need to modify the import path
func (r pathResolver) ChangeMovedFileImportPath(line string, importPath string) string {
	var newPath string
	importPathDir := filepath.Dir(importPath)
	// Count how many "../" (go to parent we have)
	aORb := regexp.MustCompile("\\.\\./") // Match ../
	matches := aORb.FindAllString(importPath, -1)
	// Remove all the "../" and add "./"
	newPath = strings.Replace(line, "../", "", len(matches)-1)
	newPath = strings.Replace(newPath, "../", "./", 1)
	// Remove the sub directories by replacing with the import path directory
	newPath = strings.Replace(newPath, importPathDir+"/", "", 1)

	return newPath
}

func (r pathResolver) getEntryAbs() string {
	return filepath.Join(
		r.projectRoot,
		filepath.Dir(r.entryFilePath),
	)
}

func (r pathResolver) getAbsEntryFilePath() string {
	return filepath.Join(r.projectRoot, r.entryFilePath)
}

func (r pathResolver) getEntryFolderName() string {
	return filepath.Base(filepath.Dir(r.entryFilePath))
}
