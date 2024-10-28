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
	Token   string
	DryMode bool
}

func NewCubox(url, token string, DryMode bool) *Cubox {
	return &Cubox{Url: url, Token: token, DryMode: DryMode}
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
	jsonBody, err := c.prepareRequestBody(content)
	if err != nil {
		return err
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

var ErrDuplicatedLink = errors.New("link already saved")

func (c *Cubox) prepareRequestBody(content Content) (io.Reader, error) {
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

		if c.checkForDuplicatedLink(urlMetadata.Title) {
			return nil, ErrDuplicatedLink
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

type SearchResponse struct {
	Items []struct {
		Id  string `json:"userSearchEngineID"`
		Url string `json:"targetURL"`
	} `json:"data"`
}

func (c *Cubox) checkForDuplicatedLink(link string) bool {
	httpClient := &http.Client{}

	searchURL := buildSearchURL(link)
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		log.Printf("[ERROR] failed to create request: %v", err)
		return false
	}
	req.Header.Add("Authorization", c.Token)
	bodyResp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("[ERROR] failed to send request: %v", err)
		return false
	}
	defer func() {
		if err := bodyResp.Body.Close(); err != nil {
			fmt.Printf("[ERROR] Error while closing response body: %v", err)
		}
	}()

	var searchResp SearchResponse
	body, err := io.ReadAll(bodyResp.Body)
	if err != nil {
		log.Printf("[ERROR] failed to read response body: %v", err)
		return false
	}

	if err := json.Unmarshal(body, &searchResp); err != nil {
		log.Printf("[ERROR] failed to unmarshal response: %v", err)
		return false
	}

	if len(searchResp.Items) > 0 {
		return true
	}

	return false
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

func buildSearchURL(query string) string {
	encodedQuery := url.QueryEscape(query)
	return fmt.Sprintf("https://cubox.cc/c/api/search?page=1&pageSize=50&keyword=%s&filters=&archiving=false", encodedQuery)
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
