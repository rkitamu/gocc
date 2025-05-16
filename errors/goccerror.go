package errors

import (
	"fmt"
	"strings"
)

type PosError struct {
	Message string
	Input   string
	Pos     int
}

func NewPosError(message string, input string, pos int) *PosError {
	return &PosError{
		Message: message,
		Input:   input,
		Pos:     pos,
	}
}

func (e *PosError) Error() string {
	offset := e.Pos
	spaces := strings.Repeat(" ", offset)
	return fmt.Sprintf("%s\n%s^ %s", e.Input, spaces, e.Message)
}
