package controller

import (
	"io"
	"os"
	"os/exec"
)

func EditBody() ([]byte, error) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	tempFile, err := os.CreateTemp("", "kuroneko-*")
	if err != nil {
		return nil, err
	}

	defer func() {
		tempFile.Close()
		os.Remove(tempFile.Name())
	}()

	cmd := exec.Command(editor, tempFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err = cmd.Run()

	body, err := io.ReadAll(tempFile)

	return body, err
}
