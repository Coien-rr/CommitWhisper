package git

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
)

func GetGitDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--staged")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()

	return assertDiff(out.String(), err)
}

func assertDiff(diff string, err error) (string, error) {
	if err != nil {
		return "", fmt.Errorf("ERROR: Run 'git diff HEAD' get: %v", err)
	}

	if diff == "" {
		return "", errors.New("ERROR: Run 'git diff HEAD' get empty res")
	}

	return diff, nil
}
