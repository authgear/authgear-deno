package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/authgear/authgear-deno/pkg/deno"
)

type RunRequest struct {
	Script string      `json:"script"`
	Input  interface{} `json:"input"`
}

type Stream struct {
	String    string `json:"string,omitempty"`
	Truncated bool   `json:"truncated,omitempty"`
}

func NewStream(stdStream deno.StdStream) *Stream {
	return &Stream{
		String:    stdStream.W.String(),
		Truncated: stdStream.Exceeded,
	}
}

type ErrorCode string

const (
	ErrorCodeRunTimout ErrorCode = "run_timeout"
	ErrorCodeUnknown   ErrorCode = "unknown"
)

type RunResponse struct {
	Error     string      `json:"error,omitempty"`
	ErrorCode ErrorCode   `json:"error_code,omitempty"`
	Output    interface{} `json:"output,omitempty"`
	Stderr    *Stream     `json:"stderr,omitempty"`
	Stdout    *Stream     `json:"stdout,omitempty"`
}

type Runner struct {
	Runner         *deno.Runner
	sema           chan struct{}
	timeoutSeconds int
}

func NewRunner(runner *deno.Runner, maxConcurrency int, timeoutSeconds int) *Runner {
	return &Runner{
		Runner:         runner,
		sema:           make(chan struct{}, maxConcurrency),
		timeoutSeconds: timeoutSeconds,
	}
}

func (t *Runner) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Wait for a slot to be available or for the request context to be done.
	select {
	case t.sema <- struct{}{}:
		// Acquired a slot; ensure release after handling.
		defer func() { <-t.sema }()
	case <-r.Context().Done():
		http.Error(w, "request canceled", http.StatusRequestTimeout)
		return
	}

	result, err := t.handle(w, r)
	if err != nil {
		t.writeError(w, r, err)
		return
	}
	t.writeResult(w, r, result)
}

func (t *Runner) handle(_ http.ResponseWriter, r *http.Request) (*deno.RunGoValueResult, error) {
	var runRequest RunRequest
	err := json.NewDecoder(r.Body).Decode(&runRequest)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(t.timeoutSeconds)*time.Second)
	defer cancel()

	result, err := t.Runner.RunGoValue(ctx, deno.RunGoValueOptions{
		TargetScript: runRequest.Script,
		Input:        runRequest.Input,
	})
	if err != nil {
		return nil, errors.Join(err, ctx.Err())
	}

	return result, nil
}

func (t *Runner) writeError(w http.ResponseWriter, r *http.Request, err error) {
	runResponse := RunResponse{
		Error: err.Error(),
	}

	var runFileError *deno.RunFileError
	if errors.As(err, &runFileError) {
		runResponse.Stderr = NewStream(runFileError.Stderr)
		runResponse.Stdout = NewStream(runFileError.Stdout)
	}
	if errors.Is(err, context.DeadlineExceeded) {
		runResponse.ErrorCode = ErrorCodeRunTimout
	} else {
		runResponse.ErrorCode = ErrorCodeUnknown
	}

	writeJSON(w, r, runResponse)
}

func (t *Runner) writeResult(w http.ResponseWriter, r *http.Request, result *deno.RunGoValueResult) {
	runResponse := RunResponse{
		Output: result.Output,
		Stderr: NewStream(result.Stderr),
		Stdout: NewStream(result.Stdout),
	}
	writeJSON(w, r, runResponse)
}

func writeJSON(w http.ResponseWriter, _ *http.Request, jsonValue interface{}) {
	w.Header().Set("Content-Type", "application/json")
	//nolint:errchkjson
	_ = json.NewEncoder(w).Encode(jsonValue)
}
