package pathResolver

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

func New(projectRoot string, entryFilePath string, outputPath string) *pathResolver {
	return &pathResolver{projectRoot: projectRoot, entryFilePath: entryFilePath, outputPath: outputPath}
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

	entryFolderName := r.GetEntryFolderName()
	fileName := filepath.Base(currentOriginFilePath)
	outputPath = filepath.Join(outputPath, entryFolderName, fileName)

	return outputPath, err
}

// ChangeMovedFileImportPath since we are moving all the references to the root
// of the function we need to modify the import path
func (r pathResolver) ChangeMovedFileImportPath(line string, importPath string) (newImportPath, newLine string) {
	importPathDir := filepath.Dir(importPath)
	newLine = strings.ReplaceAll(line, importPath, "#__#")

	// Count how many "../" (go to parent we have)
	aORb := regexp.MustCompile("\\.\\./") // Match ../
	matches := aORb.FindAllString(importPath, -1)
	// Remove all the "../" and add "./"
	newImportPath = strings.Replace(importPath, "../", "", len(matches)-1)
	newImportPath = strings.Replace(newImportPath, "../", "./", 1)
	// Remove the sub directories by replacing with the import path directory
	newImportPath = strings.Replace(newImportPath, importPathDir+"/", "", 1)

	newLine = strings.ReplaceAll(line, "#__#", newImportPath)

	return newLine, newImportPath
}

func (r pathResolver) GetEntryAbs() string {
	return filepath.Join(
		r.projectRoot,
		filepath.Dir(r.entryFilePath),
	)
}

func (r pathResolver) GetAbsEntryFilePath() string {
	return filepath.Join(r.projectRoot, r.entryFilePath)
}

func (r pathResolver) GetEntryFolderName() string {
	return filepath.Base(filepath.Dir(r.entryFilePath))
}

func (r pathResolver) IsCommonJs(line string) bool {
	if strings.Contains(line, "require(") {
		return true
	}

	return false
}

func (r pathResolver) IsES6Module(line string) bool {
	if strings.Contains(line, "import") {
		return true
	}

	return false
}

func (r pathResolver) IsNodeModule(importPath string) bool {
	if strings.HasPrefix(importPath, "./") || strings.HasPrefix(importPath, "../") {
		return false
	}

	return true
}

func (r pathResolver) IsUnnamedES6Import(line string) bool {
	return !strings.Contains(line, "from")
}


type PathResolver interface {
	GetRawImportPathForES6Module(line string) string
	ExtractImportPathFromLine(line string) string
	CleanRawImportPath(rawImportPath string) string
	GetOriginFileLocation(currentOriginFilePath, importPath string) string
	GetDestFileLocation(currentOriginFilePath string) (string, error)
	ChangeMovedFileImportPath(line string, importPath string) (newLine string, newImportPath string)
	GetEntryAbs() string
	GetAbsEntryFilePath() string
	GetEntryFolderName() string
	IsCommonJs(line string) bool
	IsES6Module(line string) bool
	IsNodeModule(importPath string) bool
	IsUnnamedES6Import(line string) bool
}
