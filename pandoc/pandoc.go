// Package pandoc calls the pandoc command line tool to convert content from one format to another.
package pandoc

import (
	"io"
	"os/exec"
)

// Launch pandoc and returns content translated to given format
func Launch(content, format string) (string, error) {
	cmd := exec.Command("pandoc", "-f", "html", "-t", format)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	if err := cmd.Start(); err != nil {
		return "", err
	}

	stdin.Write([]byte(content))
	defer stdin.Close()

	o, err := io.ReadAll(stdout)
	if err != nil {
		return "", err
	}

	return string(o), nil
}
