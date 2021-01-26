package fileCollector

import (
	"github.com/onsi/gomega"
	"testing"
)

func Test_fileCollector_rewritePath(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	type args struct {
		path string
	}
	tests := []struct {
		name   string
		args   args
		want   string
	}{
		{name: "should resolve path with double", args: args{path: "../../index.js"}, want: "../index.js"},
		{name: "should resolve path with triple", args: args{path: "../../../index.js"}, want: "../../index.js"},
		{name: "should resolve path with only one", args: args{path: "../index.js"}, want: "./index.js"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fileCollector{}

			newPath := f.rewritePath(tt.args.path)

			g.Expect(newPath).Should(gomega.Equal(tt.want))
		})
	}
}
