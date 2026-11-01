package realtime

import (
	dto "adinfluence_service_app/internal/constant/model/dto/marketplace"
	"context"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

const (
	sendBuffer   = 64
	writeTimeout = 10 * time.Second
	// pingInterval keeps idle connections alive through proxies (which often
	// drop sockets silent for ~60s) and detects dead peers.
	pingInterval = 25 * time.Second
	pingTimeout  = 10 * time.Second
)

// Close codes sent to clients. 4000-4999 are application-defined; clients
// should reconnect (after refreshing their token for 4001).
const (
	CloseTokenExpired       websocket.StatusCode = 4001
	CloseAccountUnavailable websocket.StatusCode = 4003
)

type Client struct {
	UserID string
	Conn   *websocket.Conn
	Send   chan dto.WSOutgoingMessage

	closeOnce sync.Once
}

func NewClient(userID string, conn *websocket.Conn) *Client {
	return &Client{
		UserID: userID,
		Conn:   conn,
		Send:   make(chan dto.WSOutgoingMessage, sendBuffer),
	}
}

// TrySend queues a message without blocking. It returns false when the
// client is not keeping up.
func (c *Client) TrySend(msg dto.WSOutgoingMessage) bool {
	select {
	case c.Send <- msg:
		return true
	default:
		return false
	}
}

// Close sends a close frame with the given code once; the read loop then
// ends and the handler unregisters the client.
func (c *Client) Close(code websocket.StatusCode, reason string) {
	c.closeOnce.Do(func() {
		go c.Conn.Close(code, reason)
	})
}

// WriteLoop delivers queued messages and pings the peer until the
// connection ends.
func (c *Client) WriteLoop(ctx context.Context) {
	defer c.Conn.CloseNow()

	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()

	for {
		select {
		case msg, ok := <-c.Send:
			if !ok {
				return
			}
			writeCtx, cancel := context.WithTimeout(ctx, writeTimeout)
			err := wsjson.Write(writeCtx, c.Conn, msg)
			cancel()
			if err != nil {
				return
			}

		case <-ticker.C:
			pingCtx, cancel := context.WithTimeout(ctx, pingTimeout)
			err := c.Conn.Ping(pingCtx)
			cancel()
			if err != nil {
				return
			}

		case <-ctx.Done():
			return
		}
	}
}
