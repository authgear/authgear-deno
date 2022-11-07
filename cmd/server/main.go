package main

import (
	"net/http"
	"time"

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

	server := &http.Server{
		Addr:              cfg.ListenAddr,
		ReadHeaderTimeout: 3 * time.Second,
	}

	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
