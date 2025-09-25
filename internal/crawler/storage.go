package crawler

import (
	"encoding/json"
	"fmt"
	"os"
)

func SaveResults(results []Page, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("cannot create file: %w", err)
	}
	defer file.Close()
	
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	
	if err := encoder.Encode(results); err != nil {
		return fmt.Errorf("cannot encode JSON: %w", err)
	}
	
	return nil
}

func LoadResults(filename string) ([]Page, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %w", err)
	}
	defer file.Close()
	
	var results []Page
	decoder := json.NewDecoder(file)
	
	if err := decoder.Decode(&results); err != nil {
		return nil, fmt.Errorf("cannot decode JSON: %w", err)
	}
	
	return results, nil
}