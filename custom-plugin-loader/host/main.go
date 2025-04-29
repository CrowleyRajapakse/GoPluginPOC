package main

import (
	"fmt"
	"io"
	"log"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/uuid" // v1.1.1
)

type stdioRW struct {
	in  io.Reader
	out io.Writer
}

func (s *stdioRW) Read(p []byte) (int, error)  { return s.in.Read(p) }
func (s *stdioRW) Write(p []byte) (int, error) { return s.out.Write(p) }
func (s *stdioRW) Close() error                { return nil }

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

func loadPlugin(path string) (*rpc.Client, *exec.Cmd) {
	cmd := exec.Command(path)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalf("stdin pipe: %v", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("stdout pipe: %v", err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatalf("start %s: %v", path, err)
	}
	rw := &stdioRW{in: stdout, out: stdin}
	client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(rw))
	return client, cmd
}

func main() {
	pluginDir := "./plugins"
	entries, err := os.ReadDir(pluginDir)
	if err != nil {
		log.Fatalf("read dir: %v", err)
	}

	// initial headers with host-generated UUID
	headers := Headers{"RequestID": uuid.New().String()}

	var clients []*rpc.Client
	var cmds []*exec.Cmd

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		path := filepath.Join(pluginDir, e.Name())
		info, _ := os.Stat(path)
		if info.Mode()&0111 == 0 {
			continue
		}
		client, cmd := loadPlugin(path)
		clients = append(clients, client)
		cmds = append(cmds, cmd)

		switch {
		case strings.Contains(e.Name(), "addheader"):
			var out Headers
			args := AddArgs{Key: "X-Trace-Id", Value: "", Headers: headers}
			if err := client.Call("Plugin.Add", args, &out); err != nil {
				log.Printf("Add error: %v", err)
			} else {
				fmt.Printf("[%s] after Add: %#v\n", e.Name(), out)
				headers = out
			}

		case strings.Contains(e.Name(), "removeheader"):
			var out Headers
			args := RemoveArgs{Key: "RequestID", Headers: headers}
			if err := client.Call("Plugin.Remove", args, &out); err != nil {
				log.Printf("Remove error: %v", err)
			} else {
				fmt.Printf("[%s] after Remove: %#v\n", e.Name(), out)
				headers = out
			}

		default:
			log.Printf("skipping %q", e.Name())
		}
	}

	for _, c := range clients {
		c.Close()
	}
	for _, c := range cmds {
		c.Wait()
	}
}
