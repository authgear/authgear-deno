package deno

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/creack/pty"
)

type RunOptions struct {
	// TargetScript is the filename of the target script.
	TargetScript string
	// Input is the filename of the input.
	Input string
	// Output is the filename of the output.
	Output string
}

type Runner struct {
	// RunnerScript is the runner script that will import the target script
	// and execute the default function.
	RunnerScript string
	// Permissioner manages the permissions of the target script.
	Permissioner Permissioner
}

func (r *Runner) runnerScript() string {
	if r.RunnerScript != "" {
		return r.RunnerScript
	}
	return "./runner.ts"
}

func (r *Runner) RunFile(ctx context.Context, opts RunOptions) error {
	runnerScript := r.runnerScript()

	targetScript, err := filepath.Abs(opts.TargetScript)
	if err != nil {
		return err
	}
	input, err := filepath.Abs(opts.Input)
	if err != nil {
		return err
	}
	output, err := filepath.Abs(opts.Output)
	if err != nil {
		return err
	}

	cmd := exec.CommandContext(
		ctx,
		"deno",
		"run",
		"--quiet",
		fmt.Sprintf("--allow-read=%v,%v", targetScript, input),
		fmt.Sprintf("--allow-write=%v", output),
		runnerScript,
		targetScript,
		input,
		output,
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
