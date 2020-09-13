package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config contains all the configuration for the build.
type Config struct {
	// Artifacts lists the images you're going to be building.
	Artifacts []*Artifact `yaml:"artifacts,omitempty"`
}

// Artifact are the items that need to be built, along with the context in which they should be built.
type Artifact struct {
	// Name of the artifact
	Name string `yaml:"name,omitempty"`

	// Context is the directory containing the artifact's sources. Defaults to `.`.
	Context string `yaml:"context,omitempty"`

	// Dependencies are the file dependencies should watch for building.
	Dependencies *Dependencies `yaml:"dependencies,omitempty"`
}

// Dependencies is used to specify dependencies for an service built.
type Dependencies struct {
	// Paths should be set to the file dependencies for this artifact.
	Paths []string `yaml:"paths,omitempty"`

	// Ignore specifies the paths that should be ignored. If a file exists in both `paths` and in `ignore`, it will be ignored.
	// Will only work in conjunction with `paths`.
	Ignore []string `yaml:"ignore,omitempty"`
}

// NewConfig returns a new decoded Config struct.
func NewConfig(path string) (*Config, error) {
	// Create config structure
	config := &Config{}

	// Open config file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
