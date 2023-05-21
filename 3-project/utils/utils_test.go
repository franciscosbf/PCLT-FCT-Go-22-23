package utils

import (
	"bufio"
	"os"
	"strconv"
	"testing"
)

func TestFreshBuild(t *testing.T) {
	s := "foo"
	os.Remove(s)
	_, err := Build(s)
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
	f, _ := os.Create(s)
	f.WriteString(strconv.Itoa(10) + " times built.\n")
	f.Close()
	Build(s)
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
	os.Remove(s)
	_, err := Status(s)
	if err == nil {
		t.Error("File was not there, should error.")
	}

}

func TestStatusExisting(t *testing.T) {
	s := "bar"
	f, _ := os.Create(s)
	f.Close()
	time1, err := Status(s)
	if err != nil {
		t.Error("File was there, should not have errored.")
	}
	time2, _ := Status(s)
	os.Remove(s)
	if time1 != time2 {
		t.Error("Times should match.")
	}
}
