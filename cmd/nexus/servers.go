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

func (s *ModelServers) Add(name string, url string) error {
	if s.Find(name) {
		return fmt.Errorf("server %v already exists", name)
	}

	server := Server{
		Name:      name,
		Url:       url,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	// TODO: Call api to check status and get more information (/v1/models)

	*s = append(*s, server)
	return s.Store()
}

func (s *ModelServers) Remove(name string) error {
	if !s.Find(name) {
		return fmt.Errorf("server %v does not exists", name)
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
			{Align: simpletable.AlignCenter, Text: "UPDATED AT"},
			// {Align: simpletable.AlignCenter, Text: "STATUS"},
		},
	}

	for _, server := range *s {
		// status, _ := s.CheckStatus(server.Url)
		// Check status only returns nil right now
		// if err != nil {
		// 	fmt.Println(err.Error())
		// }

		row := []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: server.Name},
			{Align: simpletable.AlignCenter, Text: server.Url},
			{Align: simpletable.AlignCenter, Text: server.UpdatedAt.Format(time.ANSIC)},
			// {Align: simpletable.AlignCenter, Text: fmt.Sprintf("%v", status)},
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
