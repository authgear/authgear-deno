package deno

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type CheckFileOptions struct {
	// TargetScript is the filename of the target script.
	TargetScript string
}

type CheckSnippetOptions struct {
	// TargetScript is the content of the target script.
	TargetScript string
}

type CheckFileError struct {
	Inner  error
	Stderr string
}

func (e *CheckFileError) Error() string {
	return fmt.Sprintf("%v\n%v", e.Inner, e.Stderr)
}

func (e *CheckFileError) Unwrap() error {
	return e.Inner
}

type Checker struct{}

func (c *Checker) CheckFile(ctx context.Context, opts CheckFileOptions) error {
	targetScript, err := filepath.Abs(opts.TargetScript)
	if err != nil {
		return err
	}

	cmd := exec.CommandContext( // #nosec G204
		ctx,
		"deno",
		"check",
		"--quiet",
		targetScript,
	)

	// Tell deno not to output ASCII escape code.
	cmd.Env = append(cmd.Environ(), "NO_COLOR=1")

	stderr := &bytes.Buffer{}
	cmd.Stderr = stderr

	err = cmd.Run()
	if err != nil {
		return &CheckFileError{
			Inner:  err,
			Stderr: c.fixStderr(targetScript, stderr),
		}
	}

	return nil
}

func (c *Checker) CheckSnippet(ctx context.Context, opts CheckSnippetOptions) error {
	targetScript, err := os.CreateTemp("", "authgear-deno-script.*.ts")
	if err != nil {
		return err
	}
	defer os.Remove(targetScript.Name())

	_, err = io.Copy(targetScript, strings.NewReader(opts.TargetScript))
	if err != nil {
		return err
	}
	err = targetScript.Close()
	if err != nil {
		return err
	}

	err = c.CheckFile(ctx, CheckFileOptions{
		TargetScript: targetScript.Name(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Checker) fixStderr(targetScript string, stderr *bytes.Buffer) string {
	u := &url.URL{
		Scheme: "file",
		Path:   targetScript,
	}
	return strings.ReplaceAll(stderr.String(), u.String(), "FILE")
}
