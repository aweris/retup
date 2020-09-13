// +build tools

// Place any runtime dependencies as imports in this file.
// Go modules will be forced to download and install them.
package main

//noinspection GoInvalidPackageImport
import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
)
