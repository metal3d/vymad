// Lanch pandoc and get results
package pandoc

import (
	"fmt"
	"io/ioutil"
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
	stdin.Close()

	if err != nil {
		fmt.Println(err)
	}
	o, err := ioutil.ReadAll(stdout)
	if err != nil {
		return "", err
	}

	return string(o), nil
}
