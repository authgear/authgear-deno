package main

import (
	"github.com/kelseyhightower/envconfig"

	"github.com/authgear/authgear-deno/pkg/deno"
)

type Config struct {
	ListenAddr                      string `envconfig:"LISTEN_ADDR" default:"0.0.0.0:8090"`
	RunnerScript                    string `envconfig:"RUNNER_SCRIPT" default:"./pkg/deno/runner.ts"`
	DisallowGlobalUnicast           bool   `envconfig:"DISALLOW_GLOBAL_UNICAST" default:"false"`
	DisallowInterfaceLocalMulticast bool   `envconfig:"DISALLOW_INTERFACE_LOCAL_MULTICAST" default:"true"`
	DisallowLinkLocalUnicast        bool   `envconfig:"DISALLOW_LINK_LOCAL_UNICAST" default:"true"`
	DisallowLinkLocalMulticast      bool   `envconfig:"DISALLOW_LINK_LOCAL_MULTICAST" default:"true"`
	DisallowLoopback                bool   `envconfig:"DISALLOW_LOOPBACK" default:"true"`
	DisallowMulticast               bool   `envconfig:"DISALLOW_MULTICAST" default:"true"`
	DisallowPrivate                 bool   `envconfig:"DISALLOW_PRIVATE" default:"true"`
	DisallowUnspecified             bool   `envconfig:"DISALLOW_UNSPECIFIED" default:"true"`
}

func LoadConfigFromEnv() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) IPPolicies() []deno.IPPolicy {
	var policies []deno.IPPolicy

	if c.DisallowGlobalUnicast {
		policies = append(policies, deno.DisallowGlobalUnicast)
	}
	if c.DisallowInterfaceLocalMulticast {
		policies = append(policies, deno.DisallowInterfaceLocalMulticast)
	}
	if c.DisallowLinkLocalUnicast {
		policies = append(policies, deno.DisallowLinkLocalUnicast)
	}
	if c.DisallowLinkLocalMulticast {
		policies = append(policies, deno.DisallowLinkLocalMulticast)
	}
	if c.DisallowLoopback {
		policies = append(policies, deno.DisallowLoopback)
	}
	if c.DisallowMulticast {
		policies = append(policies, deno.DisallowMulticast)
	}
	if c.DisallowPrivate {
		policies = append(policies, deno.DisallowPrivate)
	}
	if c.DisallowUnspecified {
		policies = append(policies, deno.DisallowUnspecified)
	}

	return policies
}
