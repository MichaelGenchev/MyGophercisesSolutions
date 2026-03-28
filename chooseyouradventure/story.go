package main

import "encoding/json"

type Adventure struct {
	Start string         `json:"start_arc"`
	Arcs  map[string]Arc `json:"arcs"`
}
type Arc struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Option `json:"options"`
}

type Option struct {
	Text    string `json:"text"`
	ArcName string `json:"arc"`
}

func NewAdventure(jsonFile []byte) (*Adventure, error) {
	var adventure Adventure

	err := json.Unmarshal(jsonFile, &adventure)
	if err != nil {
		return nil, err
	}

	return &adventure, nil
}
