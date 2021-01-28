package compiler

import (
	"fmt"
	"packyou/pku/pathResolver"
	"strings"
)

type compiler struct {
	pathResolver pathResolver.PathResolver
}

func New(pathResolver pathResolver.PathResolver) *compiler {
	return &compiler{
		pathResolver: pathResolver,
	}
}

func (c compiler) TransformImport(line string, importPath string) string {
	if c.pathResolver.IsUnnamedES6Import(line) {
		return fmt.Sprintf("require('%v');", importPath)
	} else {
		varDeclaration := strings.Split(line, "from")[0]
		varDeclaration = strings.Replace(varDeclaration, "import", "const", 1)
		return fmt.Sprintf("%v= require('%v');", varDeclaration, importPath)
	}
}

func (c compiler) TransformExport(line string) string {
	panic("not implemented")
}

type Compiler interface {
	TransformImport(line string, importPath string) string
	TransformExport(line string) string
}
