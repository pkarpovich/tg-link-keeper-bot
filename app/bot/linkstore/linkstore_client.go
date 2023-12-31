package linkstore

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkarpovich/tg-linkding/app/bot"
	"github.com/pkarpovich/tg-linkding/app/bot/metadata"
	"io"
	"log"
	"net/http"
)

type Client struct {
	Url string
}

func NewLinkStoreClient(url string) *Client {
	return &Client{
		Url: url,
	}
}

func (lc *Client) OnMessage(msg bot.Message) error {
	if err := lc.saveLink(msg.Text); err != nil {
		return fmt.Errorf("failed to save link: %w", err)
	}

	return nil
}

func (lc *Client) saveLink(link string) error {
	urlMetadata, err := metadata.Prepare(link)
	if err != nil {
		return fmt.Errorf("failed to prepare link metadata: %w", err)
	}
	log.Printf("[DEBUG] url metadata: %+v", urlMetadata)

	values := map[string]string{
		"description": urlMetadata.Description,
		"title":       urlMetadata.Title,
		"content":     urlMetadata.Url,
		"type":        "url",
	}
	jsonValue, _ := json.Marshal(values)

	req, err := http.NewRequest("POST", lc.Url, bytes.NewBuffer(jsonValue))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	bodyStr := string(body)
	log.Printf("[DEBUG] save link response: %s", bodyStr)

	return nil
}
