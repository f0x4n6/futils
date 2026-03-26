// Package sys functions.
package sys

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"strings"
)

const (
	EX_OK       = 0
	EX_ERROR    = 1
	EX_USAGE    = 2
	EX_DATAERR  = 3
	EX_NOINPUT  = 4
	EX_NOTEXEC  = 126
	EX_NOTFOUND = 127
)

const (
	MODE_ALL  = 0777
	MODE_DIR  = 0755
	MODE_FILE = 0644
)

var Progress Any = Print

type Any func(a ...any)

func Print(a ...any) {
	_, _ = fmt.Fprintln(os.Stdout, a...)
}

func Error(a ...any) {
	_, _ = fmt.Fprintln(os.Stderr, a...)
}

func Final(a ...any) {
	_, _ = fmt.Fprintln(os.Stdout, a...)
	os.Exit(EX_OK)
}

func Fatal(a ...any) {
	_, _ = fmt.Fprintln(os.Stderr, a...)
	os.Exit(EX_ERROR)
}

func Usage(u string) {
	_, _ = fmt.Fprintln(os.Stdout, "Usage:", u)
	os.Exit(EX_USAGE)
}

func Debug(d string) {
	_, f, no, ok := runtime.Caller(1)

	if ok {
		_, _ = fmt.Fprintf(os.Stdout, "%s:%d %s\n", f, no, d)
	}
}

func Args() (args, xargs []string) {
	if flag.NArg() > 0 {
		args = flag.Args()
	} else {
		stdin, err, code := Stdin()

		if err != nil {
			Error(err)
			os.Exit(code)
		}

		if len(stdin) > 0 {
			args = strings.Split(stdin, "\n")
		}
	}

	i := slices.IndexFunc(args, func(a string) bool {
		return a == "--"
	})

	if i > -1 && i < len(args)-1 {
		return args[:i], args[i+1:]
	}

	return
}

func Stdin() (in string, err error, code int) {
	fi, err := os.Stdin.Stat()

	if err != nil {
		code = EX_NOINPUT
		return
	}

	if (fi.Mode() & os.ModeCharDevice) == 0 {
		var b []byte

		b, err = io.ReadAll(os.Stdin)

		if err != nil {
			code = EX_DATAERR
			return
		}

		in = strings.TrimSpace(string(b))
	}

	return
}

func StdCall(name string, args ...string) (string, error) {
	stdout := new(strings.Builder)
	stderr := new(strings.Builder)

	if call(stdout, stderr, name, args...) != 0 {
		return stdout.String(), errors.New(stderr.String())
	}

	return stdout.String(), nil
}

func ExitCall(name string, args ...string) {
	os.Exit(call(os.Stdout, os.Stderr, name, args...))
}

func call(stdout, stderr io.Writer, name string, args ...string) int {
	bin, err := exec.LookPath(name)

	if err != nil {
		Error(err)
		return EX_NOTFOUND
	}

	cmd := exec.Command(bin, args...)

	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err = cmd.Run()

	if ee, ok := errors.AsType[*exec.ExitError](err); ok {
		return ee.ExitCode()
	} else if err != nil {
		return EX_NOTEXEC
	}

	return EX_OK
}
