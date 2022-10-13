package handlers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/rs/zerolog"
)

// This will have a websocket handler.

func WS(l zerolog.Logger) fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		var msg messageResponse

		var (
			mt    int
			recvd []byte
			err   error
		)
		for {
			if mt, recvd, err = c.ReadMessage(); err != nil {
				l.Err(err).Msg("wbesocket read message broke")
				_ = c.Close()
				break
			}

			err = json.Unmarshal(recvd, &msg)
			if err != nil {
				l.Err(err).Msg("parsing incoming message to struct broke")
				_ = c.Close()
				break
			}

			out, err := json.Marshal(messageResponse{Message: "pong"})
			if err != nil {
				_ = c.Close()
				break
			}

			if err = c.WriteMessage(mt, out); err != nil {
				l.Err(err).Msg("sending message back broke")
				_ = c.Close()
				break
			}
		}
	})
}

func WSUpgradeMW() fiber.Handler {
	return func(c *fiber.Ctx) error {

		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	}
}
