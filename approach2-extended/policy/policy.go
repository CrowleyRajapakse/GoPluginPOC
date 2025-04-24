package policy

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

// Policy is the common interface for header-processing plugins.
type Policy interface {
	ProcessRequestHeaders(map[string]string) (map[string]string, error)
}

// PolicyPlugin wraps a Policy for go-plugin.
type PolicyPlugin struct {
	Impl Policy
	plugin.NetRPCUnsupportedPlugin
}

func (p *PolicyPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &RPCServer{Impl: p.Impl}, nil
}

func (PolicyPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &RPCClient{client: c}, nil
}

// RPCServer bridges RPC to the real implementation.
type RPCServer struct{ Impl Policy }

func (s *RPCServer) ProcessRequestHeaders(args map[string]string, resp *map[string]string) error {
	out, err := s.Impl.ProcessRequestHeaders(args)
	if err != nil {
		return err
	}
	*resp = out
	return nil
}

// RPCClient calls RPC methods on the plugin.
type RPCClient struct{ client *rpc.Client }

func (c *RPCClient) ProcessRequestHeaders(headers map[string]string) (map[string]string, error) {
	var out map[string]string
	err := c.client.Call("Plugin.ProcessRequestHeaders", headers, &out)
	return out, err
}

// Handshake config for host/plugin compatibility.
var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "POLICY_PLUGIN",
	MagicCookieValue: "policy",
}
