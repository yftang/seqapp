package rna_test

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/yftang/seqapp/batch/rna"
)

// ComparisonSliceEqual determines whether two `Comparison` slices are equal
func ComparisonSliceEqual(a, b []Comparison) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil || len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i].Case != b[i].Case || a[i].Control != b[i].Control {
			return false
		}
	}

	return true
}

// GroupSliceEqual determines whether two `Group` slices are equal
func GroupSliceEqual(a, b []Group) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil || len(a) != len(b) {
		return false
	}

	for i := range a {
		GroupEqual(a[i], b[i])
	}

	return true
}

// GroupEqual determines where two `Group`s are equal
func GroupEqual(g1, g2 Group) bool {
	if fmt.Sprintf("%T", g1) != "Group" || fmt.Sprintf("%T", g2) != "Group" {
		return false
	}
	if g1.Name != g2.Name {
		return false
	}
	if !StringSliceEqual(g1.Samples, g2.Samples) {
		return false
	}

	return true
}

// StringSliceEqual determines whether two slices are equal
func StringSliceEqual(a, b []string) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil || len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

// TestDiffConfigFile404 tests if `DiffConfig` throws `file not found` error when
// perform `Parse` method with non-exist `ConfigFile`.
func TestDiffConfigFile404(t *testing.T) {
	dc := DiffConfig{ConfigFile: "file_not_exist.json"}

	err := dc.Parse()

	if err == nil || strings.Contains(fmt.Sprintf("%v", err), "file not found") {
		t.Error("DiffConfig should raise 'EOF' error if ConfigFile not existes.")
	}
}

// TestDiffConfig tests `DiffConfig` functionality.
func TestDiffConfig(t *testing.T) {
	testConfigFile := "../../test_files/diff_config_test.json"
	absPath, _ := filepath.Abs(testConfigFile)
	absPath = filepath.Clean(absPath)

	dc := DiffConfig{ConfigFile: testConfigFile}
	err := dc.Parse()
	if err != nil {
		t.Error("DiffConfig should not throw any error when `Parse()` with the correct `ConfigFile`.")
	}

	g1 := Group{Name: "Oh-U", Samples: []string{"CHG016348", "CHG016351", "CHG016354"}}
	g2 := Group{Name: "Oh-M", Samples: []string{"CHG016349", "CHG016352", "CHG016355"}}
	g3 := Group{Name: "Oh-B", Samples: []string{"CHG016350", "CHG016353", "CHG016356"}}
	c1 := Comparison{Case: g1.Name, Control: g2.Name}
	c2 := Comparison{Case: g1.Name, Control: g3.Name}
	c3 := Comparison{Case: g2.Name, Control: g3.Name}
	if dc.CuffdiffExec != "/online/software/cuffdiff" {
		t.Error("CuffdiffExec should be '/online/software/cuffdiff', is ", dc.CuffdiffExec)
	}
	if dc.ConfigFile != absPath {
		t.Error("ConfigFile should be ", absPath, " is ", dc.ConfigFile)
	}
	if dc.ProjectPath != "/home/sam/GoTest" {
		t.Error("ProjectPath should be '/home/sam/GoTest', is ", dc.ProjectPath)
	}
	if !GroupSliceEqual(dc.Groups, []Group{g1, g2, g3}) {
		t.Error("Groups is ", dc.Groups)
	}
	if !ComparisonSliceEqual(dc.Comparisons, []Comparison{c1, c2, c3}) {
		t.Error("Comparisons is ", dc.Comparisons)
	}
}
