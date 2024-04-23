package models

import (
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"io"
	"log"
	"os"
)

type ExecutionerModel struct {
}

func (em ExecutionerModel) Execute(code string, output chan string) {
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

	_, err := i.Eval(code)

	if err != nil {
		output <- err.Error()
	}

	// read stdout
	if err = w.Close(); err != nil {
		log.Print("Failing to close write pipe")
	}
	stdout, err := io.ReadAll(r)
	if err != nil {
		log.Print("Error reading from pipe")
	}

	output <- string(stdout)
}
