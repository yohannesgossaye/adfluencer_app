package realtime

import (
	dto "adinfluence_service_app/internal/constant/model/dto/marketplace"
	"context"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type Client struct {
	UserID string
	Conn   *websocket.Conn
	Send   chan dto.WSOutgoingMessage
	Hub    *Hub
}

func NewClient(userID string, conn *websocket.Conn, hub *Hub) *Client {
	return &Client{
		UserID: userID,
		Conn:   conn,
		Send:   make(chan dto.WSOutgoingMessage, 64),
		Hub:    hub,
	}
}

func (c *Client) WriteLoop(ctx context.Context) {
	defer c.Conn.CloseNow()

	for {
		select {
		case msg, ok := <-c.Send:
			if !ok {
				return
			}

			writeCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			err := wsjson.Write(writeCtx, c.Conn, msg)
			cancel()

			if err != nil {
				return
			}

		case <-ctx.Done():
			return
		}
	}
}
