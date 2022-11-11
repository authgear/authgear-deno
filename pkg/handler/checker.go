package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/authgear/authgear-deno/pkg/deno"
)

type CheckRequest struct {
	Script string `json:"script"`
}

type CheckResponse struct {
	Stderr string `json:"stderr,omitempty"`
}

type Checker struct {
	Checker *deno.Checker
}

func (t *Checker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	err := t.handle(w, r)
	if err != nil {
		t.writeError(w, r, err)
		return
	}
	t.writeResult(w, r)
}

func (t *Checker) handle(_ http.ResponseWriter, r *http.Request) error {
	var checkRequest CheckRequest
	err := json.NewDecoder(r.Body).Decode(&checkRequest)
	if err != nil {
		return err
	}

	err = t.Checker.CheckSnippet(r.Context(), deno.CheckSnippetOptions{
		TargetScript: checkRequest.Script,
	})
	if err != nil {
		return err
	}

	return nil
}

func (t *Checker) writeError(w http.ResponseWriter, r *http.Request, err error) {
	checkResponse := CheckResponse{}

	var checkError *deno.CheckFileError
	if errors.As(err, &checkError) {
		checkResponse.Stderr = checkError.Stderr
	}

	writeJSON(w, r, checkResponse)
}

func (t *Checker) writeResult(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, r, CheckResponse{})
}
