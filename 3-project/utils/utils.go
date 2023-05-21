package utils

import (
	"bufio"
	"os"
	"strconv"
	"time"
)

// GetModTime Returns the target modification time, assuming the target already exists
func GetModTime(target string) time.Time {
	fs, _ := os.Stat(target)
	return fs.ModTime()
}

// Status Checks if file given by path exists and returns its modification time if so, errors otherwise.
func Status(path string) (time.Time, error) {
	fs, err := os.Stat(path)
	if err != nil {
		return time.Time{}, err
	}
	return fs.ModTime(), nil
}

// Build Fake builds the object file and returns its modification time.
func Build(object string) (time.Time, error) {
	f, err := os.Open(object)
	var n int
	if err == nil { // File existed, read n
		scanner := bufio.NewScanner(f)
		scanner.Split(bufio.ScanWords)
		scanner.Scan()
		n, _ = strconv.Atoi(scanner.Text())
		n++
		f.Close()
	}

	f, err = os.Create(object)
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
