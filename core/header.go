package core

import (
	"fmt"
	"strings"
)

type headerMap map[string]string

func HeaderMapFromList(headerList []string) headerMap {
	result := make(headerMap)
	for _, v := range headerList {
		sp := strings.Split(v, ":")
		result[sp[0]] = sp[1]
	}
	return result
}

func (h headerMap) ToList() []string {
	result := make([]string, 0)
	for k, v := range h {
		result = append(result, fmt.Sprintf("%s:%s", k, v))
	}
	return result
}
func (h headerMap) AddNameValue(nameValue string) {
	sp := strings.Split(nameValue, ":")
	h[sp[0]] = sp[1]
}
func (h headerMap) RemoveNameValue(nameValue string) {
	sp := strings.Split(nameValue, ":")
	delete(h, sp[0])
}
