package builder

import (
	"cpl_go_proj22/parser"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"testing"
	"time"
)

func init() {
	if value := os.Getenv("DISABLE_LOG"); value != "" {
		log.SetOutput(io.Discard)
	}
}

var missing = errors.New("missing file")

type buildError struct {
	filename string
}

func (e *buildError) Error() string {
	return fmt.Sprintf("build error on file %q", e.filename)
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
		return time.Time{}, &buildError{filename: filename}
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

func TestBuildWithoutErrors(t *testing.T) {
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

	if err.filename != "d1" && err.filename != "d5" {
		t.Fatalf("Expecting build error from d1 or d5. got=%s", err.filename)
	}
}

// TestBuildWithBigGraph assumeds that
// everything is ok
func TestBuildWithBigDepGraph(t *testing.T) {
	fileScan := &fakeScan{
		files: map[string]*fakeFileInfo{
			  "r": {time: convertTime("20")},
			 "d1": {time: convertTime("04")},
			 "d2": {time: convertTime("11")},
			 "d3": {time: convertTime("06")},
			 "d4": {},
			 "d5": {},
			 "d6": {time: convertTime("16")},
			 "d7": {},
			 "d8": {time: convertTime("01")},
			 "d9": {time: convertTime("04")},
			"d10": {time: convertTime("13")},
			"d11": {time: convertTime("06")},
			"d12": {},
			"d13": {},
			"d14": {time: convertTime("02")},
			"d15": {},
			"d16": {time: convertTime("01")},
			"d17": {time: convertTime("26")},
			"d18": {time: convertTime("03")},
			"d19": {time: convertTime("14")},
			"d20": {},
			"d21": {},
			"d22": {time: convertTime("02")},
			"d23": {},
			"d24": {time: convertTime("01")},
			"d25": {},
			"d26": {},
			"d27": {time: convertTime("06")},
			"d28": {},
			"d29": {time: convertTime("01")},
			"d30": {time: convertTime("20")},
			"d31": {},
			"d32": {time: convertTime("02")},
			"d33": {},
			"d34": {time: convertTime("24")},
			"d35": {},
			"d36": {},
			"d37": {time: convertTime("02")},
			"d38": {},
			"d39": {time: convertTime("25")},
			"d40": {time: convertTime("01")},
			"d41": {time: convertTime("14")},
			"d42": {time: convertTime("01")},
			"d43": {time: convertTime("15")},
			"d44": {time: convertTime("16")},
			"d45": {time: convertTime("01")},
		},
	}

	s := `
r   <- d14 d13 d5;
d5  <- d12 d24 d1;
d13 <- d5 d11 d15 d33;
d14 <- d32 d33;
d33 <- d15 d44 d19 d32;
d1  <- d2 d45 d6;
d24 <- d6 d10;
d12 <- d11;
d15 <- d16 d44;
d44 <- d19;
d19 <- d16 d43 d20 d34;
d2  <- d3 d26;
d45 <- d26;
d6  <- d26 d25;
d11 <- d10 d42 d16;
d20 <- d35;
d43 <- d31;
d16 <- d31;
d31 <- d21 d35;
d21 <- d35 d36;
d42 <- d21;
d10 <- d17 d41;
d25 <- d10 d41 d9;
d26 <- d7;
d3  <- d27 d7;
d7  <- d27 d4 d28 d29;
d28 <- d40;
d4  <- d40;
d29 <- d8;
d8  <- d39;
d9  <- d23;
d41 <- d23;
d17 <- d30;
d30 <- d18 d36;
d36 <- d37;
d18 <- d37;
d23 <- d8 d38 d37 d30;
`
	
	dFile, _ := parser.Parse(s)
	
	tunnel := MakeController(dFile, fileScan)
	<-tunnel
}

