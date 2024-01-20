package checker

import (
	"bytes"
	"strings"
)

type (
	SqlValidationError struct {
		Line int
		Msg  string
	}

	SqlValidationErrors []SqlValidationError
)

func (s SqlValidationErrors) Error() string {
	buff := bytes.NewBufferString("")

	for i := 0; i < len(s); i++ {

		buff.WriteString(s[i].Error())
		buff.WriteString("\n")
	}

	return strings.TrimSpace(buff.String())
}

func (s SqlValidationError) Error() string {
	return s.Msg
}
