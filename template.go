package ansible

import (
	"io/ioutil"
)

type Template struct {
	Path string
	Data []byte
}

type Templates []*Template

func NewTemplate(path string, data []byte) *Template {

	return &Template{
		Path: path,
		Data: data,
	}
}

func ReadTemplate(sourcePath, targetPath string) (*Template, error) {

	data, err := ioutil.ReadFile(sourcePath)
	if err != nil {
		return nil, err
	}

	return NewTemplate(targetPath, data), nil
}
