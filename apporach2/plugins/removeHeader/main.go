package main

import (
    "github.com/hashicorp/go-plugin"
    "example.com/mypkg/policy"
)

// RemoveHeaderPolicy removes a specific header if present.
type RemoveHeaderPolicy struct{}

func (p *RemoveHeaderPolicy) ProcessRequestHeaders(headers map[string]string) (map[string]string, error) {
    // remove header "X-To-Remove" if exists
    delete(headers, "X-To-Remove")
    return headers, nil
}

func main() {
    plugin.Serve(&plugin.ServeConfig{
        HandshakeConfig: policy.Handshake,
        Plugins: map[string]plugin.Plugin{
            "policy": &policy.PolicyPlugin{Impl: &RemoveHeaderPolicy{}},
        },
        GRPCServer: plugin.DefaultGRPCServer,
    })
}