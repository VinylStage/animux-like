package pet

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

type State struct {
	Name        string      `json:"name"`
	Species     SpeciesType `json:"species"`
	Hunger      int         `json:"hunger"`      // 0-100, 100 is full
	Happiness   int         `json:"happiness"`   // 0-100, 100 is happy
	Cleanliness int         `json:"cleanliness"` // 0-100, 100 is clean
	IsSick      bool        `json:"is_sick"`
	LastUpdate  time.Time   `json:"last_update"`
}

var ErrNoPet = errors.New("no pet found, please adopt one first")

func GetSavePath() string {
	xdgData := os.Getenv("XDG_DATA_HOME")
	if xdgData == "" {
		home, _ := os.UserHomeDir()
		xdgData = filepath.Join(home, ".local", "share")
	}
	return filepath.Join(xdgData, "animux", "state.json")
}

func LoadState() (*State, error) {
	path := GetSavePath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNoPet
		}
		return nil, err
	}
	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func SaveState(s *State) error {
	path := GetSavePath()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	s.LastUpdate = time.Now()
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func (s *State) Clamp() {
	if s.Hunger < 0 { s.Hunger = 0 }
	if s.Hunger > 100 { s.Hunger = 100 }
	if s.Happiness < 0 { s.Happiness = 0 }
	if s.Happiness > 100 { s.Happiness = 100 }
	if s.Cleanliness < 0 { s.Cleanliness = 0 }
	if s.Cleanliness > 100 { s.Cleanliness = 100 }
}
