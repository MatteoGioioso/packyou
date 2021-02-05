package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	pathRes  pathResolver.PathResolver
	comp     compiler.Compiler
)

// Naming conventions:
// FileImportPath: is the actual import path appearing in the import declaration: "./path/to/file"
// Line: the line of code from the file

func initializeCommand(cmd *cobra.Command) {
	entry := cmd.Flag("entry").Value.String()
	projectRoot := cmd.Flag("project-root").Value.String()
	output := cmd.Flag("output").Value.String()
	compileToCommonjs := cmd.Flag("compile-commonjs").Value.String()
	addExtension := cmd.Flag("add-extension").Value.String()

	viper.Set("entry", entry)
	viper.Set("output", output)
	viper.Set("projectRoot", projectRoot)
	viper.Set("compileCommonjs", compileToCommonjs)
	viper.Set("addExtension", addExtension)

	fileRepo = fileRepository.New()
	pathRes = pathResolver.New(projectRoot, entry, output)
	comp = compiler.New(pathRes)

	entryFileLocation := filepath.Join(projectRoot, entry)
	if err := collect(entryFileLocation); err != nil {
		log.Fatal(err)
	}
}

func collect(originFileLocation string) error {
	file, err := fileRepo.GetFile(originFileLocation)
	if err != nil {
		return errorPkg.New(err, "collect")
	}

	lines := strings.Split(string(file), "\n")
	for _, line := range lines {
		if pathRes.IsES6Module(line) {
			if err := parseES6Module(line, originFileLocation); err != nil {
				return errorPkg.New(err, "collect")
			}
		}

		if pathRes.IsCommonJs(line) {
			// Implement commonJs
		}
	}

	newFile := fileRepo.RewriteToFile(string(file))
	destFileLocation, err := pathRes.GetDestFileLocation(originFileLocation)
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
	} else {
		importPath = pathRes.ExtractImportPathFromLine(line)
	}

	if pathRes.IsNodeModule(importPath) {
		// TODO: for now just copy node_modules into the dest folder
		// no-op
		return nil
	}

	originFileLocation := pathRes.GetOriginFileLocation(currentOriginFilePath, importPath)
	newImportPath := pathRes.ChangeFileImportPathToNewLocation(importPath, currentOriginFilePath)
	var newLine string
	if viper.Get("compileCommonjs") == "true" {
		newLine = comp.TransformImportToCommonJs(line, newImportPath)
	} else {
		newLine = pathRes.ChangeLineWithNewImportPath(line, newImportPath)
	}

	fileRepo.AddRewrite(line, newLine)

	return collect(originFileLocation)
}
