package main

import (
	"fmt"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"

	"github.com/google/uuid" // v1.2.0
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
	Headers Headers
}

type Plugin struct{}

func (p *Plugin) Remove(args Args, reply *Headers) error {
	h := args.Headers
	delete(h, args.Key)
	// for demo: log removal with a UUID
	fmt.Fprintf(os.Stderr, "removeheader[%s]: removed %q at %s\n",
		uuid.New().String(), args.Key, os.Args[0])
	*reply = h
	return nil
}

func main() {
	rpc.RegisterName("Plugin", new(Plugin))
	stdio := &stdioRW{os.Stdin, os.Stdout}
	rpc.ServeCodec(jsonrpc.NewServerCodec(stdio))
}
