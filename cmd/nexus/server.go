package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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

	if responseData.Data == nil {
		return fmt.Errorf("could not refresh %v API. Is %v the correct url?", s.Name, s.Url)
	}

	s.ModelName = responseData.Data[0].ID
	s.Status = Healthy
	s.UpdatedAt = time.Now()

	return nil
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (s *Server) Chat(temperature, max_tokens int) error {
	fmt.Printf("Starting chat with %v temperature and %v max tokens\n", temperature, max_tokens)

	// Create a scanner to read user input
	scanner := bufio.NewScanner(os.Stdin)

	messages := &[]Message{}

	for {
		// Prompt user for input
		fmt.Print("You: ")
		scanner.Scan()
		userInput := scanner.Text()

		// Check if the user wants to exit
		if strings.ToLower(userInput) == "exit" {
			fmt.Println("Goodbye!")
			break
		}

		*messages = append(*messages, Message{
			Role:    "user",
			Content: userInput,
		})
		// Get response from the chatbot
		response := generateResponse(*messages, s.ModelName, s.ApiKey, s.Url, s.InsecureSkipTLSVerify)
		*messages = append(*messages, Message{
			Role:    "assistant",
			Content: response,
		})

		// Print the chatbot's response
		fmt.Printf("%v: %v\n", s.ModelName, response)
	}

	return nil
}

type OpenAIChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type OpenAIChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func generateResponse(messages []Message, modelName, apiKey, url string, insecureSkipTLSVerify bool) string {

	requestBody := OpenAIChatRequest{
		Model:    modelName, // Replace with your model
		Messages: messages,
	}

	reqBody, err := json.Marshal(requestBody)
	if err != nil {
		return "Error: Unable to marshal request body."
	}

	req, err := http.NewRequest("POST", url+"/v1/chat/completions", bytes.NewBuffer(reqBody))
	if err != nil {
		return "Error: Unable to create request."
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureSkipTLSVerify},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		return "Error: Unable to complete request."
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "Error: Unable to read response body."
	}

	var openaiResp OpenAIChatResponse
	err = json.Unmarshal(body, &openaiResp)
	if err != nil {
		return "Error: Unable to unmarshal response body."
	}

	return openaiResp.Choices[0].Message.Content
}
