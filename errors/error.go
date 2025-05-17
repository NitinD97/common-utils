package errors

import (
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

type Tracer struct {
	cause           error
	originalMessage string
	userMessage     string
	stackTrace      []uintptr
	isWrapped       bool
}

// New creates a new error with captured stack trace.
func New(message string) error {
	return &Tracer{
		originalMessage: message,
		userMessage:     "",
		stackTrace:      callers(),
		isWrapped:       false,
	}
}

// Wrap wraps an existing error with a new message and stack trace.
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}

	return &Tracer{
		cause:           err,
		originalMessage: err.Error(),
		userMessage:     message,
		stackTrace:      callers(),
		isWrapped:       true,
	}
}

// Message extracts the top-level user message from the Tracer.
func Message(err error) string {
	if err == nil {
		return ""
	}

	var et *Tracer
	if errors.As(err, &et) {
		return et.userMessage
	}
	return err.Error()
}

// Error implements the error interface.
func (e *Tracer) Error() string {
	if e.userMessage != "" {
		return e.userMessage + ": " + e.originalMessage
	}
	return e.originalMessage
}

// Unwrap returns the cause of the error if it's a Tracer.
func Unwrap(err error) error {
	if err == nil {
		return nil
	}
	var et *Tracer
	if errors.As(err, &et) {
		return et.cause
	}
	return nil
}

func Cause(err error) error {
	if err == nil {
		return nil
	}

	// loop through the error chain, unwrapping until we find the root cause
	for {
		terr, ok := err.(*Tracer)
		if !ok {
			return err
		}

		if terr.isWrapped {
			err = terr.cause
			continue
		}

		return err
	}
}

// Format implements fmt.Formatter for %+v formatting.
func (e *Tracer) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		if f.Flag('+') {
			if len(e.stackTrace) > 0 {
				fmt.Fprintf(f, "%s\n", e.Error())
				fmt.Fprintf(f, "%s", e.printStack())
			}

			if e.cause != nil {
				if tracer, ok := e.cause.(*Tracer); ok {
					fmt.Fprintf(f, "%+v", tracer)
				} else {
					fmt.Fprintf(f, "%v", e.cause)
				}
			}
			return
		} else {
			fmt.Fprint(f, e.Error())
		}
	case 's':
		fmt.Fprint(f, e.Error())
	case 'q':
		fmt.Fprintf(f, "%q", e.Error())
	}
}

// printStack returns a formatted string of the captured stack trace.
func (e *Tracer) printStack() string {
	var sb strings.Builder

	for _, pc := range e.stackTrace {
		fn := runtime.FuncForPC(pc - 1)
		if fn == nil {
			continue
		}
		file, line := fn.FileLine(pc - 1)
		sb.WriteString(fn.Name())
		sb.WriteString("\n\t")
		sb.WriteString(file)
		sb.WriteString(":")
		sb.WriteString(strconv.Itoa(line))
		sb.WriteString("\n")
	}

	return sb.String()
}

// callers captures the call stack.
func callers() []uintptr {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	return pcs[:n]
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target *interface{}) bool {
	return errors.As(err, target)
}
