package html_template

import (
	"encoding/json"
	"os"
)

type Definition struct {
	Key     string `json:"Key"`
	Prompt  string `json:"Prompt"`
	Default any    `json:"Default"`
	Type    string `json:"Type"`
}

type Definitions []Definition

type HTMLTemplateConfiguration struct {
	Name        string
	Description string
	Definitions Definitions
}

func GetHTMLTemplateConfigurationFromFile(path string) (*HTMLTemplateConfiguration, error) {
	configDataBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var result HTMLTemplateConfiguration
	if err := json.Unmarshal(configDataBytes, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (result *HTMLTemplateConfiguration) IsKeySet(key string, rawData map[string]interface{}) bool {
	return true
}
