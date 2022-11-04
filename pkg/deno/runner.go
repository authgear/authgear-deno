package deno

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/creack/pty"
)

type RunFileOptions struct {
	// TargetScript is the filename of the target script.
	TargetScript string
	// Input is the filename of the input.
	Input string
	// Output is the filename of the output.
	Output string
}

type RunGoValueOptions struct {
	// TargetScript is the content of the target script.
	TargetScript string
	// Input is the input.
	Input interface{}
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

func (r *Runner) RunFile(ctx context.Context, opts RunFileOptions) error {
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

func (r *Runner) RunGoValue(ctx context.Context, opts RunGoValueOptions) (out interface{}, err error) {
	targetScript, err := os.CreateTemp("", "authgear-deno-script.*.ts")
	if err != nil {
		return nil, err
	}
	defer os.Remove(targetScript.Name())

	input, err := os.CreateTemp("", "authgear-deno-input.*.json")
	if err != nil {
		return nil, err
	}
	defer os.Remove(input.Name())

	output, err := os.CreateTemp("", "authgear-deno-output.*.json")
	if err != nil {
		return nil, err
	}
	defer os.Remove(output.Name())

	_, err = io.Copy(targetScript, strings.NewReader(opts.TargetScript))
	if err != nil {
		return nil, err
	}
	err = targetScript.Close()
	if err != nil {
		return nil, err
	}

	err = json.NewEncoder(input).Encode(opts.Input)
	if err != nil {
		return nil, err
	}
	err = input.Close()
	if err != nil {
		return nil, err
	}

	err = r.RunFile(ctx, RunFileOptions{
		TargetScript: targetScript.Name(),
		Input:        input.Name(),
		Output:       output.Name(),
	})
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(output).Decode(&out)
	if err != nil {
		return nil, err
	}
	err = output.Close()
	if err != nil {
		return nil, err
	}

	return
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
