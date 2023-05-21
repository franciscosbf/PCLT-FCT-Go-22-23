package parser

import "testing"

func TestBasic(t *testing.T) {
	s := "root <- dep1 dep2 dep3;"

	res, err := Parse(s)
	if err != nil {
		t.Error(err)
		return
	}
	if res.Rules[0].Object != "root" {
		t.Error("Head is not root.")
	}
	if len(res.Rules[0].Deps) != 3 {
		t.Error("Expected 3 deps, found", len(res.Rules[0].Deps))
	}
}

func TestMultiline1(t *testing.T) {
	s := `root <- dep1 dep2;dep1 <- dep3;dep2 <- dep3;`
	res, err := Parse(s)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if len(res.Rules) != 3 {
		t.Error("Failed to parse 3 rules")
		return
	}
	if res.Rules[0].Object != "root" {
		t.Error("First head is not root.")
	}
	if res.Rules[1].Object != "dep1" {
		t.Error("Second head is not dep1.")
	}
	if res.Rules[2].Object != "dep2" {
		t.Error("Second head is not dep2.")
	}
}

func TestMultiline2(t *testing.T) {
	s := `root <- dep1 dep2;
dep1 <- dep3;
dep2 <- dep3;`
	res, err := Parse(s)
	if err != nil {
		t.Error(err)
		return
	}
	if len(res.Rules) != 3 {
		t.Error("Failed to parse 3 rules")
		return
	}
	if res.Rules[0].Object != "root" {
		t.Error("First head is not root.")
	}
	if res.Rules[1].Object != "dep1" {
		t.Error("Second head is not dep1.")
	}
	if res.Rules[2].Object != "dep2" {
		t.Error("Third head is not dep2.")
	}
}

func TestParseFile(t *testing.T) {
	res, err := ParseFile("sample.df")
	if err != nil {
		t.Error(err)
		return
	}
	if len(res.Rules) != 4 {
		t.Error("Failed to parse 4 rules")
		return
	}
	if res.Rules[0].Object != "root" {
		t.Error("First head is not root.")
	}
	if res.Rules[1].Object != "dep1" {
		t.Error("Second head is not dep1.")
	}
	if res.Rules[2].Object != "dep2" {
		t.Error("Third head is not dep2.")
	}
	if res.Rules[3].Object != "dep3" {
		t.Error("Fourth head is not dep3.")
	}
}

func TestParseExt1(t *testing.T) {
	s := `root <- dep1.c dep2.h;dep1.c <- dep3.o;dep2.h <- dep3.o;`
	res, err := Parse(s)

	if err != nil {
		t.Error(err)
		return
	}
	if len(res.Rules) != 3 {
		t.Error("Failed to parse 3 rules")
		return
	}
	if res.Rules[0].Object != "root" {
		t.Error("First head is not root.")
	}
	if res.Rules[1].Object != "dep1.c" {
		t.Error("Second head is not dep1.c.")
	}
	if res.Rules[2].Object != "dep2.h" {
		t.Error("Second head is not dep2.h.")
	}
}

func TestParseExt2(t *testing.T) {
	s := `root <- dep1.c dep2.h;
dep1.c <- dep3.o;
dep2.h <- dep3.o;
`
	res, err := Parse(s)

	if err != nil {
		t.Error(err)
		return
	}
	if len(res.Rules) != 3 {
		t.Error("Failed to parse 3 rules")
		return
	}
	if res.Rules[0].Object != "root" {
		t.Error("First head is not root.")
	}
	if res.Rules[1].Object != "dep1.c" {
		t.Error("Second head is not dep1.c.")
	}
	if res.Rules[2].Object != "dep2.h" {
		t.Error("Second head is not dep2.h.")
	}
}
