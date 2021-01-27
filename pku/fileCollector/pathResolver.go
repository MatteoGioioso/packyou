package fileCollector

import (
	"fmt"
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

	return outputPath, err
}

func (r pathResolver) IsExternalReference(path string) bool {
	if strings.Contains(path, "../") {
		return true
	}

	return false
}

// ChangeMovedFileImportPath since we are moving all the external references to the root
// of the function we need to modify the import path
func (r pathResolver) ChangeMovedFileImportPath(line string, currentOriginPath string, importPath string) string {
	var newPath string
	aORb := regexp.MustCompile("\\.\\./") // Match ../

	// We calculate the difference between the two paths so we can identify how many
	// times the external reference is far from our entry folder
	base := strings.ReplaceAll(importPath, "../", "")
	pathAfter := strings.ReplaceAll(currentOriginPath, base, "")
	diff := len(strings.Split(r.getEntryAbs(), "/")) - len(strings.Split(pathAfter, "/"))

	// If bigger than one means that we are 2 times outside the entry folder
	if diff > 0 {
		// If we remain without "../" we need to add a relative import "./"
		// otherwise nodejs will look in node_modules
		replacedPath := strings.Replace(line, "../", "", diff+1)
		matches := aORb.FindAllString(replacedPath, -1)

		if len(matches) == 0 {
			newPath = strings.Replace(line, "../", "", diff)
			newPath = strings.Replace(newPath, "../", "./", 1)
		} else {
			newPath = replacedPath
		}
	} else {
		// If 0 then it means we are in the entry folder
		matches := aORb.FindAllStringIndex(line, -1)
		if len(matches) > 1 {
			newPath = strings.Replace(line, "../", "", 1)
		} else {
			newPath = strings.Replace(line, "../", "./", 1)
		}
	}

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
