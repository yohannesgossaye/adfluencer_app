package mail

import (
	"bufio"
	"net"
	"strings"
	"testing"
	"time"
)

// fakeSMTP speaks just enough SMTP for one message and returns what it received.
func fakeSMTP(t *testing.T, stall bool) (port string, got chan string) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { ln.Close() })
	got = make(chan string, 1)
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		if stall {
			time.Sleep(2 * time.Second)
			return
		}
		r := bufio.NewReader(conn)
		write := func(s string) { conn.Write([]byte(s + "\r\n")) }
		write("220 fake ESMTP")
		var data strings.Builder
		inData := false
		for {
			line, err := r.ReadString('\n')
			if err != nil {
				return
			}
			if inData {
				if line == ".\r\n" {
					inData = false
					write("250 queued")
					continue
				}
				data.WriteString(line)
				continue
			}
			switch cmd := strings.ToUpper(strings.Fields(line)[0]); cmd {
			case "EHLO":
				write("250-fake")
				write("250 AUTH PLAIN")
			case "AUTH":
				write("235 ok")
			case "MAIL", "RCPT":
				write("250 ok")
			case "DATA":
				inData = true
				write("354 go")
			case "QUIT":
				write("221 bye")
				got <- data.String()
				return
			}
		}
	}()
	_, port, _ = net.SplitHostPort(ln.Addr().String())
	return port, got
}

func TestSendDeliversMessage(t *testing.T) {
	port, got := fakeSMTP(t, false)
	c := New(Config{Host: "127.0.0.1", Port: port, User: "u", Password: "p", From: "noreply@adfluence.test"})
	if err := c.Send("user@example.com", "Hello", "Your code is 123456"); err != nil {
		t.Fatalf("Send: %v", err)
	}
	msg := <-got
	if !strings.Contains(msg, "Subject: Hello") || !strings.Contains(msg, "Your code is 123456") {
		t.Fatalf("unexpected message:\n%s", msg)
	}
}

func TestSendTimesOutOnStalledServer(t *testing.T) {
	old := sendTimeout
	sendTimeout = 300 * time.Millisecond
	t.Cleanup(func() { sendTimeout = old })

	port, _ := fakeSMTP(t, true)
	c := New(Config{Host: "127.0.0.1", Port: port, User: "u", Password: "p", From: "noreply@adfluence.test"})
	start := time.Now()
	if err := c.Send("user@example.com", "Hello", "body"); err == nil {
		t.Fatal("expected a timeout error")
	}
	if d := time.Since(start); d > time.Second {
		t.Fatalf("Send took %v; the timeout did not apply", d)
	}
}

func TestBuildMessageHeaders(t *testing.T) {
	msg := buildMessage("yohannesgossaye73@gmail.com", "creator@example.com", "Problem reported: የቡና ማስታወቂያ", "Hello\nline two")
	head, body, ok := strings.Cut(msg, "\r\n\r\n")
	if !ok {
		t.Fatal("no header/body separator")
	}
	for _, want := range []string{
		`From: "Adfluence" <yohannesgossaye73@gmail.com>`,
		"To: creator@example.com",
		"Subject: =?utf-8?q?",
		"Date: ",
		"@gmail.com>",
		"Content-Type: text/plain; charset=UTF-8",
	} {
		if !strings.Contains(head, want) {
			t.Errorf("headers missing %q:\n%s", want, head)
		}
	}
	if strings.Contains(head, "የቡና") {
		t.Error("subject must be encoded, not raw UTF-8")
	}
	if body != "Hello\r\nline two" {
		t.Errorf("body = %q, want CRLF line endings", body)
	}
	// A From that already has a name is kept as it is.
	if m := buildMessage(`"Adfluence Team" <team@example.com>`, "a@b.c", "s", "b"); !strings.Contains(m, `From: "Adfluence Team" <team@example.com>`) {
		t.Errorf("named From changed: %s", m)
	}
}
