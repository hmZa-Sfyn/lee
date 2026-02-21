// errors.go
package main

import "fmt"

type EsoError struct {
	Msg  string
	Line int
	Col  int
}

func (e *EsoError) Error() string {
	if e.Line == 0 {
		return e.Msg
	}
	return fmt.Sprintf("%s (line %d, col %d)", e.Msg, e.Line, e.Col)
}

func NewEsoError(msg string, line, col int) *EsoError {
	return &EsoError{Msg: msg, Line: line, Col: col}
}

func NewEsoErrorf(line, col int, format string, args ...any) *EsoError {
	return NewEsoError(fmt.Sprintf(format, args...), line, col)
}
