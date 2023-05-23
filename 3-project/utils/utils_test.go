package utils

import (
	"bufio"
	"os"
	"strconv"
	"testing"
)

var fileScan *FileScan

func init() {
	path := os.Getenv("UTILS_TEST_PATH")
	if path == "" {
		path = "./"
	}
	var err error
	if fileScan, err = NewFileScan(path); err != nil {
		panic(err)
	}
}

func TestFreshBuild(t *testing.T) {
	s := "foo"
	fileScan.remove(s)
	_, err := fileScan.Build(s)
	if err != nil {
		t.Error("Build failed,")
		return
	}
	f, err := os.Open(s)
	if err != nil {
		t.Error("Something went wrong opening the file", err)
		return
	}
	scan := bufio.NewScanner(f)
	scan.Split(bufio.ScanWords)
	scan.Scan()
	n, err := strconv.Atoi(scan.Text())
	if err != nil {
		t.Error("Something broke converting write number.")
	}
	if n != 0 {
		t.Error("File was not fresh.")
	}

}

func TestBuildInc(t *testing.T) {
	s := "foo"
	f := fileScan.create(s)
	f.WriteString(strconv.Itoa(10) + " times built.\n")
	f.Close()
	fileScan.Build(s)
	f, _ = os.Open(s)
	scan := bufio.NewScanner(f)
	scan.Split(bufio.ScanWords)
	scan.Scan()
	n, _ := strconv.Atoi(scan.Text())
	f.Close()
	os.Remove(s)
	if n != 11 {
		t.Error("File was not incremented properly.")
	}

}

func TestStatusNew(t *testing.T) {
	s := "bar"
	fileScan.remove(s)
	_, err := fileScan.Status(s)
	if err == nil {
		t.Error("File was not there, should error.")
	}

}

func TestStatusExisting(t *testing.T) {
	s := "bar"
	f, _ := os.Create(s)
	f.Close()
	time1, err := fileScan.Status(s)
	if err != nil {
		t.Error("File was there, should not have errored.")
	}
	time2, _ := fileScan.Status(s)
	os.Remove(s)
	if time1 != time2 {
		t.Error("Times should match.")
	}
}
