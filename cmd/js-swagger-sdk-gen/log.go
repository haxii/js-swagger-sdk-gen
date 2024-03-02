package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func fprint(w io.Writer, format string, a ...any) {
	if !strings.HasSuffix(format, "\n") {
		format = format + "\n"
	}
	fmt.Fprintf(w, format, a...)
}

func log(format string, a ...any) {
	fprint(os.Stdout, format, a...)
}

func debug(format string, a ...any) {
	if opt.Verbose {
		log(format, a...)
	}
}

func warn(format string, a ...any) {
	fprint(os.Stderr, format, a...)
}

func fatal(format string, a ...any) {
	warn(format, a...)
	os.Exit(1)
}
