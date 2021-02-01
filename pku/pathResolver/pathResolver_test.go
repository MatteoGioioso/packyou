package pathResolver

import (
	"github.com/onsi/gomega"
	"testing"
)

func Test_pathResolver_ChangeMovedFileImportPath(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "should resolve path with double", args: args{path: "../../index.js"}, want: "../index.js"},
		{name: "should resolve path with triple", args: args{path: "../../../index.js"}, want: "../../index.js"},
		{name: "should resolve path with only one", args: args{path: "../index.js"}, want: "./index.js"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := pathResolver{}

			newPath, _ := f.ChangeFileImportPathToNewLocation(tt.args.path, "", "")

			g.Expect(newPath).Should(gomega.Equal(tt.want))
		})
	}
}

func Test_pathResolver_GetDestFileLocation(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	type fields struct {
		projectRoot   string
		entryFilePath string
		outputPath    string
	}
	type args struct {
		currentOriginFilePath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "should get the right output path",
			fields: fields{
				projectRoot:   "/Users/madeo/Development/Projects/gfg-gpe-daas-catalogue-api",
				entryFilePath: "services/catalogueSearch/handler.js",
				outputPath:    "../dist",
			},
			args: args{
				currentOriginFilePath: "/Users/madeo/Development/Projects/gfg-gpe-daas-catalogue-api/services/catalogueSearch/repositories/DaasCatalogueRepository.js",
			},
			want:    "/Users/madeo/Development/Tools/packyou/pku/dist/repositories/DaasCatalogueRepository.js",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := pathResolver{
				projectRoot:   tt.fields.projectRoot,
				entryFilePath: tt.fields.entryFilePath,
				outputPath:    tt.fields.outputPath,
			}
			got, err := r.GetDestFileLocation(tt.args.currentOriginFilePath)
			g.Expect(got).Should(gomega.Equal(tt.want))

			if tt.wantErr {
				g.Expect(err).Should(gomega.HaveOccurred())
			} else {
				g.Expect(err).ShouldNot(gomega.HaveOccurred())
			}
		})
	}
}

func Test_pathResolver_ChangeMovedFileImportPath1(t *testing.T) {
	type fields struct {
		projectRoot   string
		entryFilePath string
		outputPath    string
	}
	type args struct {
		line       string
		importPath string
	}
	tests := []struct {
		name              string
		fields            fields
		args              args
		wantNewImportPath string
		wantNewLine       string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := pathResolver{
				projectRoot:   tt.fields.projectRoot,
				entryFilePath: tt.fields.entryFilePath,
				outputPath:    tt.fields.outputPath,
			}
			gotNewImportPath, gotNewLine := r.ChangeFileImportPathToNewLocation(tt.args.line, tt.args.importPath, "")
			if gotNewImportPath != tt.wantNewImportPath {
				t.Errorf("ChangeFileImportPathToNewLocation() gotNewImportPath = %v, want %v", gotNewImportPath, tt.wantNewImportPath)
			}
			if gotNewLine != tt.wantNewLine {
				t.Errorf("ChangeFileImportPathToNewLocation() gotNewLine = %v, want %v", gotNewLine, tt.wantNewLine)
			}
		})
	}
}
