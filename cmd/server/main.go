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

	runHandler := handler.NewRunner(&deno.Runner{
		Permissioner: deno.DisallowIPPolicy(cfg.IPPolicies()...),
	}, cfg.RunMaxConcurrency)
	http.Handle("/run", runHandler)
	http.Handle("/check", &handler.Checker{
		Checker: &deno.Checker{},
	})

	server := &http.Server{
		Addr:              cfg.ListenAddr,
		ReadHeaderTimeout: 3 * time.Second,
	}

	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
