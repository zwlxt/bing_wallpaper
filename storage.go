package main

import (
	"io/ioutil"
	"os"
)

type StorageManager interface {
	Load(name string) []byte
	Save(img []byte, name string)
}

type FileSystemStorage struct {
	Dir string
}

func (fs *FileSystemStorage) Save(img []byte, name string) {
	if _, err := os.Stat(fs.Dir); os.IsNotExist(err) {
		os.Mkdir(fs.Dir, 0755)
	}
	path := fs.Dir + "/" + name
	err := ioutil.WriteFile(path, img, 0755)
	if err != nil {
		panic(err)
	}
}

func (fs *FileSystemStorage) Load(name string) []byte {
	f, err := ioutil.ReadFile(fs.Dir + "/" + name)
	if err != nil {
		panic(err)
	}
	return f
}
