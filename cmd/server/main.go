package main

import (
	"net/http"

	"github.com/authgear/authgear-deno/pkg/deno"
	"github.com/authgear/authgear-deno/pkg/handler"
)

func main() {
	cfg, err := LoadConfigFromEnv()
	if err != nil {
		panic(err)
	}

	permissioner := deno.DisallowIPPolicy(cfg.IPPolicies()...)
	runner := &deno.Runner{
		RunnerScript: cfg.RunnerScript,
		Permissioner: permissioner,
	}
	h := handler.New(runner)
	http.Handle("/", h)
	err = http.ListenAndServe(cfg.ListenAddr, nil)
	if err != nil {
		panic(err)
	}
}
