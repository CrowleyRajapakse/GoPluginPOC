package main

import (
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/google/uuid"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

// loadPlugin interprets the Go plugin at path with modules disabled
func loadPlugin(path string) (reflect.Value, error) {
	dir := filepath.Dir(path)
	cwd, err := os.Getwd()
	if err != nil {
		return reflect.Value{}, err
	}
	defer os.Chdir(cwd)

	// chdir into plugin dir so vendor/ is seen
	if err := os.Chdir(dir); err != nil {
		return reflect.Value{}, err
	}

	// disable modules so interpreter uses vendor/
	env := append(os.Environ(), "GO111MODULE=off")
	i := interp.New(interp.Options{Env: env})
	i.Use(stdlib.Symbols)

	// interpret the plugin source file
	if _, err := i.EvalPath(filepath.Base(path)); err != nil {
		return reflect.Value{}, err
	}
	// lookup main.Process
	v, err := i.Eval("main.Process")
	if err != nil {
		return reflect.Value{}, err
	}
	return v, nil
}

func main() {
	// stamp a unique Host-UUID using uuid dependency
	baseHeaders := map[string]string{
		"Host-UUID": uuid.New().String(),
		"Hello":     "World",
		"X-Remove":  "bye",
	}

	pluginsRoot := filepath.Join("..", "plugins")
	log.Printf("Scanning plugins in %s", pluginsRoot)

	dirs, err := os.ReadDir(pluginsRoot)
	if err != nil {
		log.Fatalf("read plugins root: %v", err)
	}

	for _, d := range dirs {
		if !d.IsDir() {
			continue
		}
		pluginDir := filepath.Join(pluginsRoot, d.Name())
		files, _ := os.ReadDir(pluginDir)
		for _, f := range files {
			if !strings.HasSuffix(f.Name(), ".go") {
				continue
			}
			path := filepath.Join(pluginDir, f.Name())
			log.Printf("Loading plugin: %s", path)

			procVal, err := loadPlugin(path)
			if err != nil {
				log.Printf(" load error: %v", err)
				continue
			}

			// copy baseHeaders for isolation
			headers := make(map[string]string)
			for k, v := range baseHeaders {
				headers[k] = v
			}

			// interpret and invoke Process
			results := procVal.Call([]reflect.Value{reflect.ValueOf(headers)})
			newH := results[0].Interface().(map[string]string)
			if !results[1].IsNil() {
				log.Printf(" plugin error: %v", results[1].Interface().(error))
				continue
			}
			log.Printf(" plugin %s output: %v", f.Name(), newH)
		}
	}
}
