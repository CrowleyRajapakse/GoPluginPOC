package main

import (
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"example.com/policy"
	"github.com/google/uuid"
	"github.com/hashicorp/go-plugin"
)

func main() {
	// initial header uses host's uuid v1.0.0
	baseHeaders := map[string]string{"Host-UUID": uuid.New().String()}

	pluginsDir := "plugins/"
	entries, err := os.ReadDir(pluginsDir)
	if err != nil {
		log.Fatalf("read plugins dir: %v", err)
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		path := filepath.Join(pluginsDir, e.Name())
		info, err := os.Stat(path)
		if err != nil || info.Mode()&fs.ModePerm&0111 == 0 {
			continue
		}

		log.Printf("Loading plugin %s…", e.Name())
		client := plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig: policy.Handshake,
			Plugins:         map[string]plugin.Plugin{"policy": &policy.PolicyPlugin{}},
			Cmd:             exec.Command(path),
		})
		defer client.Kill()

		rpcClient, err := client.Client()
		if err != nil {
			log.Printf("Client error: %v", err)
			continue
		}
		raw, err := rpcClient.Dispense("policy")
		if err != nil {
			log.Printf("Dispense error: %v", err)
			continue
		}
		pol := raw.(policy.Policy)

		// each plugin uses its own uuid dependency version
		headers := make(map[string]string)
		for k, v := range baseHeaders {
			headers[k] = v
		}

		out, err := pol.ProcessRequestHeaders(headers)
		if err != nil {
			log.Printf("Plugin error: %v", err)
			continue
		}
		log.Printf("Plugin %s output: %v", e.Name(), out)
	}
}
