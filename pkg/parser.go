package handlers

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Converts a YAML file to Struct Type.
func convertYAMLToType(file string, object any) error {
	yamlFile, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("error reading YAML file: %w", err)
	}
	err = yaml.Unmarshal(yamlFile, object)
	if err != nil {
		return fmt.Errorf("error unmarshalling YAML file: %w", err)
	}
	return nil
}

// Converts a JSON file to Struct Type.
func convertJSONToType(file string, object any) error {
	jsonFile, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("error reading JSON file: %w", err)
	}
	err = json.Unmarshal(jsonFile, object)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON file: %w", err)
	}
	return nil
}
