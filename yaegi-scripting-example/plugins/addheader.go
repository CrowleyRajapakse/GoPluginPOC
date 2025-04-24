// plugins/addheader.go
package main

func Process(headers map[string]string) (map[string]string, error) {
	headers["X-Added-By"] = "AddHeaderPlugin"
	return headers, nil
}
