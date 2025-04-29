package main

import (
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"

	"github.com/google/uuid" // v1.3.0
)

type stdioRW struct {
	in  *os.File
	out *os.File
}

func (s *stdioRW) Read(p []byte) (n int, err error)  { return s.in.Read(p) }
func (s *stdioRW) Write(p []byte) (n int, err error) { return s.out.Write(p) }
func (s *stdioRW) Close() error                      { return nil }

type Headers map[string]string

type Args struct {
	Key     string
	Value   string
	Headers Headers
}

type Plugin struct{}

func (p *Plugin) Add(args Args, reply *Headers) error {
	h := args.Headers
	if h == nil {
		h = make(Headers)
	}
	// if no value provided, generate a UUID
	if args.Value == "" {
		args.Value = uuid.New().String()
	}
	h[args.Key] = args.Value
	*reply = h
	return nil
}

func main() {
	rpc.RegisterName("Plugin", new(Plugin))
	stdio := &stdioRW{os.Stdin, os.Stdout}
	rpc.ServeCodec(jsonrpc.NewServerCodec(stdio))
}
