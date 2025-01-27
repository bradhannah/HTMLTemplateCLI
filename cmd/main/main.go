package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"strings"
	"time"
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

func stringJoin(slice interface{}, sep string) string {
	arr, ok := slice.([]interface{})
	if !ok {
		return ""
	}
	var strItems []string
	for _, item := range arr {
		// Convert every element to a string (fmt.Sprintf)
		strItems = append(strItems, fmt.Sprintf("%v", item))
	}
	return strings.Join(strItems, sep)
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

func GetHtmlTemplate(path string) (*template.Template, error) {
	htmlTemplateData, err := os.ReadFile(path)
	htmlTemplate := string(htmlTemplateData)

	if err != nil {
		return nil, err
		// log.Fatalf("Error reading file: %v", err)
	}

	funcMap := template.FuncMap{
		"join": func(sep string, items []interface{}) string {
			return stringJoin(items, sep)
			// return strings.Join(items, sep)
		},
		"currentDate": func(format string) string {
			return time.Now().Format(format)
		},
	}

	tmpl, err := template.New("template").Funcs(funcMap).Parse(htmlTemplate)

	return tmpl, err
}

func GetJsonFileAsMap(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
		// log.Fatalf("Error reading file: %v", err)
	}

	// Prepare a map to hold the JSON data
	var result map[string]interface{}

	// Unmarshal (deserialize) the JSON into the map
	err = json.Unmarshal(data, &result)
	return result, err
}

func main() {
	templateConfig, err := GetHTMLTemplateConfigurationFromFile("input/CoverLetterInputDefinition.json")
	if err != nil {
		panic(err)
		// log.Fatal(err)
	}

	_ = templateConfig

	result, err := GetJsonFileAsMap("input/BradHannah.json")
	if err != nil {
		panic(err)
		// log.Fatalf("Error unmarshaling JSON: %v", err)
	}

	// Parse the template string into a *template.Template.

	tmpl, err := GetHtmlTemplate("input/BradHannahCoverLetter.gohtml")

	if err != nil {
		panic(err)
	}

	if err := tmpl.Execute(os.Stdout, &result); err != nil {
		panic(err)
	}
}
