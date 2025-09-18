package templates

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

type CpuK3sDeployment struct {
	Host    string `json:"host"`
	Port    int16  `json:"port"`
	Token   string `json:"token"`
	Version string `json:"version"`
}

func (ck3d *CpuK3sDeployment) Execute() (string, error) {
	currentDir, err := os.Getwd()

	if err != nil {
		zap.L().Error(fmt.Sprintf("Failed to get current working directory err=%v", err))
		return "", err
	}

	templatesDir := filepath.Join(currentDir, "templates")

	templatePath := filepath.Join(templatesDir, "cpu_k3s_deployment", "template.yaml")

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Failed to parse template file=%s, err=%v\n", templatePath, err))
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, ck3d)

	if err != nil {
		zap.L().Error(fmt.Sprintf("Failed to execute the template engine err=%v\n", err))
		return "", err
	}

	encodedStr := base64.StdEncoding.EncodeToString(buf.Bytes())

	return encodedStr, nil
}
