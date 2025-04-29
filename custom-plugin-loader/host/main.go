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

// loadPlugin starts the binary at path and returns an RPC client + the cmd handle
func loadPlugin(path string) (*rpc.Client, *exec.Cmd) {
	cmd := exec.Command(path)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalf("error opening stdin pipe for %s: %v", path, err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("error opening stdout pipe for %s: %v", path, err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatalf("error starting plugin %s: %v", path, err)
	}
	rw := &stdioRW{in: stdout, out: stdin}
	client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(rw))
	return client, cmd
}

func main() {
	pluginDir := "./plugins"

	entries, err := os.ReadDir(pluginDir)
	if err != nil {
		log.Fatalf("failed to read plugin dir %q: %v", pluginDir, err)
	}

	// keep track of clients + cmds so we can clean up
	var (
		clients   []*rpc.Client
		cmds      []*exec.Cmd
		pluginIDs []string
	)

	// 1) scan & load every executable in pluginDir
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		fullPath := filepath.Join(pluginDir, e.Name())
		// optional: skip non-executable files
		if info, err := os.Stat(fullPath); err != nil || info.Mode()&0111 == 0 {
			continue
		}

		client, cmd := loadPlugin(fullPath)
		clients = append(clients, client)
		cmds = append(cmds, cmd)
		pluginIDs = append(pluginIDs, e.Name())
	}

	// 2) demo: for each plugin, decide which RPC to call by filename
	initial := Headers{"User-Agent": "golang"}
	for i, name := range pluginIDs {
		client := clients[i]

		switch {
		case strings.Contains(name, "addheader"):
			var out Headers
			args := AddArgs{Key: "X-Trace-Id", Value: "abc123", Headers: initial}
			if err := client.Call("Plugin.Add", args, &out); err != nil {
				log.Printf("[%s] Add error: %v", name, err)
			} else {
				fmt.Printf("[%s] after Add: %#v\n", name, out)
				initial = out
			}

		case strings.Contains(name, "removeheader"):
			var out Headers
			args := RemoveArgs{Key: "User-Agent", Headers: initial}
			if err := client.Call("Plugin.Remove", args, &out); err != nil {
				log.Printf("[%s] Remove error: %v", name, err)
			} else {
				fmt.Printf("[%s] after Remove: %#v\n", name, out)
				initial = out
			}

		default:
			log.Printf("unknown plugin %q; skipping", name)
		}
	}

	// 3) cleanup
	for _, c := range clients {
		c.Close()
	}
	for _, c := range cmds {
		c.Wait()
	}
}
