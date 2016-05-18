package core

import (
	"io/ioutil"
	"os"

	"github.com/inkyblackness/res/serial"
)

type fileAccess struct {
	readDir    func(dirname string) ([]os.FileInfo, error)
	readFile   func(filename string) ([]byte, error)
	createFile func(filename string) (serial.SeekingWriteCloser, error)
}

var realFileAccess = fileAccess{
	readDir:    ioutil.ReadDir,
	readFile:   ioutil.ReadFile,
	createFile: func(filename string) (serial.SeekingWriteCloser, error) { return os.Create(filename) }}
