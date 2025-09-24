package templates

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

type TemplateExecutor struct {
	TemplateId string
	Parameters map[string]any
}

func (t *TemplateExecutor) Execute() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		zap.L().Error(fmt.Sprintf("Failed to get current working directory err=%v", err))
		return "", err
	}

	templatePath := filepath.Join(currentDir, "internal/server/templates", t.TemplateId, "template.yaml")

	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		zap.L().Error(fmt.Sprintf("Template file not found: %s", templatePath))
		return "", fmt.Errorf("template file not found: %s", templatePath)
	}

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Failed to parse template file %s err=%v", templatePath, err))
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, t.Parameters)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Failed to execute template %s err=%v", templatePath, err))
		return "", err
	}

	encodedStr := base64.StdEncoding.EncodeToString(buf.Bytes())
	return encodedStr, nil
}

func RetrieveExecutor(templateId string, parameters string) (*TemplateExecutor, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current working directory: %w", err)
	}

	templateDir := filepath.Join(currentDir, "internal/server/templates", templateId)

	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("template id %v does not exist", templateId)
	}

	descriptionPath := filepath.Join(templateDir, "description.json")
	if _, err := os.Stat(descriptionPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("template description not found for id %v", templateId)
	}

	templatePath := filepath.Join(templateDir, "template.yaml")
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("template file not found for id %v", templateId)
	}

	var params map[string]any
	if err := json.Unmarshal([]byte(parameters), &params); err != nil {
		return nil, fmt.Errorf("failed to parse body: %v", err)
	}

	zap.L().Info(fmt.Sprintf("Creating generic template executor for template: %s", templateId))

	return &TemplateExecutor{
		TemplateId: templateId,
		Parameters: params,
	}, nil
}
