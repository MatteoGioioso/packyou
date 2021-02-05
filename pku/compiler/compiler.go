package compiler

import (
	"fmt"
	"github.com/robertkrimen/otto/parser"
	"github.com/robertkrimen/otto/token"
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

func (c compiler) TransformImportToCommonJs(line string, importPath string) string {
	if c.pathResolver.IsUnnamedES6Import(line) {
		return fmt.Sprintf("require('%v');", importPath)
	} else {
		varDeclaration := strings.Split(line, "from")[0]
		varDeclaration = strings.Replace(varDeclaration, "import", "const", 1)
		return fmt.Sprintf("%v= require('%v');", varDeclaration, importPath)
	}
}

func (c compiler) TransformExportToCommonJs(line string) string {
	return c.prettyPrinter(line)
}

func (c compiler) prettyPrinter(line string) string {
	const variableTemplate = "module.exports.#__name__# = #__name__# = #__expression__#"
	const declarationTemplate = "module.exports.#__name__# = #__keyword__# #__name__# #__expression__#"
	const template = "module.exports = { #__name__# }"
	if strings.Contains(line, "export default") {
		out := fmt.Sprintf("module.exports =")
		return strings.Replace(line, "export default", out, 1)
	}

	return ""
}

func (c compiler) parser(line string) parsedExport {
	p := parsedExport{
		identifier: "module.exports",
	}
	psr := parser.NewParser("", line)
	for {
		scan, literal, idx := psr.Scan()
		fmt.Println(scan.String(), ": ", literal, " ", idx)

		switch scan.String() {
		case token.IDENTIFIER.String():
			p.variableName = literal
		case token.EOF.String():
			return p
		default:

		}
	}
}

type parsedExport struct {
	identifier   string
	expression   string
	keyword      string
	variableName string
}

type Compiler interface {
	TransformImportToCommonJs(line string, importPath string) string
	TransformExportToCommonJs(line string) string
}
