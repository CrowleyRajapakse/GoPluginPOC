package main

import (
	"example.com/mypkg/policy"
	"github.com/hashicorp/go-plugin"
)

// MyPolicy is your concrete Policy implementation.
type MyPolicy struct{}

func (p *MyPolicy) ProcessRequestHeaders(headers map[string]string) (map[string]string, error) {
	headers["X-Added-By"] = "MyPolicy"
	return headers, nil
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: policy.Handshake,
		Plugins: map[string]plugin.Plugin{
			"policy": &policy.PolicyPlugin{Impl: &MyPolicy{}},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
