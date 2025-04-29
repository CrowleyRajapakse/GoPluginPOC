package main

import "github.com/google/uuid"

// Process deletes X-Remove, adds X-Removed-By using uuid v1.1.1 from vendor/
func Process(headers map[string]string) (map[string]string, error) {
	delete(headers, "X-Remove")
	headers["X-Removed-By"] = uuid.New().String()
	return headers, nil
}
