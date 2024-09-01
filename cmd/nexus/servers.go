package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/alexeyco/simpletable"
)

type ModelServers []Server

func (s *ModelServers) Add(server *Server) error {
	if s.Find(server.Name) {
		return fmt.Errorf("API %v already exists", server.Name)
	}

	err := server.GetInfo()
	if err != nil {
		return err
	}

	*s = append(*s, *server)
	return s.Store()
}

func (s *ModelServers) Remove(name string) error {
	if !s.Find(name) {
		return fmt.Errorf("API %v does not exists", name)
	}

	updated := ModelServers{}

	for _, item := range *s {
		if item.Name != name {
			updated = append(updated, item)
		}
	}

	*s = updated
	return s.Store()
}

func (s *ModelServers) Update(name string, server Server) error {
	updated := ModelServers{}

	for _, item := range *s {
		if item.Name != name {
			updated = append(updated, item)
		} else {
			updated = append(updated, server)
		}
	}

	*s = updated
	return s.Store()
}

func (s *ModelServers) GetServer(name string) (*Server, error) {
	for _, server := range *s {
		if server.Name == name {
			return &server, nil
		}
	}

	return nil, fmt.Errorf("API %v not found", name)
}

func (s *ModelServers) List() error {
	table := simpletable.New()
	httpErrors := ""

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "NAME"},
			{Align: simpletable.AlignCenter, Text: "URL"},
			{Align: simpletable.AlignCenter, Text: "MODEL"},
			{Align: simpletable.AlignCenter, Text: "UPDATED AT"},
			{Align: simpletable.AlignCenter, Text: "STATUS"},
		},
	}

	for _, server := range *s {
		row := []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: server.Name},
			{Align: simpletable.AlignCenter, Text: server.Url},
			{Align: simpletable.AlignCenter, Text: server.ModelName},
			{Align: simpletable.AlignCenter, Text: server.UpdatedAt.Format(time.ANSIC)},
			{Align: simpletable.AlignCenter, Text: fmt.Sprintf("%v", server.Status)},
		}

		table.Body.Cells = append(table.Body.Cells, row)
	}

	if len(httpErrors) > 0 {
		return errors.New(httpErrors)
	}

	fmt.Println(table.String())
	return nil
}

func (s *ModelServers) Load() error {
	NEXUS_SERVERS_PATH := os.Getenv("NEXUS_SERVERS_PATH")

	file, err := os.ReadFile(NEXUS_SERVERS_PATH)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return err
	}

	if len(file) == 0 {
		return err
	}

	err = json.Unmarshal(file, s)
	if err != nil {
		return err
	}

	return nil
}

func (s *ModelServers) Store() error {
	NEXUS_SERVERS_PATH := os.Getenv("NEXUS_SERVERS_PATH")

	file, err := json.Marshal(*s)
	if err != nil {
		return err
	}

	err = os.WriteFile(NEXUS_SERVERS_PATH, file, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (s *ModelServers) Find(name string) bool {
	for _, server := range *s {
		if server.Name == name {
			return true
		}
	}

	return false
}

type ChatCompletionsRequestBody struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Stream      bool      `json:"stream"`
	Temperature float64   `json:"temperature"`
}

func (s *ModelServers) Serve(out io.Writer, port string) error {
	// ENDPOINTS

	// GET
	// /health - If endpoint is healthy
	// /v1/models - list of all API/Models managed by Nexus
	// /version - version of this app?

	// POST
	// /v1/chat/completions - can use same request body and vllm endpoint but substitute model for Nexus API name to route traffic
	// /v1/completion? - MAYBE. probably not in first iteration

	// http.HandleFunc("/v1/models", modelsHandler)
	// http.HandleFunc("/versions", versionsHandler)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(out, "Handling GET /health\n")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Nexus server is running\n"))
	})

	http.HandleFunc("/v1/chat/completions", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(405)
			w.Write([]byte("Method not accepted\n"))
			return
		}
		fmt.Fprintf(out, "Handling POST /v1/chat/completions\n")

		var body ChatCompletionsRequestBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		var server *Server
		var err error
		if server, err = s.GetServer(body.Model); err != nil {
			fmt.Fprintf(w, "API %v not found\n", body.Model)
			return
		}

		body.Model = server.ModelName

		requestBytes, err := json.Marshal(body)
		if err != nil {
			http.Error(w, "Could not create openAI request body", http.StatusInternalServerError)
			return
		}

		req, err := http.NewRequest("POST", server.Url+"/v1/chat/completions", bytes.NewBuffer(requestBytes))
		if err != nil {
			http.Error(w, "Could not create http POST request", http.StatusBadRequest)
			return
		}

		client := server.createClient(req)

		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), resp.StatusCode)
			return
		}
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		additionalData := []byte("\n")
		bodyBytes = append(bodyBytes, additionalData...)

		w.Write(bodyBytes)
	})

	// http.HandleFunc("/v1/completions", completionsHandler)

	fmt.Printf("Starting server on port %v\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		return fmt.Errorf("error starting server: %v", err.Error())
	}

	return nil
}

// func modelsHandler(w http.ResponseWriter, r *http.Request) {
// curl -X 'GET' \
// 'https://vllm-predictor-rhsaia-model-serving.apps.rhsaia.vg6c.p1.openshiftapps.com/v1/models' \
// -H 'accept: application/json'

// {
// 	"object": "list",
// 	"data": [
// 	  {
// 		"id": "Mistral-7B-Instruct-v0.3",
// 		"object": "model",
// 		"created": 1725169224,
// 		"owned_by": "vllm",
// 		"root": "Mistral-7B-Instruct-v0.3",
// 		"parent": null,
// 		"permission": [
// 		  {
// 			"id": "modelperm-104eeac59e0f453e97e963682a94feab",
// 			"object": "model_permission",
// 			"created": 1725169224,
// 			"allow_create_engine": false,
// 			"allow_sampling": true,
// 			"allow_logprobs": true,
// 			"allow_search_indices": false,
// 			"allow_view": true,
// 			"allow_fine_tuning": false,
// 			"organization": "*",
// 			"group": null,
// 			"is_blocking": false
// 		  }
// 		]
// 	  }
// 	]
// }
// }

// func versionsHandler(w http.ResponseWriter, r *http.Request) {
// 	curl -X 'GET' \
//   'https://vllm-predictor-rhsaia-model-serving.apps.rhsaia.vg6c.p1.openshiftapps.com/version' \
//   -H 'accept: application/json'

// {
// 	"version": "0.4.2"
// }
// }

// func chatCompletionsHandler(w http.ResponseWriter, r *http.Request) {

// REQUEST
// https://platform.openai.com/docs/api-reference/chat/create

// RESPONSE
// {
// 	"object": "error",
// 	"message": "[{'type': 'value_error', 'loc': ('body',), 'msg': \"Value error, You can only use one kind of guided decoding ('guided_json', 'guided_regex' or 'guided_choice').\", 'input': {'messages': [{'content': 'string', 'role': 'system', 'name': 'string'}, {'content': 'string', 'role': 'user', 'name': 'string'}, {'content': 'string', 'role': 'assistant', 'name': 'string', 'function_call': {'arguments': 'string', 'name': 'string'}, 'tool_calls': [{'id': 'string', 'function': {'arguments': 'string', 'name': 'string'}, 'type': 'function'}]}, {'content': 'string', 'role': 'tool', 'name': 'string', 'tool_call_id': 'string'}, {'content': 'string', 'role': 'function', 'name': 'string'}], 'model': 'string', 'frequency_penalty': 0, 'logit_bias': {'additionalProp1': 0, 'additionalProp2': 0, 'additionalProp3': 0}, 'logprobs': False, 'top_logprobs': 0, 'max_tokens': 0, 'n': 1, 'presence_penalty': 0, 'response_format': {'type': 'text'}, 'seed': 0, 'stop': 'string', 'stream': False, 'temperature': 0.7, 'top_p': 1, 'user': 'string', 'best_of': 0, 'use_beam_search': False, 'top_k': -1, 'min_p': 0, 'repetition_penalty': 1, 'length_penalty': 1, 'early_stopping': False, 'ignore_eos': False, 'min_tokens': 0, 'stop_token_ids': [0], 'skip_special_tokens': True, 'spaces_between_special_tokens': True, 'echo': False, 'add_generation_prompt': True, 'include_stop_str_in_output': False, 'guided_json': 'string', 'guided_regex': 'string', 'guided_choice': ['string'], 'guided_grammar': 'string', 'guided_decoding_backend': 'string', 'guided_whitespace_pattern': 'string'}, 'ctx': {'error': ValueError(\"You can only use one kind of guided decoding ('guided_json', 'guided_regex' or 'guided_choice').\")}}]",
// 	"type": "BadRequestError",
// 	"param": null,
// 	"code": 400
// }

// }

// func completionsHandler(w http.ResponseWriter, r *http.Request) {
// 	REQUEST
// https://platform.openai.com/docs/guides/completions

// RESPONSE

// }
