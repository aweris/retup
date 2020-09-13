package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/docker/docker/pkg/fileutils"
	"github.com/spf13/pflag"
)

var (
	// version flags.
	version = "dev"
	// nolint:gochecknoglobals
	commit = "none"
	// nolint:gochecknoglobals
	date = "unknown"
	// nolint:gochecknoglobals
	builtBy = "unknown"
)

func main() {
	var (
		// config flags
		cfgFile string
		distDir string

		// other
		showVersion bool
	)

	pflag.StringVar(&cfgFile, "config", "./retup.yaml", "config file")
	pflag.StringVar(&distDir, "output", "dist", "output directory")

	pflag.BoolVar(&showVersion, "version", false, "Prints version info")

	bindEnv(pflag.Lookup("config"), "RETUP_CONFIG")
	bindEnv(pflag.Lookup("output"), "RETUP_OUTPUT")

	pflag.Parse()

	if showVersion {
		fmt.Printf("Version    : %s\n", version)
		fmt.Printf("Git Commit : %s\n", commit)
		fmt.Printf("Build Date : %s\n", date)
		fmt.Printf("Build By   : %s\n", builtBy)
		os.Exit(0)
	}

	cfg, err := NewConfig(cfgFile)
	if err != nil {
		log.Fatal(err)
	}

	// create brand new output directory
	if err := ensureNewDir(distDir); err != nil {
		log.Fatalf("failed to create dir: %s err: %v\n", distDir, err)
	}

	for _, artifact := range cfg.Artifacts {
		// create output directory for artifact
		artifactDir := path.Join(distDir, artifact.Name)
		if err := ensureNewDir(artifactDir); err != nil {
			log.Fatalf("failed to create dir: %s err: %v\n", distDir, err)
		}

		// absolute path of the artifactDir
		artifactDirAbs, err := filepath.Abs(artifactDir)
		if err != nil {
			log.Fatalf("failed to get absolute path: %s err: %v\n", artifactDir, err)
		}

		// absolute path for the artifact context
		context, err := filepath.Abs(artifact.Context)
		if err != nil {
			log.Fatalf("failed to get context path: %s err: %v\n", context, err)
		}

		// get file list for the artifact
		files, err := WalkWorkspace(context, artifact.Dependencies.Ignore, artifact.Dependencies.Paths)
		if err != nil {
			log.Fatalf("failed to walk context dependencies: %s err: %v\n", context, err)
		}

		// copy files to output directory with relative paths
		for file := range files {
			src := path.Join(context, file)
			dst := path.Join(artifactDirAbs, file)

			if err := fileutils.CreateIfNotExists(dst, false); err != nil {
				log.Fatalf("failed to create file: %s err: %v\n", dst, err)
			}

			if _, err := fileutils.CopyFile(src, dst); err != nil {
				log.Fatalf("failed to copy file: %s err: %v\n", src, err)
			}
		}
	}
}

func bindEnv(fn *pflag.Flag, env string) {
	if fn == nil || fn.Changed {
		return
	}

	val := os.Getenv(env)

	if len(val) > 0 {
		if err := fn.Value.Set(val); err != nil {
			log.Fatalf("failed to bind env: %v\n", err)
		}
	}
}

// ensures given directory exist and empty.
func ensureNewDir(path string) error {
	// remove directory if exist
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		if err := os.RemoveAll(path); err != nil {
			log.Fatalf("failed to remove path: %s err: %v\n", path, err)
		}
	}
	// create new directory
	if err := os.Mkdir(path, os.ModePerm); err != nil {
		return err
	}

	return nil
}
