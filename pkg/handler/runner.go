package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/authgear/authgear-deno/pkg/deno"
)

type RunRequest struct {
	Script               string      `json:"script"`
	Input                interface{} `json:"input"`
	IsUnstableAPIAllowed bool        `json:"is_unstable_api_allowed"`
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

type RunResponse struct {
	Error  string      `json:"error,omitempty"`
	Output interface{} `json:"output,omitempty"`
	Stderr *Stream     `json:"stderr,omitempty"`
	Stdout *Stream     `json:"stdout,omitempty"`
}

type Runner struct {
	Runner *deno.Runner
}

func (t *Runner) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
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

	result, err := t.Runner.RunGoValue(r.Context(), deno.RunGoValueOptions{
		TargetScript:         runRequest.Script,
		Input:                runRequest.Input,
		IsUnstableAPIAllowed: runRequest.IsUnstableAPIAllowed,
	})
	if err != nil {
		return nil, err
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
