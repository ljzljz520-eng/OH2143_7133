package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type Settings struct {
	DataPath       string `json:"data_path"`
	DefaultProfile string `json:"default_profile"`
	DefaultLocale  string `json:"default_locale"`
	DefaultTheme   string `json:"default_theme"`
	DefaultFitMode string `json:"default_fit_mode"`
	DimBackground  bool   `json:"dim_background"`
	Fullscreen     bool   `json:"fullscreen"`
}

func DefaultSettings(base string) Settings {
	return Settings{DataPath: filepath.Join(base, "wedding-sign.db"), DefaultProfile: "ceremony", DefaultLocale: "zh-CN", DefaultTheme: "garden", DefaultFitMode: "cover", DimBackground: true, Fullscreen: true}
}

func (s Settings) Validate() error {
	if strings.TrimSpace(s.DataPath) == "" {
		return errors.New("data path is required")
	}
	if s.DefaultFitMode != "cover" && s.DefaultFitMode != "contain" {
		return errors.New("default fit mode is invalid")
	}
	if s.DefaultProfile == "" || s.DefaultLocale == "" || s.DefaultTheme == "" {
		return errors.New("default profile, locale, and theme are required")
	}
	return nil
}

func Load(path string) (Settings, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return Settings{}, err
	}
	var settings Settings
	if err := json.Unmarshal(raw, &settings); err != nil {
		return Settings{}, err
	}
	if err := settings.Validate(); err != nil {
		return Settings{}, err
	}
	return settings, nil
}

func Save(path string, settings Settings) error {
	if err := settings.Validate(); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	raw, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, raw, 0o644)
}
