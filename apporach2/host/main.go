package main

import (
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"example.com/mypkg/policy"
	"github.com/hashicorp/go-plugin"
)

func main() {
	pluginsDir := "./plugins"

	entries, err := os.ReadDir(pluginsDir)
	if err != nil {
		log.Fatalf("failed to read plugins dir: %v", err)
	}

	for _, e := range entries {
		if e.IsDir() || e.Name()[0] == '.' {
			continue
		}

		pluginPath := filepath.Join(pluginsDir, e.Name())
		info, err := os.Stat(pluginPath)
		if err != nil {
			log.Printf("skipping %s: %v", pluginPath, err)
			continue
		}

		// only executables
		if info.Mode()&fs.ModePerm&0111 == 0 {
			continue
		}

		log.Printf("loading plugin %s…", e.Name())

		client := plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig: policy.Handshake,
			Plugins: map[string]plugin.Plugin{
				"policy": &policy.PolicyPlugin{},
			},
			Cmd: exec.Command(pluginPath),
		})
		defer client.Kill()

		rpcClient, err := client.Client()
		if err != nil {
			log.Printf(" failed to start client: %v", err)
			continue
		}

		raw, err := rpcClient.Dispense("policy")
		if err != nil {
			log.Printf(" dispense error: %v", err)
			continue
		}

		pol := raw.(policy.Policy)

		// sample headers
		headers := map[string]string{
			"Hello":       "World",
			"X-To-Remove": "bye",
		}

		out, err := pol.ProcessRequestHeaders(headers)
		if err != nil {
			log.Printf(" plugin error: %v", err)
			continue
		}

		log.Printf(" plugin %s returned: %v", e.Name(), out)
	}
}
