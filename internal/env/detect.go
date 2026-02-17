package env

import (
	"os"
	"runtime"
)

type Env struct {
	OS    string
	Shell string
}

func Detect() Env {
	shell := os.Getenv("SHELL") // Unix-like
	if shell == "" {
		shell = os.Getenv("COMSPEC") // Windows
	}
	if shell == "" {
		shell = "unknown"
	}

	return Env{
		OS:    runtime.GOOS,
		Shell: shell,
	}
}
