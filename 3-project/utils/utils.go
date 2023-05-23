package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type FileScan struct {
	basePath string
}

type InvalidDir struct {
	path string
}

func (e *InvalidDir) Error() string {
	return fmt.Sprintf("path %q isn't a directory", e.path)
}

type Scan interface {
	Status(string) (time.Time, error)
	Build(string) (time.Time, error)
}

// NewFileScan returns a file scan given
// a base path. Returns an error if it couldn't
// validate the path or it doesn't point to a dir
func NewFileScan(path string) (*FileScan, error) {
	// Validates path
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, &InvalidDir{path: path}
	}

	return &FileScan{basePath: path}, nil
}

// join appends the base dir to path
func (fscan *FileScan) join(path string) string {
	return filepath.Join(fscan.basePath, path)
}

// NOTE: Testing purposes
func (fscan *FileScan) create(filename string) (info *os.File) {
	info, _ = os.Create(fscan.join(filename))
	return
}

// NOTE: Testing purposes
func (fscan *FileScan) remove(filename string) {
	os.Remove(fscan.join(filename))
}

// Status Checks if file given by path exists and returns its modification time if so, errors otherwise.
func (fscan *FileScan) Status(filename string) (time.Time, error) {
	fs, err := os.Stat(fscan.join(filename))
	if err != nil {
		return time.Time{}, err
	}
	return fs.ModTime(), nil
}

// Build Fake builds the object file and returns its modification time.
func (fscan *FileScan) Build(filename string) (time.Time, error) {
	filename = fscan.join(filename)

	f, err := os.Open(filename)
	var n int
	if err == nil { // File existed, read n
		scanner := bufio.NewScanner(f)
		scanner.Split(bufio.ScanWords)
		scanner.Scan()
		n, _ = strconv.Atoi(scanner.Text())
		n++
		f.Close()
	}

	f, err = os.Create(filename)
	defer f.Close()
	if err != nil {
		return time.Time{}, err
	}
	_, err = f.WriteString(strconv.Itoa(n) + " times built.\n")
	if err != nil {
		return time.Time{}, err
	}

	fs, _ := f.Stat()
	t := fs.ModTime()

	return t, nil

}
