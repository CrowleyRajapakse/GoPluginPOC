package main

import (
	"example.com/mypkg/policy"
	"github.com/hashicorp/go-plugin"
)

// AddHeaderPolicy adds a custom header.
type AddHeaderPolicy struct{}

func (p *AddHeaderPolicy) ProcessRequestHeaders(headers map[string]string) (map[string]string, error) {
	headers["X-Added-By"] = "AddHeaderPolicy"
	return headers, nil
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: policy.Handshake,
		Plugins: map[string]plugin.Plugin{
			"policy": &policy.PolicyPlugin{Impl: &AddHeaderPolicy{}},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
