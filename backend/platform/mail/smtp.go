package mail

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"mime"
	"net"
	netmail "net/mail"
	"net/smtp"
	"strings"
	"time"
)

// sendTimeout caps one whole SMTP conversation; net/smtp.SendMail has no
// timeout of its own and can hang forever on a stalled server.
var sendTimeout = 20 * time.Second

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

	addr := net.JoinHostPort(c.cfg.Host, c.cfg.Port)
	auth := smtp.PlainAuth("", c.cfg.User, c.cfg.Password, c.cfg.Host)
	msg := buildMessage(c.cfg.From, to, subject, body)

	conn, err := net.DialTimeout("tcp", addr, sendTimeout)
	if err != nil {
		return fmt.Errorf("smtp dial: %w", err)
	}
	if err := conn.SetDeadline(time.Now().Add(sendTimeout)); err != nil {
		conn.Close()
		return err
	}

	client, err := smtp.NewClient(conn, c.cfg.Host)
	if err != nil {
		conn.Close()
		return fmt.Errorf("smtp handshake: %w", err)
	}
	defer client.Close()

	// Same sequence as smtp.SendMail, on a connection with a deadline.
	if ok, _ := client.Extension("STARTTLS"); ok {
		if err := client.StartTLS(&tls.Config{ServerName: c.cfg.Host}); err != nil {
			return fmt.Errorf("smtp starttls: %w", err)
		}
	}
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("smtp auth: %w", err)
	}
	if err := client.Mail(c.cfg.From); err != nil {
		return err
	}
	if err := client.Rcpt(to); err != nil {
		return err
	}
	w, err := client.Data()
	if err != nil {
		return err
	}
	if _, err := w.Write([]byte(msg)); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	return client.Quit()
}

// buildMessage writes the headers mail providers expect (Date, Message-ID,
// a sender name) so messages are not filed as spam, and encodes the subject
// so non-English titles, such as Amharic contract names, arrive intact.
func buildMessage(from, to, subject, body string) string {
	sender := from
	if addr, err := netmail.ParseAddress(from); err == nil && addr.Name == "" {
		addr.Name = "Adfluence"
		sender = addr.String()
	}
	domain := "adfluence.local"
	if at := strings.LastIndex(from, "@"); at >= 0 {
		domain = strings.Trim(from[at+1:], "> ")
	}
	id := make([]byte, 12)
	_, _ = rand.Read(id)

	var b strings.Builder
	b.WriteString("From: ")
	b.WriteString(sender)
	b.WriteString("\r\n")
	b.WriteString("To: ")
	b.WriteString(to)
	b.WriteString("\r\n")
	b.WriteString("Subject: ")
	b.WriteString(mime.QEncoding.Encode("utf-8", subject))
	b.WriteString("\r\n")
	b.WriteString("Date: ")
	b.WriteString(time.Now().Format(time.RFC1123Z))
	b.WriteString("\r\n")
	b.WriteString("Message-ID: <")
	b.WriteString(hex.EncodeToString(id))
	b.WriteString("@")
	b.WriteString(domain)
	b.WriteString(">\r\n")
	b.WriteString("MIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\nContent-Transfer-Encoding: 8bit\r\n\r\n")
	b.WriteString(strings.ReplaceAll(strings.ReplaceAll(body, "\r\n", "\n"), "\n", "\r\n"))
	return b.String()
}
