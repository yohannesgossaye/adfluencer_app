package mail

import (
	"fmt"
	"net/smtp"
	"strings"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	From     string
}

type Client struct {
	cfg Config
}

func New(cfg Config) *Client {
	return &Client{cfg: cfg}
}

func (c *Client) Enabled() bool {
	return c.cfg.Host != "" && c.cfg.User != "" && c.cfg.Password != ""
}

func (c *Client) Send(to, subject, body string) error {
	if !c.Enabled() {
		return fmt.Errorf("smtp is not configured")
	}

	addr := fmt.Sprintf("%s:%s", c.cfg.Host, c.cfg.Port)
	auth := smtp.PlainAuth("", c.cfg.User, c.cfg.Password, c.cfg.Host)

	msg := buildMessage(c.cfg.From, to, subject, body)
	return smtp.SendMail(addr, auth, c.cfg.From, []string{to}, []byte(msg))
}

func buildMessage(from, to, subject, body string) string {
	var b strings.Builder
	b.WriteString("From: ")
	b.WriteString(from)
	b.WriteString("\r\nTo: ")
	b.WriteString(to)
	b.WriteString("\r\nSubject: ")
	b.WriteString(subject)
	b.WriteString("\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n")
	b.WriteString(body)
	return b.String()
}
