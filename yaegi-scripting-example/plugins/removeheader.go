// plugins/removeheader.go
package main

func Process(headers map[string]string) (map[string]string, error) {
	delete(headers, "X-Remove")
	return headers, nil
}
