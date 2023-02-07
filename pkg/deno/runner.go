package deno

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/creack/pty"

	"github.com/authgear/authgear-deno/pkg/ioutil"
)

type StdStream = *ioutil.LimitedWriter[*bytes.Buffer]

// StdStreamLimit is 1MiB.
const StdStreamLimit int64 = 1 * 1024 * 1024

type RunFileResult struct {
	Stdout StdStream
	Stderr StdStream
}

func (r *RunFileResult) Wrap(err error) error {
	return &RunFileError{
		Inner:  err,
		Stdout: r.Stdout,
		Stderr: r.Stderr,
	}
}

type RunFileError struct {
	Inner  error
	Stdout StdStream
	Stderr StdStream
}

func (e *RunFileError) Error() string {
	return e.Inner.Error()
}

func (e *RunFileError) Unwrap() error {
	return e.Inner
}

type RunFileOptions struct {
	// TargetScript is the filename of the target script.
	TargetScript string
	// Input is the filename of the input.
	Input string
	// Output is the filename of the output.
	Output string
	// Allow using unstable api in deno script if true
	IsUnstableAPIAllowed bool
}

type RunGoValueResult struct {
	Output interface{}
	Stdout StdStream
	Stderr StdStream
}

type RunGoValueOptions struct {
	// TargetScript is the content of the target script.
	TargetScript string
	// Input is the input.
	Input interface{}
	// Allow using unstable api in deno script if true
	IsUnstableAPIAllowed bool
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

func (r *Runner) RunFile(ctx context.Context, opts RunFileOptions) (*RunFileResult, error) {
	runnerScript := r.runnerScript()

	targetScript, err := filepath.Abs(opts.TargetScript)
	if err != nil {
		return nil, err
	}
	input, err := filepath.Abs(opts.Input)
	if err != nil {
		return nil, err
	}
	output, err := filepath.Abs(opts.Output)
	if err != nil {
		return nil, err
	}

	cmdArgs := []string{
		"run",
		"--quiet",
		fmt.Sprintf("--allow-read=%v,%v", targetScript, input),
		fmt.Sprintf("--allow-write=%v", output),
	}

	if opts.IsUnstableAPIAllowed {
		cmdArgs = append(cmdArgs, "--unstable")
	}

	cmdArgs = append(cmdArgs,
		runnerScript,
		targetScript,
		input,
		output)

	cmd := exec.CommandContext( //nolint:gosec
		ctx,
		"deno",
		cmdArgs...,
	)

	// Tell deno not to output ASCII escape code.
	cmd.Env = append(cmd.Environ(), "NO_COLOR=1")

	stdout := ioutil.LimitWriter(&bytes.Buffer{}, StdStreamLimit)
	stderr := ioutil.LimitWriter(&bytes.Buffer{}, StdStreamLimit)

	// Separate stdout and stderr.
	cmd.Stdout = stdout

	// Allocate a pty, connect stdin and stderr to the pty, and start the command.
	f, err := pty.Start(cmd)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Read stderr
	go func() {
		scanner := bufio.NewScanner(io.TeeReader(f, stderr))
		scanner.Split(ScanStderr)
		for scanner.Scan() {
			line := scanner.Text()
			// Start of permission prompt
			if strings.HasPrefix(line, "⚠️  ┌ Deno requests ") {
				var granted bool
				d, ok := LineToPermissionDescriptor(line)
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
		return nil, &RunFileError{
			Inner:  err,
			Stdout: stdout,
			Stderr: stderr,
		}
	}

	return &RunFileResult{
		Stdout: stdout,
		Stderr: stderr,
	}, nil
}

func (r *Runner) RunGoValue(ctx context.Context, opts RunGoValueOptions) (*RunGoValueResult, error) {
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

	runFileResult, err := r.RunFile(ctx, RunFileOptions{
		TargetScript:         targetScript.Name(),
		Input:                input.Name(),
		Output:               output.Name(),
		IsUnstableAPIAllowed: opts.IsUnstableAPIAllowed,
	})
	if err != nil {
		return nil, err
	}

	var out interface{}
	err = json.NewDecoder(output).Decode(&out)
	if err != nil {
		return nil, runFileResult.Wrap(err)
	}
	err = output.Close()
	if err != nil {
		return nil, runFileResult.Wrap(err)
	}

	return &RunGoValueResult{
		Output: out,
		Stdout: runFileResult.Stdout,
		Stderr: runFileResult.Stderr,
	}, nil
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
