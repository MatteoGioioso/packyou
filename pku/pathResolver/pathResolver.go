package pathResolver

import (
	"fmt"
	"math/rand"
	"packyou/pku/errorPkg"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type pathResolver struct {
	projectRoot   string
	entryFilePath string
	outputPath    string
	fileMap       map[string]string
	lookupFileMap map[string]string
}

func New(projectRoot string, entryFilePath string, outputPath string) *pathResolver {
	return &pathResolver{
		projectRoot: projectRoot,
		entryFilePath: entryFilePath,
		outputPath: outputPath,
		// This is for checking if there are filename collisions
		fileMap: make(map[string]string, 0),
		// This is for retrieving the new file name
		lookupFileMap: make(map[string]string, 0),
	}
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
	fileName := r.lookupFileMap[currentOriginFilePath]
	if fileName == "" {
		fileName = filepath.Base(currentOriginFilePath)
	}
	outputPath = filepath.Join(outputPath, entryFolderName, fileName)

	return outputPath, err
}

// getNewFilenameInCaseOfCollision if it detects a naming collision it will generate a new file name
// and store the new file name into the lookupFileMap for retrieval and fileMap for checking
func (r *pathResolver) getNewFilenameInCaseOfCollision(currentOriginFilePath string) string {
	fileName := filepath.Base(currentOriginFilePath)

	if originFilePath, ok := r.fileMap[fileName]; ok {
		// If the filename is the same but the origin path is different
		// means that there is going to be a collision
		if originFilePath != currentOriginFilePath {
			fileName = r.generateFileNameUniqueId(fileName)
			r.fileMap[fileName] = currentOriginFilePath
		}
	} else {
		r.fileMap[fileName] = currentOriginFilePath
	}

	r.lookupFileMap[currentOriginFilePath] = fileName

	return fileName
}

func (r pathResolver) generateFileNameUniqueId(filename string) string {
	rand.Seed(time.Now().UnixNano())
	randomString := strconv.Itoa(rand.Intn(10))
	ext := filepath.Ext(filename)
	name := strings.Replace(filename, ext, "", 1)
	return name +randomString+ext
}

func (r pathResolver) getNewImportPathInCaseOfCollision(importPath, currentOriginFileLocation string) string {
	newFileName := r.getNewFilenameInCaseOfCollision(r.GetOriginFileLocation(currentOriginFileLocation, importPath))
	oldFileName := filepath.Base(importPath)
	return strings.Replace(importPath, oldFileName, newFileName, 1)
}

// ChangeMovedFileImportPath since we are moving all the references to the root
// of the function we need to modify the import path
func (r pathResolver) ChangeFileImportPathToNewLocation(line string, importPath string, currentOriginFileLocation string) (newLine string, newImportPath string) {
	importPath = r.getNewImportPathInCaseOfCollision(importPath, currentOriginFileLocation)
	importPathDir := filepath.Dir(importPath)
	// Temporary replace the importPath in the line with a placeholder
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
	ChangeFileImportPathToNewLocation(line string, importPath string, originFileLocation string) (newLine string, newImportPath string)
	GetEntryAbs() string
	GetAbsEntryFilePath() string
	GetEntryFolderName() string
	IsCommonJs(line string) bool
	IsES6Module(line string) bool
	IsNodeModule(importPath string) bool
	IsUnnamedES6Import(line string) bool
}
