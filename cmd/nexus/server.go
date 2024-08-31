package main

import (
	"fmt"
	"net/http"
	"time"
)

type serverStatus string // "healthy", "unhealthy"
type serverType string   // "openAI (/v1/models), ollama (/api/tags)"

type Server struct {
	Name   string
	Url    string
	Status serverStatus
	Type   serverType

	// InsecureSkipTLSVerify bool
	// Add fields that will be returned by /v1/models

	CreatedAt time.Time
	UpdatedAt time.Time
}

// type apiData struct {
// }

func (s *Server) GetInfo() {
	client := http.Client{
		Timeout: 3 * time.Second,
	}

	resp, err := client.Get(s.Url + "/v1/models")
	if err != nil {

	}

	fmt.Println(resp.Body)
}
