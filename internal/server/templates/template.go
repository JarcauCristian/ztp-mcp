package templates

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
)

type GenericTemplate struct {
	Id              string      `json:"id" jsonschema_description:"The id of the template, should be lowercased and separated by underscores."`
	Name            string      `json:"name" jsonschema_description:"The name of the template, the same as the id, but with each word capitalized and replace the underscores with spaces."`
	Parameters      []Parameter `json:"parameters" jsonschema_description:"The parameters that will be placed inside the template.yaml to customize each deployment."`
	Description     string      `json:"description" jsonschema_description:"The description of the template."`
	UpdatePackages  bool        `json:"update_packages" jsonschema_description:"If true will update all the packages."`
	UpgradePackages bool        `json:"upgrade_packages" jsonschema_description:"If true will upgrade all the packages."`
	Packages        []string    `json:"packages" jsonschema_description:"The packages to install on the system."`
	Commands        []string    `json:"commands" jsonschema_description:"The commands to run when the system is installed."`
	Files           []File      `json:"files" jsonschema_description:"Specify the files that needs to be available on the system, such as config files and other files needed by the installed packages and applications."`
}

type Parameter struct {
	Name        string `json:"name" jsonschema_description:"The name of the parameter, needs to be written in Pascal case. If include it in template.yaml as templates needs to be done conform to Go html/template conventions."`
	Description string `json:"description" jsonschema_description:"The description about what the parameter is about."`
}

type File struct {
	Path    string `json:"path" jsonschema_description:"The path where the file will be created on the system."`
	Content string `json:"content" jsonschema_description:"The content of the files that will be written to the system."`
}

func Capitalize(value string) string {
	if value == "" {
		return ""
	}
	return strings.ToUpper(string(value[0])) + strings.ToLower(value[1:])
}

func CreateTemplate(genericTemplate GenericTemplate) error {
	currentDir, err := os.Getwd()
	if err != nil {
		zap.L().Error(fmt.Sprintf("Failed to get current working directory err=%v", err))
		return err
	}

	templateDir := filepath.Join(currentDir, "internal/server/templates/template")
	outputDir := filepath.Join(currentDir, "internal/server/templates", genericTemplate.Id)

	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Failed to create output directory %s err=%v", outputDir, err))
		return err
	}

	var returnErr error
	defer func() {
		if returnErr != nil {
			zap.L().Info(fmt.Sprintf("Cleaning up output directory %s due to error", outputDir))
			if removeErr := os.RemoveAll(outputDir); removeErr != nil {
				zap.L().Error(fmt.Sprintf("Failed to cleanup output directory %s err=%v", outputDir, removeErr))
			}
		}
	}()

	templateFiles, err := os.ReadDir(templateDir)
	if err != nil {
		returnErr = err
		zap.L().Error(fmt.Sprintf("Failed to read template directory %s err=%v", templateDir, err))
		return err
	}

	for _, file := range templateFiles {
		if file.IsDir() {
			continue
		}

		err := executeTemplateFile(templateDir, outputDir, file, genericTemplate)
		if err != nil {
			returnErr = err
			return err
		}
	}

	zap.L().Info(fmt.Sprintf("Successfully created template files for %s in %s", genericTemplate.Id, outputDir))
	return nil
}

func DeleteTemplate(templateId string) error {
	currentDir, err := os.Getwd()
	if err != nil {
		zap.L().Error(fmt.Sprintf("Failed to get current working directory err=%v", err))
		return err
	}

	templateDir := filepath.Join(currentDir, "internal/server/templates", templateId)

	_, err = os.ReadDir(templateDir)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Failed to retrieve template directory %s err=%v", templateDir, err))
		return err
	}

	err = os.RemoveAll(templateDir)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Failed to cleanup output directory %s err=%v", templateDir, err))
	}

	return nil
}

func executeTemplateFile(templateDir, outputDir string, file os.DirEntry, templ GenericTemplate) error {
	funcMap := template.FuncMap{
		"Capitalize": Capitalize,
		"ToLower":    strings.ToLower,
		"sub": func(a, b int) int {
			return a - b
		},
	}

	templatePath := filepath.Join(templateDir, file.Name())

	tmpl, err := template.New(file.Name()).Funcs(funcMap).ParseFiles(templatePath)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Failed to parse template file %s err=%v", templatePath, err))
		return err
	}

	outputFileName := strings.TrimSuffix(file.Name(), ".templ")

	outputPath := filepath.Join(outputDir, outputFileName)

	outputFile, err := os.Create(outputPath)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Failed to create output file %s err=%v", outputPath, err))
		return err
	}

	err = tmpl.Execute(outputFile, templ)
	if err != nil {
		outputFile.Close()
		zap.L().Error(fmt.Sprintf("Failed to execute template %s err=%v", templatePath, err))
		return err
	}

	outputFile.Close()
	zap.L().Info(fmt.Sprintf("Generated file: %s", outputPath))
	return nil
}
