package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"plugin"

	extprocv3 "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
	"github.com/wso2/apk/gateway/enforcer/pkg/plugins"
)

func main() {
	// find the .so next to the executable
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("os.Executable: %v", err)
	}
	soPath := filepath.Join(filepath.Dir(exePath), "addheader.so")
	log.Printf("loading plugin %q…", soPath)

	// 1) Open the plugin file
	p, err := plugin.Open(soPath)
	if err != nil {
		log.Fatalf("plugin.Open: %v", err)
	}

	// 2) Lookup the "Plugin" symbol
	sym, err := p.Lookup("Plugin")
	if err != nil {
		log.Fatalf("Lookup(\"Plugin\"): %v", err)
	}

	// 3) Because your plugin declares
	//       var Plugin plugins.Policy = &AddHeader{}
	//    Lookup returns *plugins.Policy, so we need to:
	policyPtr, ok := sym.(*plugins.Policy)
	if !ok {
		log.Fatalf("symbol Plugin has wrong type: %T (expected *plugins.Policy)", sym)
	}

	// 4) Dereference to get the actual Policy
	policy := *policyPtr
	log.Printf("successfully loaded policy %q", policy.Name())

	// 5) Invoke it directly (no RegisterPolicy() needed)
	reqHdrs, err := policy.ApplyRequestHeaders(context.Background(), &extprocv3.ProcessingRequest{})
	if err != nil {
		log.Fatalf("ApplyRequestHeaders: %v", err)
	}
	for _, h := range reqHdrs {
		log.Printf("→ add request header %s=%s", h.Header.Key, h.Header.Value)
	}

	resHdrs, err := policy.ApplyResponseHeaders(context.Background(), &extprocv3.ProcessingRequest{})
	if err != nil {
		log.Fatalf("ApplyResponseHeaders: %v", err)
	}
	for _, h := range resHdrs {
		log.Printf("→ add response header %s=%s", h.Header.Key, h.Header.Value)
	}
}
