package gersby

import (
	"fmt"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/google/go-cmp/cmp"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestName(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("could not get testdata directory")
	}
	testdata := path.Join(path.Dir(filename), "testdata")

	tests := []struct {
		desc   string
		walkFn func(tracker *walkTracker) filepath.WalkFunc
	}{
		{"simple walk", simpleWalk},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			fs := osfs.New(filepath.Join(testdata))
			billyTracker, filePathTracker := &walkTracker{}, &walkTracker{}
			billyWalk, filepathWalk := test.walkFn(billyTracker), test.walkFn(filePathTracker)

			fpErr := filepath.Walk(testdata, filepathWalk)
			billyErr := Walk(fs, "/", billyWalk)

			filePathTracker.normalize(testdata)

			if !cmp.Equal(billyErr, fpErr) {
				tt.Errorf("\nwanted %v\n   got %v", fpErr, billyErr)
			}
			if !cmp.Equal(billyTracker.Invocations, filePathTracker.Invocations, cmp.Comparer(compareStats)) {
				tt.Errorf("\nwanted %v\n   got %v", filePathTracker.Invocations, billyTracker.Invocations)
			}
		})
	}

}

func simpleWalk(tracker *walkTracker) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		tracker.Invocations = append(tracker.Invocations, &walkInvocation{
			Path: path,
			Info: info,
			Err:  err,
		})
		return nil
	}
}

type walkTracker struct {
	Invocations []*walkInvocation
}

type walkInvocation struct {
	Path string
	Info os.FileInfo
	Err  error
}

func compareStats(a, b os.FileInfo) bool {
	return a.Name() == b.Name() && a.Mode() == b.Mode() && a.ModTime() == b.ModTime() && a.Size() == b.Size() && a.IsDir() == b.IsDir() && cmp.Equal(a.Sys(), b.Sys())
}

func (w walkInvocation) String() string {
	return fmt.Sprintf("{ path: %s, info: %+v, err: %+v }", w.Path, w.Info, w.Err)
}

func (t *walkTracker) normalize(root string) {
	for _, invocation := range t.Invocations {
		invocation.Path = filepath.Clean(strings.ReplaceAll(invocation.Path, root, "/"))
	}
}
