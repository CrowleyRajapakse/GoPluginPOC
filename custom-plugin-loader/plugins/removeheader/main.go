package main

import (
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
)

type stdioReadWriteCloser struct {
	in  *os.File
	out *os.File
}

func (s *stdioReadWriteCloser) Read(p []byte) (n int, err error)  { return s.in.Read(p) }
func (s *stdioReadWriteCloser) Write(p []byte) (n int, err error) { return s.out.Write(p) }
func (s *stdioReadWriteCloser) Close() error                      { return nil }

type Headers map[string]string

type Args struct {
	Key     string
	Headers Headers
}

type Plugin struct{}

func (p *Plugin) Remove(args Args, reply *Headers) error {
	h := args.Headers
	delete(h, args.Key)
	*reply = h
	return nil
}

func main() {
	rpc.RegisterName("Plugin", new(Plugin))
	stdio := &stdioReadWriteCloser{os.Stdin, os.Stdout}

	// <-- CORRECTED here as well
	rpc.ServeCodec(jsonrpc.NewServerCodec(stdio))
}
