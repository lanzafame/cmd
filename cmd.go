package cmd

import (
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var lineBreak = "\n"

// Command is a structure representing a shell command to be run in the
// specified directory.
type Command struct {
	dir   string
	parts []string
}

// New creates a new command using the pwd and its cwd.
func New(parts ...string) Command {
	return NewWithDir("./", parts...)
}

// NewWithDir creates a new command using the specified directory as its cwd.
func NewWithDir(dir string, parts ...string) Command {
	return Command{
		dir:   dir,
		parts: parts,
	}
}

// RunCmd runs a Command.
func RunCmd(c Command) {
	parts := c.parts
	if len(parts) == 1 {
		parts = strings.Split(parts[0], " ")
	}

	name := strings.Join(parts, " ")
	cmd := exec.Command(parts[0], parts[1:]...) // #nosec
	cmd.Dir = c.dir
	log.Println(name)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		if _, err = io.Copy(os.Stderr, stderr); err != nil {
			panic(err)
		}
	}()
	go func() {
		defer wg.Done()
		if _, err = io.Copy(os.Stdout, stdout); err != nil {
			panic(err)
		}
	}()

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	wg.Wait()
	if err := cmd.Wait(); err != nil {
		log.Fatalf("Command '%s' failed: %s\n", name, err)
	}
}

// RunCapture runs a command and captures the output and
// returns it as a string.
func RunCapture(name string) string {
	args := strings.Split(name, " ")
	cmd := exec.Command(args[0], args[1:]...) // #nosec
	log.Println(name)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Command '%s' failed: %s\n", name, err)
	}

	return strings.Trim(string(output), lineBreak)
}
