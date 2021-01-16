package modutil

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
)

// exe runs a program and returns stdout, stderr, and any
// error that might have occurred
func exe(name string, arg ...string) (stdout, stderr string, err error) {
	cmd := exec.Command(name, arg...) //nolint: gosec,gocritic // it's fine, this is not a user input

	// we pipe stderr to get the error message if something goes wrong
	stderrReader, err := cmd.StderrPipe()
	if err != nil {
		return "", "", fmt.Errorf("could not pipe stderr: %w", err)
	}
	// we pipe stdout to get the output of the script
	stdoutReader, err := cmd.StdoutPipe()
	if err != nil {
		return "", "", fmt.Errorf("could not pipe stdout: %w", err)
	}

	// We start the command
	if err = cmd.Start(); err != nil {
		return "", "", fmt.Errorf("could not run the command: %w", err)
	}

	// we read all stderr to get the error message (if any)
	stderrByte, err := ioutil.ReadAll(stderrReader)
	if err != nil {
		return "", "", fmt.Errorf("could not read stderr: %w", err)
	}
	// we read all stdout to get the output of the script (if any)
	stdoutByte, err := ioutil.ReadAll(stdoutReader)
	if err != nil {
		return "", "", fmt.Errorf("could not read stdout: %w", err)
	}

	stdout = strings.TrimSuffix(string(stdoutByte), "\n")
	stderr = strings.TrimSuffix(string(stderrByte), "\n")
	return stdout, stderr, cmd.Wait()
}

// run runs a program and returns stderr as error
func run(name string, arg ...string) (string, error) {
	stdout, stderr, err := exe(name, arg...)

	if err != nil && stderr != "" {
		return stdout, errors.New(stderr)
	}
	return stdout, err
}
