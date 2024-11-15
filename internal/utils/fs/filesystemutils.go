package fs_utils

import (
	"os"
	"path"
	"path/filepath"

	"github.com/charmbracelet/log"
)

type FilePickerDataType int

const (
	FileOnlyDataType FilePickerDataType = iota
	DirOnlyDataType
	DirOrFileDataType
)

type FSDirectory struct {
	Path     string
	DataType FilePickerDataType
}

func (f FSDirectory) getAllImmediateChildren() []string {
	c, err := os.ReadDir(f.Path)
	if err != nil {
		log.Fatal(err)
	}
	var items []string
	for _, entry := range c {
		items = append(items, entry.Name())
	}
	return items
}

func (f FSDirectory) getImmediateChildDirs() []string {
	c, err := os.ReadDir(f.Path)
	if err != nil {
		log.Fatal(err)
	}
	var items []string
	for _, entry := range c {
		if entry.IsDir() {
			items = append(items, entry.Name())
		}
	}
	return items
}

func (f FSDirectory) getImmediateChildFiles() []string {
	c, err := os.ReadDir(f.Path)
	if err != nil {
		log.Fatal(err)
	}
	var items []string
	for _, entry := range c {
		if !entry.IsDir() {
			items = append(items, entry.Name())
		}
	}
	return items
}

func (f *FSDirectory) GetDataFromAbsolutePath(dirPath string) []string {
	f.Path = dirPath
	switch f.DataType {
	case FileOnlyDataType:
		return f.getImmediateChildFiles()
	case DirOnlyDataType:
		return f.getImmediateChildDirs()
	default:
		return f.getAllImmediateChildren()
	}
}

func (f *FSDirectory) GetNewData(dirPath string) []string {
	var dir string
	if f.Path != "" {
		dir = path.Join(f.Path, dirPath)
	} else {
		dir = dirPath
	}
	return f.GetDataFromAbsolutePath(dir)
}

func (f *FSDirectory) GetParentData() []string {
	dir := filepath.Dir(f.Path)
	return f.GetDataFromAbsolutePath(dir)
}
