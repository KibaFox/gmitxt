// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Lint runs golangci-lint on the project.
func Lint() error {
	if err := run("golangci-lint", "run"); err != nil {
		return fmt.Errorf("problem linting project: %w", err)
	}

	if err := run("golangci-lint", "run", "magefile.go"); err != nil {
		return fmt.Errorf("problem linting magefile.go: %w", err)
	}

	return nil
}

// Test runs all tests for the project.
func Test() error {
	err := run("go", "test", "-coverprofile=coverage.out", "-bench=.", "./...")
	if err != nil {
		return fmt.Errorf("problem testing project: %w", err)
	}

	return nil
}

func run(cmd string, args ...string) (err error) {
	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin

	fmt.Println("exec:", cmd, strings.Join(args, " "))

	return c.Run()
}
