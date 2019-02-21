//+build mage

package main

import (
	"errors"
	"fmt"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"os"
	"strings"
	"sync"
)

const goexe = "go"

//var Default = Build

func init() {
	// We want to use Go 1.11 modules even if the source lives inside GOPATH.
	// The default is "auto".
	os.Setenv("GO111MODULE", "on")
}

// Runs dep ensure and then installs the binary.
func Build() error {
	return sh.Run(goexe, "install", "./...")
}

// Cross build using gox
func Crossbuild() error {
	mg.Deps(gox)
	return sh.Run("gox", "-output", "dist/{{.Dir}}_{{.OS}}_{{.Arch}}", "./cmd/b2b")
}

func Clean() error {
	err := os.RemoveAll("dist")
	return err
}

// Run tests
func Test() error {
	return sh.Run(goexe, "test", "./...")
}

var (
	pkgPrefixLen = len("b2b-go")
	pkgs         []string
	pkgsInit     sync.Once
)

func packageList() ([]string, error) {
	var err error
	pkgsInit.Do(func() {
		var s string
		s, err = sh.Output(goexe, "list", "./...")
		if err != nil {
			return
		}
		pkgs = strings.Split(s, "\n")
		for i := range pkgs {
			pkgs[i] = "." + pkgs[i][pkgPrefixLen:]
		}
	})
	return pkgs, err
}

// List packages
func Packages() {
	pkgs, _ := packageList()
	for _, p := range pkgs {
		println(p)
	}
}

func gox() error {
	return sh.Run(goexe, "get", "-u", "github.com/mitchellh/gox")
}

func golint() error {
	return sh.Run(goexe, "get", "-u", "golang.org/x/lint/golint")
}

// Run golint linter
func Lint() error {
	mg.Deps(golint)

	pkgs, err := packageList()
	if err != nil {
		return err
	}
	failed := false
	for _, pkg := range pkgs {
		// We don't actually want to fail this target if we find golint errors,
		// so we don't pass -set_exit_status, but we still print out any failures.
		if _, err := sh.Exec(nil, os.Stderr, nil, "golint", pkg); err != nil {
			fmt.Printf("ERROR: running go lint on %q: %v\n", pkg, err)
			failed = true
		}
	}
	if failed {
		return errors.New("errors running golint")
	}
	return nil
}
