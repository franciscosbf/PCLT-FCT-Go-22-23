package builder

import (
	"cpl_go_proj22/parser"
	"errors"
	"fmt"
	"testing"
	"time"
)

var missing = errors.New("missing file")

type buildError struct {
	file string
}

func (e *buildError) Error() string {
	return ""
}

type fakeFileInfo struct {
	fail bool
	time *time.Time
}

// Scan represents the worst mock-up ever seen
type fakeScan struct {
	files map[string]*fakeFileInfo
}

func (s *fakeScan) Status(filename string) (time.Time, error) {
	info := s.files[filename]
	if info.time == nil {
		return time.Time{}, missing
	}

	return *info.time, nil
}

func (s *fakeScan) Build(filename string) (time.Time, error) {
	info := s.files[filename]
	if info.fail {
		return time.Time{}, &buildError{file: filename}
	}

	return time.Now(), nil // current time used to force build on dependants
}

// convertTime expects convertTime in the format %d%d
func convertTime(day string) *time.Time {
	t, _ := time.Parse("2001-01-01", fmt.Sprintf("2001-01-%s", day))
	return &t
}

func TestGraphContent(t *testing.T) {
	s := `
r  <- d1 d2;
d1 <- d3;
d2 <- d3 d4;
`

	dFile, _ := parser.Parse(s)

	dG := buildGraph(dFile)
	if dG == nil {
		t.Fatal("Graph is nil")
	}
	
	if len(dG.nodes) != 5 {
		t.Errorf("Invalid nodes map. expect=[r, d1 ,d2, d3, d4] got=%v", dG.nodes)
	}
	for _, f := range []string{"r", "d1", "d2", "d3", "d4"} {
		if _, ok := dG.nodes[f]; !ok {
			t.Errorf("Missing %q in nodes map", f)
		}
	} 
	
	if len(dG.leafs) != 2 {
		t.Errorf("Invalid nodes map. expect=[d3, d4] got=%v", dG.leafs)
	}
	for _, f := range []string{"d3", "d4"} {
		if _, ok := dG.leafs[f]; !ok {
			t.Errorf("Missing %q in leafs map", f)
		}
	} 
	
	if len(dG.targets) != 3 {
		t.Errorf("Invalid targets map. expect=[r, d1, d2] got=%v", dG.targets)
	}
	for _, f := range []string{"r", "d1", "d2"} {
		var found bool
		for _, t := range dG.targets {
			if t.filename == f {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Missing target %q", f)
		}
	} 
}

func TestFilesInfo(t *testing.T) {
	s := `
r  <- d1 d2;
d1 <- d3;
d2 <- d3 d4 d5;
`

	dFile, _ := parser.Parse(s)

	dG := buildGraph(dFile)
	if dG == nil {
		t.Fatal("Graph is nil")
	}

	checkInfo := func(filename string, dependencies int, dependants ...string) {
		info := dG.nodes[filename]
		if info == nil {
			t.Errorf("Info of %q is nil", filename)
			return
		}

		if info.filename != filename {
			t.Errorf(
				"Wrong filename. got=%q, expect=%q", 
				info.filename, filename,
			)
		}

		if info.dependencies != dependencies {
			t.Errorf(
				"Wrong number of dependencies of %q. got=%d, expect=%d", 
				filename, info.dependencies, dependencies,
			)
		} 

		if len(info.dependants) != len(dependants) {
			t.Errorf(
				"Wrong dependants of %q. got=%v, expect=%v", 
				filename, info.dependants, dependants,
			)
		}
		for _, dep := range dependants {
			var found bool
			for _, d := range info.dependants {
				if dep == d {
					found = true
					break
				}
			}	
			if !found {
				t.Errorf("Missing dependant %q of %q", dep, filename)
			}
		}
	}

	checkInfo("r", 2)
	checkInfo("d1", 1, "r")
	checkInfo("d2", 3, "r")
	checkInfo("d3", 0, "d1", "d2")
	checkInfo("d4", 0, "d2")
	checkInfo("d5", 0, "d2")
}

func TestBuildWihoutErrors(t *testing.T) {
	fileScan := &fakeScan{
		files: map[string]*fakeFileInfo{
			 "r": {time: convertTime("01")},
			"d1": {time: convertTime("04")},
			"d2": {time: convertTime("03")},
			"d3": {time: convertTime("06")},
			"d4": {},
			"d5": {},
			"d6": {time: convertTime("02")},
			"d7": {},
			"d8": {time: convertTime("01")},
		},
	}

	s := `
r  <- d1 d3 d5 d4 d8;
d1 <- d2 d7;
d3 <- d1 d4 d5 d7;
d5 <- d6 d8;
`

	dFile, _ := parser.Parse(s)

	tunnel := MakeController(dFile, fileScan)
	if tunnel == nil {
		t.Fatal("Channel is nil")
	}

	msg := <-tunnel
	if msg == nil {
		t.Fatal("Message is nil")
	}

	if msg.Type != BuildSuccess {
		t.Fatalf("Got an unnexpected error: %v", msg.Err)
	}
}

func TestBuildWithErrors(t *testing.T) {
	fileScan := &fakeScan{
		files: map[string]*fakeFileInfo{
			 "r": {time: convertTime("01")},
			"d1": {time: convertTime("04"), fail: true},
			"d2": {time: convertTime("03")},
			"d3": {time: convertTime("06")},
			"d4": {},
			"d5": {fail: true},
			"d6": {time: convertTime("02")},
			"d7": {},
			"d8": {time: convertTime("01")},
		},
	}

	s := `
r  <- d1 d3 d5 d4 d8;
d1 <- d2 d7;
d3 <- d1 d4 d5 d7;
d5 <- d6 d8;
`

	dFile, _ := parser.Parse(s)

	tunnel := MakeController(dFile, fileScan)
	if tunnel == nil {
		t.Fatal("Channel is nil")
	}

	msg := <-tunnel
	if msg == nil {
		t.Fatal("Message is nil")
	}

	if msg.Type != BuildError {
		t.Fatal("Expecting message of type BuildError")
	}

	err, ok := msg.Err.(*buildError)
	if !ok {
		t.Fatalf("Err isn't of type buildError: got=%v", err)
	}

	if err.file != "d1" && err.file != "d5" {
		t.Fatalf("Expecting build error from d1 or d5. got=%s", err.file)
	}
}

