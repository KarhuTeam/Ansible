package ansible

import (
	"io/ioutil"
)

type File struct {
	Path string
	Data []byte
}

type Files []*File

func NewFile(path string, data []byte) *File {

	return &File{
		Path: path,
		Data: data,
	}
}

func ReadFile(sourcePath, targetPath string) (*File, error) {

	data, err := ioutil.ReadFile(sourcePath)
	if err != nil {
		return nil, err
	}

	return NewFile(targetPath, data), nil
}
