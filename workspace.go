package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/pkg/fileutils"
	"github.com/karrick/godirwalk"
	"github.com/pkg/errors"
)

// Copied from https://raw.githubusercontent.com/GoogleContainerTools/skaffold/987ec17f8f6a58df1d615f1c3731138eb72f6615/pkg/skaffold/docker/dependencies.go

// WalkWorkspace walks the given host directories and records all files found.
// nolint:gocognit
func WalkWorkspace(workspace string, excludes, deps []string) (map[string]bool, error) {
	pExclude, err := fileutils.NewPatternMatcher(excludes)
	if err != nil {
		return nil, errors.Wrap(err, "invalid exclude patterns")
	}

	// Walk the workspace
	files := make(map[string]bool)

	for _, dep := range deps {
		dep = filepath.Clean(dep)
		absDep := filepath.Join(workspace, dep)

		fi, err := os.Stat(absDep)
		if err != nil {
			return nil, errors.Wrapf(err, "stating file %s", absDep)
		}

		switch mode := fi.Mode(); {
		case mode.IsDir():
			if err := godirwalk.Walk(absDep, &godirwalk.Options{
				Unsorted: true,
				Callback: func(fpath string, info *godirwalk.Dirent) error {
					if fpath == absDep {
						return nil
					}

					relPath, err := filepath.Rel(workspace, fpath)
					if err != nil {
						return err
					}

					ignored, err := pExclude.Matches(relPath)
					if err != nil {
						return err
					}

					if info.IsDir() {
						if !ignored {
							return nil
						}
						// exclusion handling closely follows vendor/github.com/docker/docker/pkg/archive/archive.go
						// No exceptions (!...) in patterns so just skip dir
						if !pExclude.Exclusions() {
							return filepath.SkipDir
						}

						dirSlash := relPath + string(filepath.Separator)

						for _, pat := range pExclude.Patterns() {
							if !pat.Exclusion() {
								continue
							}
							if strings.HasPrefix(pat.String()+string(filepath.Separator), dirSlash) {
								// found a match - so can't skip this dir
								return nil
							}
						}

						return filepath.SkipDir
					} else if !ignored {
						files[relPath] = true
					}

					return nil
				},
			}); err != nil {
				return nil, errors.Wrapf(err, "walking folder %s", absDep)
			}
		case mode.IsRegular():
			ignored, err := pExclude.Matches(dep)
			if err != nil {
				return nil, err
			}

			if !ignored {
				files[dep] = true
			}
		}
	}

	return files, nil
}
