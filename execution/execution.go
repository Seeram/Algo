package execution

import (
	"fmt"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"io"
	"os"
)

func ExecuteGo(file string) string {
	originalStdout := os.Stdout
	originalStderr := os.Stderr
	defer func() {
		os.Stdout = originalStdout
		os.Stderr = originalStderr
	}()
	r, w, _ := os.Pipe()
	os.Stdout = w

	i := interp.New(interp.Options{})
	i.Use(stdlib.Symbols)
	i.Use(interp.Symbols)

	_, err := i.Eval(file)

	if err != nil {
		return err.Error()
	}

	// read stdout
	if err = w.Close(); err != nil {
		fmt.Print("fail")
	}
	stdout, err := io.ReadAll(r)
	if err != nil {
		fmt.Print("Error reading.")
	}

	return string(stdout)
}
