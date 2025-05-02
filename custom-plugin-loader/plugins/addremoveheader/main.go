package main

import (
	"fmt"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"

	"github.com/google/uuid"
)

type stdioRW struct{ in, out *os.File }

func (s *stdioRW) Read(p []byte) (n int, err error)  { return s.in.Read(p) }
func (s *stdioRW) Write(p []byte) (n int, err error) { return s.out.Write(p) }
func (s *stdioRW) Close() error                      { return nil }

type Headers map[string]string

type AddArgs struct {
	Key     string
	Value   string
	Headers Headers
}

type RemoveArgs struct {
	Key     string
	Headers Headers
}

type Plugin struct{}

func (p *Plugin) Add(args AddArgs, reply *Headers) error {
	h := args.Headers
	if h == nil {
		h = make(Headers)
	}
	if args.Value == "" {
		// generate a UUID when no value provided
		args.Value = uuid.New().String()
	}
	h[args.Key] = args.Value
	*reply = h
	return nil
}

func (p *Plugin) Remove(args RemoveArgs, reply *Headers) error {
	h := args.Headers
	delete(h, args.Key)
	// log removal with a UUID for traceability
	fmt.Fprintf(os.Stderr, "[addremoveheader] removed '%s' with id %s", args.Key, uuid.New().String())
	*reply = h
	return nil
}

func main() {
	rpc.RegisterName("Plugin", new(Plugin))
	stdio := &stdioRW{os.Stdin, os.Stdout}
	rpc.ServeCodec(jsonrpc.NewServerCodec(stdio))
}
