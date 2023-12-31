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
	"net/url"
)

type Client struct {
	Url string
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

func NewLinkStoreClient(url string) *Client {
	return &Client{
		Url: url,
	}
}

func (lc *Client) OnMessage(msg bot.Message) error {
	content := prepareContent(msg)

	if err := lc.saveLink(content); err != nil {
		return fmt.Errorf("failed to save link: %w", err)
	}

	return nil
}

func (lc *Client) saveLink(content Content) error {
	client := &http.Client{}
	jsonBody, err := prepareRequestBody(content)
	if err != nil {
		return fmt.Errorf("failed to prepare request body: %w", err)
	}

	req, err := http.NewRequest("POST", lc.Url, jsonBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
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

func prepareRequestBody(content Content) (io.Reader, error) {
	values := map[string]string{
		"content": content.Value,
		"type":    "url",
	}

	if content.Type == TextType {
		values["type"] = "memo"
	}

	if content.Type == ForwardType || content.Type == LinkType {
		urlMetadata, err := metadata.Prepare(content.Value)
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

func prepareContent(msg bot.Message) Content {
	if _, err := url.ParseRequestURI(msg.Text); err == nil {
		return Content{
			Type:  LinkType,
			Value: msg.Text,
		}
	}

	if msg.ForwardFromChat != nil {
		forwardPostUrl := fmt.Sprintf("https://t.me/%s/%d", msg.ForwardFromChat.UserName, msg.ForwardFromMessageID)

		return Content{
			Type:  ForwardType,
			Value: forwardPostUrl,
		}
	}

	return Content{
		Type:  TextType,
		Value: msg.Text,
	}
}
