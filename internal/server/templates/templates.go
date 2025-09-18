package templates

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

type Description struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Parameters  map[string]string `json:"parameters"`
}

func Templates() ([]Description, error) {
	var descriptions []Description

	currentDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("Error getting working directory.")
	}

	templatesDir := filepath.Join(currentDir, "templates")

	entries, err := os.ReadDir(templatesDir)
	if err != nil {
		zap.L().Error("Error reading base_dir")
		return nil, fmt.Errorf("error reading directory %s: %w", templatesDir, err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		descriptionPath := filepath.Join(templatesDir, entry.Name(), "description.json")

		if _, err := os.Stat(descriptionPath); os.IsNotExist(err) {
			continue
		}

		fileData, err := os.ReadFile(descriptionPath)
		if err != nil {
			zap.L().Error(fmt.Sprintf("Error reading %s: %v\n", descriptionPath, err))
			continue
		}

		var description Description

		if err := json.Unmarshal(fileData, &description); err != nil {
			zap.L().Error(fmt.Sprintf("Error parsing JSON from %s: %v\n", descriptionPath, err))
			continue
		}

		descriptions = append(descriptions, description)
	}

	return descriptions, nil
}

func Template(templateId string) (Description, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return Description{}, fmt.Errorf("Error getting working directory.")
	}

	templatesDir := filepath.Join(currentDir, "templates")

	descriptionPath := filepath.Join(templatesDir, templateId, "description.json")

	if _, err := os.Stat(descriptionPath); os.IsNotExist(err) {
		return Description{}, err
	}

	fileData, err := os.ReadFile(descriptionPath)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Error reading %s: %v\n", descriptionPath, err))
		return Description{}, err
	}

	var description Description

	if err := json.Unmarshal(fileData, &description); err != nil {
		zap.L().Error(fmt.Sprintf("Error parsing JSON from %s: %v\n", descriptionPath, err))
		return description, err
	}

	return description, nil
}

func TemplateIDs() ([]string, error) {
	var ids []string

	currentDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("Error getting working directory.")
	}

	templatesDir := filepath.Join(currentDir, "templates")

	entries, err := os.ReadDir(templatesDir)
	if err != nil {
		return nil, fmt.Errorf("error reading directory %s: %w", templatesDir, err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		descriptionPath := filepath.Join(templatesDir, entry.Name(), "description.json")

		if _, err := os.Stat(descriptionPath); os.IsNotExist(err) {
			continue
		}

		fileData, err := os.ReadFile(descriptionPath)
		if err != nil {
			zap.L().Error(fmt.Sprintf("Error reading %s: %v\n", descriptionPath, err))
			continue
		}

		var description Description

		if err := json.Unmarshal(fileData, &description); err != nil {
			zap.L().Error(fmt.Sprintf("Error parsing JSON from %s: %v\n", descriptionPath, err))
			continue
		}

		ids = append(ids, description.ID)
	}

	return ids, nil
}
