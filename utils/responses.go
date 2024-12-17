package utils

import (
	"encoding/json"
)

type Metadata struct {
	mp map[string]string
}

func NewMetadata() *Metadata {
	return &Metadata{
		mp: make(map[string]string),
	}
}

func (m *Metadata) Set(key string, value string) *Metadata {
	m.mp[key] = value
	return m
}

func (m *Metadata) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.mp)
}
