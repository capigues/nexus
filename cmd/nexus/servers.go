package main

import (
	"encoding/json"
	"errors"
	"fmt"
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

func (s *ModelServers) Serve() error {
	// ENDPOINTS
	// /health - If endpoint is healthy
	// /v1/models - list of all API/Models managed by Nexus
	// /version - version of this app?
	// /v1/chat/completions - can use same request body and vllm endpoint but substitute model for Nexus API name to route traffic
	// /v1/completion? - MAYBE. probably not in first iteration
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
