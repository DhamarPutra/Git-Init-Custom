package logger

import (
	"fmt"
	"io"
	"os"
)

// Logger handles output and verbose logging.
type Logger struct {
	stdout  io.Writer
	stderr  io.Writer
	verbose bool
}

// New creates a new Logger instance.
func New(verbose bool) *Logger {
	return &Logger{
		stdout:  os.Stdout,
		stderr:  os.Stderr,
		verbose: verbose,
	}
}

// Info logs messages to stdout.
func (l *Logger) Info(format string, a ...interface{}) {
	fmt.Fprintf(l.stdout, format+"\n", a...)
}

// Error logs messages to stderr.
func (l *Logger) Error(format string, a ...interface{}) {
	fmt.Fprintf(l.stderr, "Error: "+format+"\n", a...)
}

// Verbose logs messages to stdout if verbose mode is enabled.
func (l *Logger) Verbose(format string, a ...interface{}) {
	if l.verbose {
		fmt.Fprintf(l.stdout, format+"\n", a...)
	}
}
