package core

import (
	"os"
	"time"
)

type testingFileInfo struct {
	name string
}

func testFile(name string) *testingFileInfo {
	info := &testingFileInfo{name: name}

	return info
}

func testFiles(names ...string) []os.FileInfo {
	result := make([]os.FileInfo, len(names))
	for i, name := range names {
		result[i] = testFile(name)
	}

	return result
}

func (info *testingFileInfo) Name() string {
	return info.name
}

func (info *testingFileInfo) Size() int64 {
	return 0
}

func (info *testingFileInfo) Mode() os.FileMode {
	return os.FileMode(0)
}

func (info *testingFileInfo) ModTime() time.Time {
	return time.Now()
}
func (info *testingFileInfo) IsDir() bool {
	return false
}

func (info *testingFileInfo) Sys() interface{} {
	return nil
}
