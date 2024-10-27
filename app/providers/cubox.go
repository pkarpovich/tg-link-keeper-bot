package providers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jimmysawczuk/recon"
	"io"
	"log"
	"net/http"
	"net/url"
)

type Cubox struct {
	Url     string
	DryMode bool
}

func NewCubox(url string, DryMode bool) *Cubox {
	return &Cubox{Url: url, DryMode: DryMode}
}

const (
	LinkType    = "url"
	TextType    = "text"
	ForwardType = "forward"
)

type Content struct {
	Type  string
	Value string
}

type LinkStoreResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (c *Cubox) SaveLink(content Content) error {
	client := &http.Client{}
	jsonBody, err := prepareRequestBody(content)
	if err != nil {
		return fmt.Errorf("failed to prepare request body: %w", err)
	}

	if c.DryMode {
		log.Printf("[DEBUG] dry mode enabled, skipping link save")
		return nil
	}

	req, err := http.NewRequest("POST", c.Url, jsonBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Printf("[ERROR] failed to close response body: %v", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	bodyStr := string(body)
	log.Printf("[DEBUG] save link response: %s", bodyStr)

	var linkStoreResponse LinkStoreResponse
	if err := json.Unmarshal(body, &linkStoreResponse); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if linkStoreResponse.Code < 0 {
		return errors.New(linkStoreResponse.Message)
	}

	return nil
}

func prepareRequestBody(content Content) (io.Reader, error) {
	values := map[string]string{
		"content": content.Value,
		"type":    "url",
	}

	if content.Type == TextType {
		values["type"] = "memo"
	}

	if content.Type == ForwardType || content.Type == LinkType {
		urlMetadata, err := prepareLinkMetadata(content.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to prepare link metadata: %w", err)
		}
		log.Printf("[DEBUG] url metadata: %+v", urlMetadata)

		values["description"] = urlMetadata.Description
		values["title"] = urlMetadata.Title
		values["content"] = urlMetadata.Url
	}

	jsonValue, err := json.Marshal(values)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal values: %w", err)
	}

	return bytes.NewBuffer(jsonValue), nil
}

func (c *Cubox) PrepareContent(text, msgUrl string) *Content {
	if text == "" {
		return nil
	}

	if _, err := url.ParseRequestURI(text); err == nil {
		return &Content{
			Type:  LinkType,
			Value: text,
		}
	}

	if len(msgUrl) > 0 {
		return &Content{
			Type:  ForwardType,
			Value: msgUrl,
		}
	}

	return &Content{
		Type:  TextType,
		Value: text,
	}
}

type LinkMetadata struct {
	Title       string
	Description string
	Url         string
}

func prepareLinkMetadata(inputURL string) (*LinkMetadata, error) {
	updatedUrl, err := replaceTwitterDomain(inputURL)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare url: %w", err)
	}

	res, err := recon.Parse(updatedUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %w", err)
	}

	return &LinkMetadata{
		Description: res.Description,
		Title:       res.Title,
		Url:         res.URL,
	}, nil
}

func replaceTwitterDomain(inputURL string) (string, error) {
	parsedUrl, err := url.Parse(inputURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse url: %w", err)
	}

	if parsedUrl.Host == "twitter.com" || parsedUrl.Host == "xxxx.com" {
		parsedUrl.Host = "t.fixupx.com"
	}

	return parsedUrl.String(), nil
}
