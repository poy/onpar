//go:build golden

package onpar_test

import (
	"flag"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

const goldenDir = "testdata"

var (
	update = flag.Bool("update", false, "update the golden files")

	// timingRe is a regular expression to find timings in test output.
	timingRe = regexp.MustCompile(`[0-9]+\.[0-9]{2,}s`)
)

func zeroTimings(line string) string {
	return timingRe.ReplaceAllString(line, "0.00s")
}

func goldenPath(path ...string) string {
	full := append([]string{goldenDir}, path...)
	return filepath.Join(full...)
}

func goldenFile(t *testing.T, path ...string) []byte {
	fullPath := goldenPath(path...)
	f, err := os.Open(fullPath)
	if err != nil {
		t.Fatalf("golden: could not open file %v for reading: %v", fullPath, err)
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatalf("golden: could not read file %v: %v", fullPath, err)
	}
	return b
}

func updateGoldenFile(t *testing.T, body []byte, path ...string) {
	fullPath := goldenPath(path...)
	f, err := os.Create(fullPath)
	if err != nil {
		t.Fatalf("golden: could not open file %v for writing: %v", fullPath, err)
	}
	defer f.Close()
	for len(body) > 0 {
		n, err := f.Write(body)
		if err != nil {
			t.Fatalf("golden: could not write to file %v: %v", fullPath, err)
		}
		body = body[n:]
	}
}

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func TestVerboseOutput(t *testing.T) {
	got, err := exec.Command("go", "test", "-v", "-tags", "goldenoutput", "-run", "TestNestedStructure").CombinedOutput()
	if err != nil {
		t.Fatalf("golden: tests failed: %v", err)
	}

	fn := "verbose.out"
	if *update {
		updateGoldenFile(t, got, fn)
		return
	}

	// All tests run in parallel, so we can't rely on the lines coming out in
	// the same order. But we _can_ test that all lines exist in the output,
	// complete with matching indentation.
	//
	// This is mainly to prove that t.Run calls are nested properly.
	want := goldenFile(t, fn)
	wantedLines := strings.Split(string(want), "\n")
	gotLines := strings.Split(string(got), "\n")

	if len(wantedLines) != len(gotLines) {
		t.Fatalf("expected %d lines of output; got %d", len(wantedLines), len(gotLines))
	}
	for _, wl := range wantedLines {
		wl = zeroTimings(wl)
		found := false
		for i, gl := range gotLines {
			gl = zeroTimings(gl)
			if wl == gl {
				gotLines = append(gotLines[:i], gotLines[i+1:]...)
				found = true
				break
			}
		}
		if !found {
			t.Errorf("missing line %v in actual output", wl)
		}
	}
	for _, l := range gotLines {
		t.Errorf("extra line %v in actual output", l)
	}
}
