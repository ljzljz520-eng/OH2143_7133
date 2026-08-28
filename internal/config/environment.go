package config

import (
	"os"
	"path/filepath"
	"strings"
)

func FromEnvironment() Settings {
	base := strings.TrimSpace(os.Getenv("WEDDING_SIGN_DATA"))
	if base == "" {
		base = ".data"
	}
	settings := DefaultSettings(base)
	if value := strings.TrimSpace(os.Getenv("WEDDING_SIGN_PROFILE")); value != "" {
		settings.DefaultProfile = value
	}
	if value := strings.TrimSpace(os.Getenv("WEDDING_SIGN_THEME")); value != "" {
		settings.DefaultTheme = value
	}
	if value := strings.TrimSpace(os.Getenv("WEDDING_SIGN_FIT")); value == "contain" {
		settings.DefaultFitMode = value
	}
	settings.DataPath = filepath.Clean(settings.DataPath)
	return settings
}

func EnsureDataDirectory(settings Settings) error {
	return os.MkdirAll(filepath.Dir(settings.DataPath), 0o755)
}
