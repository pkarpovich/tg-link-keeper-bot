package providers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jimmysawczuk/recon"
	"github.com/pkarpovich/tg-link-keeper-bot/app/bot"
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
	if c.DryMode {
		log.Printf("[DEBUG] dry mode enabled, skipping link save")
		return nil
	}

	client := &http.Client{}
	jsonBody, err := prepareRequestBody(content)
	if err != nil {
		return fmt.Errorf("failed to prepare request body: %w", err)
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

func (c *Cubox) PrepareContent(msg bot.Message) *Content {
	if msg.Text == "" {
		return nil
	}

	if _, err := url.ParseRequestURI(msg.Text); err == nil {
		return &Content{
			Type:  LinkType,
			Value: msg.Text,
		}
	}

	if len(msg.Url) > 0 {
		return &Content{
			Type:  ForwardType,
			Value: msg.Url,
		}
	}

	return &Content{
		Type:  TextType,
		Value: msg.Text,
	}
}

type LinkMetadata struct {
	Title       string
	Description string
	Url         string
}

func prepareLinkMetadata(url string) (*LinkMetadata, error) {
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
