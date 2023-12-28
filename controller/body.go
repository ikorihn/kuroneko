package controller

import (
	"bytes"
	"encoding/json"
	"regexp"

	"github.com/go-xmlfmt/xmlfmt"
)

var (
	reContentTypeJson = regexp.MustCompile(`(?i:(application|text)/(.*json.*)(;|$))`)
	reContentTypeXml  = regexp.MustCompile(`(?i:(application|text)/(.*xml.*)(;|$))`)
)

func FormatBody(resp History) string {
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
