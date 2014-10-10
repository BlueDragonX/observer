package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Attempt to determine the path to a provider.
func lookupProvider(name string) (path string, err error) {
	name = fmt.Sprintf("observer-%s", name)
	if path, err = exec.LookPath(name); err == nil {
		return
	}

	check := func(path string) (err error) {
		var fi os.FileInfo
		if fi, err = os.Stat(path); err == nil {
			if fi.Mode() & 111 == 0 {
				err = fmt.Errorf("%s is not executable", path)
			}
		}
		return
	}

	path = filepath.Join(filepath.Dir(os.Args[0]), name)
	if err = check(path); err == nil {
		return
	}

	var cwd string
	if cwd, err = os.Getwd(); err != nil {
		return
	}

	path = filepath.Join(cwd, name)
	if err = check(path); err == nil {
		return
	}

	err = fmt.Errorf("unable to find provider binary %s\n", name)
	return
}
