package main

import "github.com/google/uuid"

// Process adds X-Added-By using uuid v1.3.0 from vendor/
func Process(headers map[string]string) (map[string]string, error) {
	headers["X-Added-By"] = uuid.New().String()
	return headers, nil
}
