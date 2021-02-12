package files

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Filer _
type Filer interface {
	FullPath() string
	Name() string
	SetChTime(timeObj time.Time) error
	EnsureParentDirExists() error
	Remove() error
	SetDirPath(path Pather)
	SetFileName(fileName string)
	Clone() Filer
	BaseName() string
	Extension() string
	NewWithSuffix(suffix string) Filer
	BuildPath() Pather
	IsExist() bool
	ReadAllContent() (string, error)
	Create() error
	Size() (int, error)
	ReadContent() (FileReader, error)
	WriteContent() (FileWriter, error)
	Move(newFullPath string) error
	Rename(newName string) error
	MarshalYAML() (interface{}, error)
}

// File _
type File struct {
	fileName string
	dirPath  string
}

// NewFile _
func NewFile(relativePath string) *File {
	fullPath, _ := filepath.Abs(relativePath)

	dirPath, fileName := filepath.Split(fullPath)

	return &File{
		fileName: fileName,
		dirPath:  dirPath,
	}
}

// FullPath _
func (f *File) FullPath() string {
	return filepath.Join(f.dirPath, f.fileName)
}

// Name _
func (f *File) Name() string {
	return f.fileName
}

// DirPath returns file's parent directory path
func (f *File) DirPath() string {
	return f.dirPath
}

// BuildPath returns Pather instance of file's parent directory
func (f *File) BuildPath() Pather {
	return NewPath(f.DirPath())
}

// IsExist _
func (f *File) IsExist() bool {
	info, err := os.Stat(f.FullPath())
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// SetDirPath _
func (f *File) SetDirPath(path Pather) {
	f.dirPath = path.FullPath()
}

// SetFileName _
func (f *File) SetFileName(fileName string) {
	f.fileName = fileName
}

// Size _
func (f *File) Size() (int, error) {
	info, err := os.Stat(f.FullPath())

	if err != nil {
		return 0, errors.Wrap(err, "Getting file size")
	}

	return int(info.Size()), nil
}

// Clone _
func (f *File) Clone() Filer {
	newFile := &File{}
	*newFile = *f
	return newFile
}

// BaseName returns file name without extension
func (f *File) BaseName() string {
	return strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))
}

// Extension returns file extension with dot from name. Example: ".mp4"
func (f *File) Extension() string {
	return filepath.Ext(f.Name())
}

// NewWithSuffix _
func (f *File) NewWithSuffix(suffix string) Filer {
	newFile := f.Clone()

	nameWithoutExt := strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))

	newFile.SetFileName(
		fmt.Sprintf(
			"%s%s%s",
			nameWithoutExt,
			suffix,
			filepath.Ext(f.Name()),
		),
	)

	return newFile
}

// SetChTime _
func (f *File) SetChTime(timeObj time.Time) error {
	return os.Chtimes(f.FullPath(), timeObj, timeObj)
}

// EnsureParentDirExists _
func (f *File) EnsureParentDirExists() error {
	path := NewPath(".")

	return path.Create()
}

// Create creates file
func (f *File) Create() error {
	parentDir := f.BuildPath()

	err := parentDir.Create()

	if err != nil {
		return errors.Wrap(err, "Creating parent directory")
	}

	_, err = os.Create(f.FullPath())

	return err
}

// FileWriter _
type FileWriter interface {
	io.Writer
	io.StringWriter
	io.Closer
}

// FileReader _
type FileReader interface {
	io.Reader
	io.Closer
}

// ReadContent _
func (f *File) ReadContent() (FileReader, error) {
	return os.Open(f.FullPath())
}

// WriteContent _
func (f *File) WriteContent() (FileWriter, error) {
	return os.OpenFile(f.FullPath(), os.O_APPEND|os.O_WRONLY, 0644)
}

// Remove _
func (f *File) Remove() error {
	return os.Remove(f.FullPath())
}

// ReadAllContent _
func (f *File) ReadAllContent() (string, error) {
	file, err := os.Open(f.FullPath())

	if err != nil {
		return "", err
	}

	defer file.Close()

	b, err := ioutil.ReadAll(file)

	return string(b), nil
}

// Move _
func (f *File) Move(newFullPath string) error {
	return os.Rename(f.FullPath(), newFullPath)
}

// Rename _
func (f *File) Rename(newName string) error {
	return os.Rename(
		f.FullPath(),
		path.Join(f.DirPath(), newName),
	)
}

// MarshalYAML is YAML Marshaller interface implementation
func (f *File) MarshalYAML() (interface{}, error) {
	return f.FullPath(), nil
}
