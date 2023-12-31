package metadata

import (
	"fmt"
	"github.com/jimmysawczuk/recon"
)

type LinkMetadata struct {
	Title       string
	Description string
	Url         string
}

func Prepare(url string) (*LinkMetadata, error) {
	res, err := recon.Parse(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %w", err)
	}

	return &LinkMetadata{
		Description: res.Description,
		Title:       res.Title,
		Url:         res.URL,
	}, nil
}
