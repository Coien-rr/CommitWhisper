package whisper

import (
	"os/exec"
	"strings"
)

func copyToClipboard(msg string) error {
	cmd := exec.Command("pbcopy")
	cmd.Stdin = strings.NewReader(msg)
	return cmd.Run()
}
