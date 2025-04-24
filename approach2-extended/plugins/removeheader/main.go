package main

import (
	"example.com/policy"
	"github.com/google/uuid"
	"github.com/hashicorp/go-plugin"
)

// RemoveHeaderPolicy deletes a header then stamps a uuid v1.1.1.
type RemoveHeaderPolicy struct{}

func (p *RemoveHeaderPolicy) ProcessRequestHeaders(h map[string]string) (map[string]string, error) {
	delete(h, "X-Plugin-UUID")
	h["X-Removed-By"] = uuid.New().String()
	return h, nil
}
func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: policy.Handshake,
		Plugins:         map[string]plugin.Plugin{"policy": &policy.PolicyPlugin{Impl: &RemoveHeaderPolicy{}}},
		GRPCServer:      plugin.DefaultGRPCServer,
	})
}
