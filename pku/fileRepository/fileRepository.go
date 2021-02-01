package fileRepository

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"packyou/pku/errorPkg"
	"path/filepath"
	"strings"
)

type fileRepository struct {
	rewrites map[string]string
	fileMap map[string]string
}

func New() *fileRepository {
	return &fileRepository{
		rewrites: make(map[string]string, 0),
	}
}

func (f fileRepository) SaveFile(outputPath string, file []byte) error {
	if err := os.MkdirAll(filepath.Dir(outputPath), 0770); err != nil {
		return errors.Wrap(err, "MkdirAll\n")
	}

	if err := ioutil.WriteFile(outputPath, file, os.ModePerm); err != nil {
		return errors.Wrap(err, "writeFile\n")
	}

	return nil
}

func (f *fileRepository) MoveFileToDest(originFileLocation, destFileLocation string) error {
	file, err := f.GetFile(originFileLocation)
	if err != nil {
		return err
	}

	if err := f.SaveFile(destFileLocation, file); err != nil {
		return errors.Wrap(err, "saveFile\n")
	}

	return nil
}

func (f fileRepository) GetFile(importPath string) ([]byte, error) {
	file, err := ioutil.ReadFile(importPath)
	if err != nil {
		return nil, errorPkg.New(err, "getFile")
	}

	return file, err
}

func (f fileRepository) RewriteToFile(file string) string {
	for oldPath, newPath := range f.rewrites {
		file = strings.ReplaceAll(file, oldPath, newPath)
	}

	return file
}

func (f *fileRepository) AddRewrite(line string, newLine string) {
	f.rewrites[line] = newLine
}

type FileRepository interface {
	SaveFile(outputPath string, file []byte) error
	MoveFileToDest(originFileLocation, destFileLocation string) error
	GetFile(importPath string) ([]byte, error)
	RewriteToFile(file string) string
	AddRewrite(line string, newLine string)
}
