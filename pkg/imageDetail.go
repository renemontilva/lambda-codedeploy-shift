package handlers

import (
	"fmt"
	"path/filepath"
)

type ImageDetail struct {
	ImageURI string `json:"ImageURI"`
}

func NewImageDetail(folderPath, fileName string) (ImageDetail, error) {
	var imageDetail ImageDetail
	err := convertJSONToType(filepath.Join(folderPath, fileName), &imageDetail)
	if err != nil {
		return imageDetail, fmt.Errorf("error returned from NewImageDetail: %w", err)
	}
	return imageDetail, nil
}
