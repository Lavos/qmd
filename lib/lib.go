package lib

import (
	"strings"
	"fmt"
	"os"
	"os/exec"
	"path"
)

type Qmd struct {
	Cmd *exec.Cmd

	redirectFileLocation string
	outputFileLocation string

	redirectFileHandle *os.File
}

type Namespace struct {
	Prefix string
}

type V [2]string

func NewNamespace(prefix string) *Namespace {
	return &Namespace{ prefix }
}

func (n *Namespace) E(key, value string, args ...interface{}) V {
	return E(fmt.Sprintf("%s_%s", n.Prefix, key), value, args...)
}

func E(key, value string, args ...interface{}) V {
	return V{ key, fmt.Sprintf(value, args...) }
}

func NewQmd(name string, args ...string) *Qmd {
	c := exec.Command(name, args...)

	c.Stderr = os.Stderr

	return &Qmd{
		Cmd: c,
	}
}

func (q *Qmd) RedirectFile(fileLocation string) *Qmd {
	q.redirectFileLocation = fileLocation
	return q
}

func (q *Qmd) OutputFile(fileLocation string) *Qmd {
	q.outputFileLocation = fileLocation
	return q
}

func (q *Qmd) AppendEnv(vs ...V) *Qmd {
	for _, v := range vs {
		q.Cmd.Env = append(q.Cmd.Env, strings.Join(v[:], "="))
	}

	return q
}

func (q *Qmd) PipeToCommand(dcmd *exec.Cmd) error {
	// get stdin pipe from Cmd
	stdin, err := dcmd.StdinPipe()

	if err != nil {
		return err
	}

	q.Cmd.Stdout = stdin
	return nil
}

func (q *Qmd) Start() error {
	// check redirectFileLocation
	if q.redirectFileLocation != "" {
		redirectFile, err := os.Open(q.redirectFileLocation)

		if err != nil {
			return err
		}

		q.Cmd.Stdin = redirectFile
		q.redirectFileHandle = redirectFile
	}

	// check outputFileLocation
	if q.outputFileLocation != "" {
		// ensure the target directory exists
		dir := path.Dir(q.outputFileLocation)

		if dir != "." {
			err := os.MkdirAll(dir, 0700)

			if err != nil {
				return err
			}
		}

		outputFile, err := os.Create(q.outputFileLocation)

		if err != nil {
			return err
		}

		q.Cmd.Stdout = outputFile
	}

	return q.Cmd.Start()
}

func (q *Qmd) Wait() error {
	if q.redirectFileHandle != nil {
		defer q.redirectFileHandle.Close()
	}

	return q.Cmd.Wait()
}
