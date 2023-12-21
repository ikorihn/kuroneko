package controller

import (
	"os"
	"os/exec"
)

func (c *Controller) EditBody() ([]byte, error) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	tempFile, err := os.CreateTemp("", "kuroneko-*")
	if err != nil {
		return nil, err
	}
	tempFile.Close()

	defer os.Remove(tempFile.Name())

	cmd := exec.Command(editor, tempFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err = cmd.Run()

	body, err := os.ReadFile(tempFile.Name())

	return body, err
}
