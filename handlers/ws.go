package handlers

import (
	"github.com/rs/zerolog"
	"golang.org/x/net/websocket"
)

// This will have a websocket handler.

func WSStd(l zerolog.Logger) websocket.Handler {
	return func(conn *websocket.Conn) {
		// io.Copy(conn, conn)
		message := messageResponse{}
		for {
			if err := websocket.JSON.Receive(conn, &message); err != nil {
				l.Err(err).Msgf("could not receive message")
				_ = conn.Close()
				return
			}

			message = messageResponse{Message: "pong"}
			if err := websocket.JSON.Send(conn, message); err != nil {
				l.Err(err).Msgf("could not send message")
				_ = conn.Close()
				return
			}
		}
	}
}
