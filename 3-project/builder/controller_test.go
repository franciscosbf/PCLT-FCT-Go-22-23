package builder

import (
	"cpl_go_proj22/parser"
	"log"
	"testing"
)

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
		log.Printf("Invalid nodes map. expect=[d3, d4] got=%v", dG.leafs)
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

