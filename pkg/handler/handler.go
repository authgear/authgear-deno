package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/authgear/authgear-deno/pkg/deno"
)

type RunRequest struct {
	Script string      `json:"script"`
	Input  interface{} `json:"input"`
}

type RunResponse struct {
	Error  string      `json:"error,omitempty"`
	Output interface{} `json:"output,omitempty"`
	Stderr string      `json:"stderr,omitempty"`
	Stdout string      `json:"stdout,omitempty"`
}

type T struct {
	Runner *deno.Runner
}

func New(runner *deno.Runner) *T {
	return &T{
		Runner: runner,
	}
}

func (t *T) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	result, err := t.handle(w, r)
	if err != nil {
		t.writeError(w, r, err)
		return
	}
	t.writeResult(w, r, result)
}

func (t *T) handle(_ http.ResponseWriter, r *http.Request) (*deno.RunGoValueResult, error) {
	var runRequest RunRequest
	err := json.NewDecoder(r.Body).Decode(&runRequest)
	if err != nil {
		return nil, err
	}

	result, err := t.Runner.RunGoValue(r.Context(), deno.RunGoValueOptions{
		TargetScript: runRequest.Script,
		Input:        runRequest.Input,
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (t *T) writeError(w http.ResponseWriter, r *http.Request, err error) {
	runResponse := RunResponse{
		Error: err.Error(),
	}

	var runFileError *deno.RunFileError
	if errors.As(err, &runFileError) {
		runResponse.Stderr = runFileError.Stderr.String()
		runResponse.Stdout = runFileError.Stdout.String()
	}

	t.writeJSON(w, r, runResponse)

}

func (t *T) writeResult(w http.ResponseWriter, r *http.Request, result *deno.RunGoValueResult) {
	runResponse := RunResponse{
		Output: result.Output,
		Stderr: result.Stderr.String(),
		Stdout: result.Stdout.String(),
	}
	t.writeJSON(w, r, runResponse)
}

func (t *T) writeJSON(w http.ResponseWriter, _ *http.Request, jsonValue interface{}) {
	w.Header().Set("Content-Type", "application/json")
	//nolint:errchkjson
	_ = json.NewEncoder(w).Encode(jsonValue)
}
