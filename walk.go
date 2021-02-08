package gersby

import (
	"github.com/go-git/go-billy/v5"
	"os"
	"path/filepath"
)

func Walk(fs billy.Filesystem, root string, walkFn filepath.WalkFunc) error {
	info, err := fs.Lstat(root)
	if err != nil {
		err = walkFn(root, nil, err)
	} else {
		err = walk(fs, root, info, walkFn)
	}
	if err == filepath.SkipDir {
		return nil
	}
	return err
}

func walk(fs billy.Filesystem, path string, info os.FileInfo, walkFn filepath.WalkFunc) error {
	if !info.IsDir() {
		return walkFn(path, info, nil)
	}

	files, err := fs.ReadDir(path)
	err1 := walkFn(path, info, err)
	// If err != nil, walk can't walk into this directory.
	// err1 != nil means walkFn want walk to skip this directory or stop walking.
	// Therefore, if one of err and err1 isn't nil, walk will return.
	if err != nil || err1 != nil {
		// The caller's behavior is controlled by the return value, which is decided
		// by walkFn. walkFn may ignore err and return nil.
		// If walkFn returns SkipDir, it will be handled by the caller.
		// So walk should return whatever walkFn returns.
		return err1
	}

	for _, file := range files {
		filename := filepath.Join(path, file.Name())
		if err != nil {
			if err := walkFn(filename, file, err); err != nil && err != filepath.SkipDir {
				return err
			}
		} else {
			err = walk(fs, filename, file, walkFn)
			if err != nil {
				if !file.IsDir() || err != filepath.SkipDir {
					return err
				}
			}
		}
	}
	return nil
}
