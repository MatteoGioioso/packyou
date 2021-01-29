package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"packyou/pku/compiler"
	"packyou/pku/errorPkg"
	"packyou/pku/fileRepository"
	"packyou/pku/pathResolver"
	"path/filepath"
	"strings"
)

var (
	fileRepo fileRepository.FileRepository
	pathRes pathResolver.PathResolver
	comp compiler.Compiler
)

func initializeCommand(cmd *cobra.Command, getConfig func(key string) interface{}) {
	entry := cmd.Flag("entry").Value.String()
	projectRoot := cmd.Flag("projectRoot").Value.String()
	output := cmd.Flag("output").Value.String()
	fileRepo = fileRepository.New()
	pathRes = pathResolver.New(projectRoot, entry, output)
	comp = compiler.New(pathRes)

	entryFileLocation := filepath.Join(projectRoot, entry)
	if err := collect(entryFileLocation); err != nil {
		log.Fatal(err)
	}
}


func collect(originFilePath string) error {
	file, err := fileRepo.GetFile(originFilePath)
	if err != nil {
		return errorPkg.New(err, "collect")
	}

	lines := strings.Split(string(file), "\n")
	for _, line := range lines {
		if pathRes.IsES6Module(line) {
			if err := parseES6Module(line, originFilePath); err != nil {
				return errorPkg.New(err, "collect")
			}
		}

		if pathRes.IsCommonJs(line) {
			// Implement commonJs
		}
	}

	newFile := fileRepo.RewriteToFile(string(file))
	destFileLocation, err := pathRes.GetDestFileLocation(originFilePath)
	if err != nil {
		return errorPkg.New(err, "collect")
	}

	if err := fileRepo.SaveFile(destFileLocation, []byte(newFile)); err != nil {
		return errorPkg.New(err, "collect")
	}

	return nil
}

func parseES6Module(line, currentOriginFilePath string) error {
	var importPath string
	if pathRes.IsUnnamedES6Import(line) {
		// TODO: handle un-named imports "import babel/regenerator"
		// no-op
		return nil
	}

	importPath = pathRes.ExtractImportPathFromLine(line)
	if pathRes.IsNodeModule(importPath) {
		// TODO: for now just copy node_modules into the dest folder
		// no-op
		return nil
	}

	originFileLocation := pathRes.GetOriginFileLocation(currentOriginFilePath, importPath)
	destFileLocation, err := pathRes.GetDestFileLocation(originFileLocation)
	if err != nil {
		return errorPkg.New(err, "parseES6Module")
	}

	if err := fileRepo.MoveFileToDest(originFileLocation, destFileLocation); err != nil {
		return errorPkg.New(err, "parseES6Module")
	}

	newLine, newImportPath := pathRes.ChangeMovedFileImportPath(line, importPath)
	newLine = comp.TransformImport(newLine, newImportPath)
	fileRepo.AddRewrite(line, newLine)

	return collect(originFileLocation)
}
