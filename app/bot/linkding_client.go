package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type LinkdingClient struct {
	Url   string
	Token string
}

func NewLinkdingClient(token string, url string) *LinkdingClient {
	return &LinkdingClient{
		Token: token,
		Url:   url,
	}
}

func (lc *LinkdingClient) OnMessage(msg Message) error {
	if err := lc.saveLink(msg.Text); err != nil {
		return fmt.Errorf("failed to save link: %w", err)
	}

	return nil
}

func (lc *LinkdingClient) saveLink(link string) error {
	url := fmt.Sprintf("%s/api/bookmarks/", lc.Url)
	values := map[string]string{
		"url": link,
	}
	jsonValue, _ := json.Marshal(values)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", lc.Token))

	client := &http.Client{}
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
