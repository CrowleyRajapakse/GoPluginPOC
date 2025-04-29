package main

import (
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
)

// stdioReadWriteCloser lets us treat stdin/stdout as one ReadWriteCloser
type stdioReadWriteCloser struct {
	in  *os.File
	out *os.File
}

func (s *stdioReadWriteCloser) Read(p []byte) (n int, err error)  { return s.in.Read(p) }
func (s *stdioReadWriteCloser) Write(p []byte) (n int, err error) { return s.out.Write(p) }
func (s *stdioReadWriteCloser) Close() error                      { return nil }

// Headers is just a map of string→string
type Headers map[string]string

// Args passed in for adding a header
type Args struct {
	Key     string
	Value   string
	Headers Headers
}

// Plugin exposes an Add method
type Plugin struct{}

func (p *Plugin) Add(args Args, reply *Headers) error {
	h := args.Headers
	if h == nil {
		h = make(Headers)
	}
	h[args.Key] = args.Value
	*reply = h
	return nil
}

func main() {
	rpc.RegisterName("Plugin", new(Plugin))
	stdio := &stdioReadWriteCloser{os.Stdin, os.Stdout}

	// <-- CORRECTED: ServeCodec is in net/rpc, not jsonrpc
	rpc.ServeCodec(jsonrpc.NewServerCodec(stdio))
}
