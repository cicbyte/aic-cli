package utils

import (
	"fmt"
)

func InitAppDirs() error {
	config := ConfigInstance

	dirs := []string{
		config.GetAppSeriesDir(),
		config.GetAppDir(),
		config.GetConfigDir(),
		config.GetLogDir(),
	}

	for _, dir := range dirs {
		if err := EnsureDir(dir); err != nil {
			return fmt.Errorf("directory init failed: %v", err)
		}
	}

	return nil
}
