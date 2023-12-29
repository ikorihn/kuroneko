package core

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"

	"github.com/go-xmlfmt/xmlfmt"
)

// EditRequestBody edits request body using $EDITOR.
// If no EDITOR is specified, vim will open.
func EditRequestBody(curBody []byte) ([]byte, error) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	tempFile, err := os.CreateTemp("", "kuroneko-*")
	if err != nil {
		return nil, err
	}

	if len(curBody) > 0 {
		tempFile.Write(curBody)
	}

	tempFile.Close()

	defer os.Remove(tempFile.Name())

	cmd := exec.Command(editor, tempFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err = cmd.Run()

	body, err := os.ReadFile(tempFile.Name())
	body = bytes.Trim(body, "\n")

	return body, err
}

func FormatResponseBody(resp History) string {
	respContentType := resp.Header.Get(hdrContentType)
	switch {
	case reContentTypeJson.MatchString(respContentType):
		var buf bytes.Buffer
		err := json.Indent(&buf, resp.Body, "", "  ")
		if err == nil {
			return buf.String()
		}
	case reContentTypeXml.MatchString(respContentType):
		return xmlfmt.FormatXML(string(resp.Request.Body), "", "  ")
	}

	return ""
}
