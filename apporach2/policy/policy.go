package policy

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

// Policy defines the interface for header-processing plugins.
type Policy interface {
	ProcessRequestHeaders(headers map[string]string) (map[string]string, error)
}

// PolicyPlugin is the go-plugin wrapper for the Policy interface.
type PolicyPlugin struct {
	Impl Policy
	plugin.NetRPCUnsupportedPlugin
}

// Server is called when the plugin is serving.
func (p *PolicyPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &RPCServer{Impl: p.Impl}, nil
}

// Client is called when the host is dispending the plugin. Returns an RPC client with error.
func (PolicyPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &RPCClient{client: c}, nil
}

// RPCServer bridges RPC calls to the Policy implementation.
type RPCServer struct {
	Impl Policy
}

func (s *RPCServer) ProcessRequestHeaders(args map[string]string, resp *map[string]string) error {
	result, err := s.Impl.ProcessRequestHeaders(args)
	if err != nil {
		return err
	}
	*resp = result
	return nil
}

// RPCClient makes RPC calls to the server.
type RPCClient struct {
	client *rpc.Client
}

func (c *RPCClient) ProcessRequestHeaders(headers map[string]string) (map[string]string, error) {
	var resp map[string]string
	err := c.client.Call("Plugin.ProcessRequestHeaders", headers, &resp)
	return resp, err
}

// Handshake ensures that the host and plugin are compatible.
var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "POLICY_PLUGIN",
	MagicCookieValue: "policy",
}
