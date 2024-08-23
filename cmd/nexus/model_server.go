package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

type Item struct {
	Name string
	Url  string

	// InsecureSkipTLSVerify bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

type ModelServers []Item

func (s *ModelServers) Add(name string, url string) error {
	if s.Find(name) {
		return fmt.Errorf("server %v already exists", name)
	}

	server := Item{
		Name:      name,
		Url:       url,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

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

func (s *ModelServers) Update(name string, server Item) error {
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
