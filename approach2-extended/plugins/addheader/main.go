package main

import (
	"example.com/policy"
	"github.com/google/uuid"
	"github.com/hashicorp/go-plugin"
)

// AddHeaderPolicy uses uuid v1.3.0 to stamp headers.
type AddHeaderPolicy struct{}

func (p *AddHeaderPolicy) ProcessRequestHeaders(h map[string]string) (map[string]string, error) {
	h["X-Plugin-UUID"] = uuid.New().String()
	return h, nil
}
func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: policy.Handshake,
		Plugins:         map[string]plugin.Plugin{"policy": &policy.PolicyPlugin{Impl: &AddHeaderPolicy{}}},
		GRPCServer:      plugin.DefaultGRPCServer,
	})
}
