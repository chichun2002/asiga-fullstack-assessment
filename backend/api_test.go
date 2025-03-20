package main

import (
	"net/http"
	"testing"
)

func TestGetProducts(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/products")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}
