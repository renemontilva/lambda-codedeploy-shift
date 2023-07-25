package handlers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type ImageDetail struct {
	ImageURI string `json:"ImageURI"`
}

func NewImageDetail(folderPath, fileName string) (ImageDetail, error) {
	imageDetail := ImageDetail{}
	jsonFile, err := os.ReadFile(filepath.Join(folderPath, fileName))
	if err != nil {
		return imageDetail, fmt.Errorf("error reading JSON file: %w", err)
	}
	err = json.Unmarshal(jsonFile, &imageDetail)
	if err != nil {
		return imageDetail, fmt.Errorf("error unmarshalling JSON file: %w", err)
	}
	return imageDetail, nil
}
