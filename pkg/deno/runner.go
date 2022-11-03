package deno

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/creack/pty"
)

type Runner struct {
	Permissioner Permissioner
}

func NewRunner(permissioner Permissioner) *Runner {
	return &Runner{
		Permissioner: permissioner,
	}
}

func (r *Runner) RunFile(ctx context.Context, filename string) error {
	cmd := exec.CommandContext(
		ctx,
		"deno",
		"run",
		"--quiet",
		filename,
	)

	// Tell deno not to output ASCII escape code.
	cmd.Env = append(cmd.Environ(), "NO_COLOR=1")

	// Pipe stdout to io.Discard so that we get a cleaner stderr.
	cmd.Stdout = io.Discard

	// Allocate a pty, connect stdin and stderr to the pty, and start the command.
	f, err := pty.Start(cmd)
	if err != nil {
		return err
	}
	defer f.Close()

	// Read stderr
	go func() {
		scanner := bufio.NewScanner(f)
		scanner.Split(ScanStderr)
		for scanner.Scan() {
			line := scanner.Text()
			// Start of permission prompt
			if strings.HasPrefix(line, "⚠️  ┌ Deno requests ") {
				var granted bool
				d, ok := lineToPermissionDescriptor(line)
				if ok {
					granted = r.askPermission(ctx, *d)
				}
				r.skipToPrompt(scanner)
				if granted {
					fmt.Fprintf(f, "y\n")
				} else {
					fmt.Fprintf(f, "n\n")
				}
			}
		}
	}()

	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}

func (r *Runner) askPermission(ctx context.Context, d PermissionDescriptor) bool {
	if r.Permissioner == nil {
		return false
	}
	ok, err := r.Permissioner.RequestPermission(ctx, d)
	if err != nil {
		return false
	}
	return ok
}

func (r *Runner) skipToPrompt(scanner *bufio.Scanner) {
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "(y = yes, allow; n = no, deny) > ") {
			break
		}
	}
}
