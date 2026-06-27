package sms

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Config struct {
	Username string
	APIKey   string
	From     string
}

type Client struct {
	cfg Config
}

func New(cfg Config) *Client {
	return &Client{cfg: cfg}
}

func (c *Client) Enabled() bool {
	return c.cfg.Username != "" && c.cfg.APIKey != ""
}

func (c *Client) Send(to, message string) error {
	if !c.Enabled() {
		return fmt.Errorf("sms is not configured")
	}

	form := url.Values{}
	form.Set("username", c.cfg.Username)
	form.Set("to", to)
	form.Set("message", message)
	if c.cfg.From != "" {
		form.Set("from", c.cfg.From)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.africastalking.com/version1/messaging",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return fmt.Errorf("build sms request: %w", err)
	}
	req.Header.Set("apiKey", c.cfg.APIKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("send sms: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if resp.StatusCode >= 400 {
		return fmt.Errorf("sms provider returned %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	return nil
}
