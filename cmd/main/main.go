package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	ht "github.com/bradhannah/HTMLTemplateCLI/pkg/html_template"
)

var rootCmd = &cobra.Command{
	Use:   "example",
	Short: "A simple example demonstrating three required named parameters",
	RunE: func(cmd *cobra.Command, args []string) error {
		// If we've reached this point, Cobra has already verified that
		// the required flags are set, but you can do additional validation here if needed.

		// fmt.Printf("Arg1: %s\nArg2: %s\nArg3: %s\n", definitionPath, goHtmlTemplatePath, inputPath)
		return nil
	},
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

var (
	definitionPath     string
	goHtmlTemplatePath string
	inputPath          string
	outputPath         string
)

func promptAndWaitForAnswer(definition ht.Definition, separator string) (string, error) {
	reader := bufio.NewReader(os.Stdin)

	done := false
	var input string

	for !done {
		fmt.Printf("Field Details: %s\n", definition.Prompt)
		if definition.Type == "" {
			fmt.Print("Type: String\n")
		} else {
			fmt.Printf("Type: %s\n", definition.Type)
		}
		if definition.Type != "" {
			fmt.Printf("Default (press enter): %s\n", definition.Default)
		}

		fmt.Print("User Input: ")
		var err error
		input, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return "", err
		}
		input = strings.TrimSpace(input)

		if input == "" && definition.Default != "" {
			s, isString := definition.Default.(string)
			if isString {
				return s, nil
			}
			sl, isStringList := definition.Default.([]string)
			if isStringList {
				return strings.Join(sl, separator), nil
			}
			fmt.Print("Error, default input unknown")
			return "", errors.New("unknown input")
		}
		fmt.Println("")
		if input == "" {
			continue
		}
		done = true
	}
	return input, nil
}

func main() {
	rootCmd.Flags().StringVar(&definitionPath, "definition", "", "Definition file path")
	rootCmd.Flags().StringVar(&goHtmlTemplatePath, "gohtml", "", ".gohtml template file path")
	rootCmd.Flags().StringVar(&inputPath, "input", "", "Input file path")
	rootCmd.Flags().StringVar(&outputPath, "output", "", "HTML output file path")

	// Mark the flags as required
	_ = rootCmd.MarkFlagRequired("definition")
	_ = rootCmd.MarkFlagRequired("gohtml")
	_ = rootCmd.MarkFlagRequired("input")
	_ = rootCmd.MarkFlagRequired("output")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	definitionPath, err := rootCmd.Flags().GetString("definition")
	goHtmlPath, err := rootCmd.Flags().GetString("gohtml")
	inputPath, err = rootCmd.Flags().GetString("input")
	outputPath, err := rootCmd.Flags().GetString("output")

	templateConfig, err := ht.GetHTMLTemplateConfigurationFromFile(definitionPath)
	if err != nil {
		panic(err)
		// log.Fatal(err)
	}

	_ = templateConfig

	result, err := GetJsonFileAsMap(inputPath) // "inputPath/BradHannah.json")
	if err != nil {
		panic(err)
		// log.Fatalf("Error unmarshaling JSON: %v", err)
	}

	for _, definition := range templateConfig.Definitions {
		// loop through and see if any of the inputs are missing
		// if they are missing then we will do a CLI prompt

		answer, found := result[definition.Key]
		if !found {
			answer, err = promptAndWaitForAnswer(definition, goHtmlTemplatePath)
			if err != nil {
				panic(err)
			}
			result[definition.Key] = answer
		}
	}

	tmpl, err := GetHtmlTemplate(goHtmlPath) // "inputPath/BradHannahCoverLetter.gohtml")

	if err != nil {
		panic(err)
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the file when we're done
	defer outputFile.Close()

	if err := tmpl.Execute(outputFile, &result); err != nil {
		panic(err)
	}
}
