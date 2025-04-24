package main

import (
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

// loadPlugin interprets the Go plugin at path and returns its Process symbol.
func loadPlugin(path string) (reflect.Value, error) {
	i := interp.New(interp.Options{})
	i.Use(stdlib.Symbols)

	src, err := os.ReadFile(path)
	if err != nil {
		return reflect.Value{}, err
	}
	if _, err := i.Eval(string(src)); err != nil {
		return reflect.Value{}, err
	}
	v, err := i.Eval("Process")
	if err != nil {
		return reflect.Value{}, err
	}
	return v, nil
}

func main() {
	pluginsDir := "./plugins"
	entries, err := os.ReadDir(pluginsDir)
	if err != nil {
		log.Fatalf("failed to read plugins dir: %v", err)
	}

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") {
			continue
		}
		pluginPath := filepath.Join(pluginsDir, e.Name())
		log.Printf("Loading plugin %s…", e.Name())

		processSym, err := loadPlugin(pluginPath)
		if err != nil {
			log.Printf("load error: %v", err)
			continue
		}

		// initial headers
		headers := map[string]string{
			"Hello":    "World",
			"X-Remove": "bye",
		}
		// call the Process function
		results := processSym.Call([]reflect.Value{reflect.ValueOf(headers)})
		newHeaders, ok := results[0].Interface().(map[string]string)
		if !ok {
			log.Printf("unexpected type from Process")
			continue
		}
		var procErr error
		if !results[1].IsNil() {
			procErr = results[1].Interface().(error)
		}
		if procErr != nil {
			log.Printf("plugin error: %v", procErr)
			continue
		}

		log.Printf("plugin %s output: %v", e.Name(), newHeaders)
	}
}
