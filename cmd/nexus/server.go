package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type serverStatus string // "healthy", "unhealthy"

var (
	Healthy   serverStatus = "Healthy"
	Unhealthy serverStatus = "Unhealthy"
)

type Server struct {
	Name      string
	Url       string
	Status    serverStatus
	ModelName string

	ApiKey                string
	InsecureSkipTLSVerify bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Model struct {
	ID string `json:"id"`
}

type ResponseData struct {
	Data   []Model `json:"data"`
	Object string  `json:"object"`
}

type openAIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param"`
	Code    string `json:"code"`
}

type OpenAIErrorResponse struct {
	Error openAIError `json:"error"`
}

func (s *Server) GetInfo() error {
	req, err := http.NewRequest("GET", s.Url+"/v1/models", nil)
	if err != nil {
		return fmt.Errorf("%v", err.Error())
	}

	tr := &http.Transport{}
	if s.InsecureSkipTLSVerify {
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	if s.ApiKey != "" {
		req.Header.Add("Authorization", "Bearer "+s.ApiKey)
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   3 * time.Second,
	}

	// Only supports Single Model Serving APIs
	resp, err := client.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "x509: certificate signed by unknown authority") {
			return fmt.Errorf("failed to verify certificate: x509: certificate signed by unknown authority")
		}

		return fmt.Errorf("unable to connect to %v. ", s.Url)
	}
	defer resp.Body.Close()

	// Read and buffer the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Create a buffer from the read body bytes
	bodyBuffer := bytes.NewBuffer(bodyBytes)

	var errorData OpenAIErrorResponse
	if err := json.NewDecoder(bodyBuffer).Decode(&errorData); err != nil {
		return err
	}

	if strings.Contains(errorData.Error.Message, "API key") {
		return fmt.Errorf("you need to provide your API key")
	}

	// Reset the buffer to the start
	bodyBuffer.Reset()
	bodyBuffer.Write(bodyBytes)

	var responseData ResponseData
	if err := json.NewDecoder(bodyBuffer).Decode(&responseData); err != nil {
		s.Status = Unhealthy
		return err
	}

	s.ModelName = responseData.Data[0].ID
	s.Status = Healthy
	s.UpdatedAt = time.Now()

	return nil
}
