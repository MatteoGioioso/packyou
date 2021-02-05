package compiler_test

import (
	"github.com/onsi/gomega"
	"packyou/pku/compiler"
	"packyou/pku/pathResolver"
	"testing"
)

func Test_compiler_TransformImport(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	type args struct {
		line string
		importPath string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "should transform esm to commonjs",
			args: args{
				line: "import file from './services/function-level-file-1'",
				importPath: "./services/function-level-file-1",
			},
			want: "const file = require('./services/function-level-file-1');",
		},
		{
			name: "should transform esm to commonjs with destructuring",
			args: args{
				line: "import { file1, file2 } from './services/function-level-file-1'",
				importPath: "./services/function-level-file-1",
			},
			want: "const { file1, file2 } = require('./services/function-level-file-1');",
		},
		{
			name: "should transform esm to commonjs with un-named imports",
			args: args{
				line: "import './services/function-level-file-1'",
				importPath: "./services/function-level-file-1",
			},
			want: "require('./services/function-level-file-1');",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := compiler.New(pathResolver.New("", "", ""))

			res := c.TransformImportToCommonJs(tt.args.line, tt.args.importPath)

			g.Expect(res).Should(gomega.Equal(tt.want))
		})
	}
}

func Test_compiler_TransformExport(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	type args struct {
		line string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "should transform default export",
			args: args{
				line: "export default myFunc",
			},
			want: "module.exports = myFunc",
		},
		{
			name: "should transform export with function",
			args: args{
				line: "export function myFunc(){}",
			},
			want: "module.exports.myFunc = function myFunc(){}",
		},
		{
			name: "should transform export arrow function",
			args: args{
				line: "export const myFunc = (myArg1, myArg2) => {}",
			},
			want: "module.exports.myFunc = const myFunc() => {}",
		},
		{
			name: "should transform export with class",
			args: args{
				line: "export class MyClass{}",
			},
			want: "module.exports.MyClass = class MyClass{}",
		},
		{
			name: "should transform export with object",
			args: args{
				line: "export { myFunc }",
			},
			want: "module.exports = { myFunc }",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := compiler.New(pathResolver.New("", "", ""))

			res := c.TransformExportToCommonJs(tt.args.line)

			g.Expect(res).Should(gomega.Equal(tt.want))
		})
	}
}

