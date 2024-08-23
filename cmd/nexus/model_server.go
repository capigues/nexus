package main

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

type item struct {
	Name string
	Url  string

	// InsecureSkipTLSVerify bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

type ModelServers []item

func (s *ModelServers) Add(name string, url string) {
	server := item{
		Name:      name,
		Url:       url,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	*s = append(*s, server)
}

func (s *ModelServers) Remove(name string) {
	updated := ModelServers{}

	for _, item := range *s {
		if item.Name != name {
			updated = append(updated, item)
		}
	}

	*s = updated
}

func (s *ModelServers) Update(name string, server item) {
	updated := ModelServers{}

	for _, item := range *s {
		if item.Name != name {
			updated = append(updated, item)
		} else {
			updated = append(updated, server)
		}
	}

	*s = updated
}

func (s *ModelServers) Load(filename string) error {
	file, err := os.ReadFile(filename)
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

func (s *ModelServers) Save(filename string) error {
	file, err := json.Marshal(*s)
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, file, 0644)
	if err != nil {
		return err
	}

	return nil
}
