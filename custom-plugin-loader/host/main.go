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

	"github.com/google/uuid"
)

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

type stdioRW struct {
	in  io.Reader
	out io.Writer
}

func (s *stdioRW) Read(p []byte) (int, error)  { return s.in.Read(p) }
func (s *stdioRW) Write(p []byte) (int, error) { return s.out.Write(p) }
func (s *stdioRW) Close() error                { return nil }

// PluginClient holds the RPC client and process for a plugin
// Name: executable filename
// Client: JSON-RPC client
// Cmd:    running subprocess

type PluginClient struct {
	Name   string
	Client *rpc.Client
	Cmd    *exec.Cmd
}

// loadPlugin starts a plugin executable and returns its RPC client + Cmd
func loadPlugin(path string) (*rpc.Client, *exec.Cmd) {
	cmd := exec.Command(path)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalf("stdin pipe error for %s: %v", path, err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("stdout pipe error for %s: %v", path, err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatalf("start plugin %s: %v", path, err)
	}
	rw := &stdioRW{in: stdout, out: stdin}
	client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(rw))
	return client, cmd
}

// loadClients finds all executables in dir and loads them
func loadClients(dir string) ([]PluginClient, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read plugin dir %s: %w", dir, err)
	}
	var out []PluginClient
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		full := filepath.Join(dir, e.Name())
		info, err := os.Stat(full)
		if err != nil || info.Mode()&0111 == 0 {
			continue // skip non-executables
		}
		client, cmd := loadPlugin(full)
		out = append(out, PluginClient{Name: e.Name(), Client: client, Cmd: cmd})
	}
	return out, nil
}

// processClients attempts both Add and Remove on every plugin
func processClients(pcs []PluginClient, headers Headers) Headers {
	current := headers
	for _, pc := range pcs {
		fmt.Printf("-- Processing %s --\n", pc.Name)

		// Try Add
		var afterAdd Headers
		addArgs := AddArgs{Key: "X-Trace-Id", Value: "", Headers: current}
		if err := pc.Client.Call("Plugin.Add", addArgs, &afterAdd); err == nil {
			fmt.Printf("[%s] Add -> %#v\n", pc.Name, afterAdd)
			current = afterAdd
		} else if !strings.Contains(err.Error(), "not found") {
			log.Printf("[%s] AddHeader is not implemented in this plugin: %v", pc.Name, err)
		}

		// Try Remove
		var afterRem Headers
		remArgs := RemoveArgs{Key: "RequestID", Headers: current}
		if err := pc.Client.Call("Plugin.Remove", remArgs, &afterRem); err == nil {
			fmt.Printf("[%s] Remove -> %#v\n", pc.Name, afterRem)
			current = afterRem
		} else if !strings.Contains(err.Error(), "not found") {
			log.Printf("[%s] RemoveHeader is not implemented in this plugin: %v", pc.Name, err)
		}
	}
	return current
}

func main() {
	pluginDir := "./plugins"

	// 1) load all plugins
	pcs, err := loadClients(pluginDir)
	if err != nil {
		log.Fatalf("load clients: %v", err)
	}

	// 2) initial headers
	initial := Headers{"RequestID": uuid.New().String()}

	// 3) process plugins dynamically
	final := processClients(pcs, initial)
	fmt.Printf("Final headers: %#v\n", final)

	// 4) cleanup
	for _, pc := range pcs {
		pc.Client.Close()
		pc.Cmd.Wait()
	}
}
