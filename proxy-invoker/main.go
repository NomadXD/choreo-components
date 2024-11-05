package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

type ProxyRequest struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

func main() {
	http.HandleFunc("/proxy", handleProxy)
	log.Printf("Starting proxy server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleProxy(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the incoming request
	var proxyReq ProxyRequest
	if err := json.NewDecoder(r.Body).Decode(&proxyReq); err != nil {
		log.Printf("Error parsing request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create the proxied request
	client := &http.Client{}
	req, err := http.NewRequest(proxyReq.Method, proxyReq.URL, strings.NewReader(proxyReq.Body))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	// Add headers
	for key, value := range proxyReq.Headers {
		req.Header.Add(key, value)
	}

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error executing request: %v", err)
		http.Error(w, "Error executing request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Set response status code
	w.WriteHeader(resp.StatusCode)

	// Reads the entire response body into memory first
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return
	}

	// Send the response body
	if _, err := w.Write(body); err != nil {
		log.Printf("Error sending response: %v", err)
	}
}
