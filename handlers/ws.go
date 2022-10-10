package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"golang.org/x/net/websocket"
)

// This will have a websocket handler.

func Ping(l zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		w, r := c.Writer, c.Request

		websocket.Handler(func(conn *websocket.Conn) {
			message := messageResponse{}
			for {
				if err := websocket.JSON.Receive(conn, &message); err != nil {
					l.Err(err).Msg("ws receive message")
					_ = conn.Close()
					return
				}

				message.Message = "pong"
				if err := websocket.JSON.Send(conn, message); err != nil {
					l.Err(err).Msg("ws send message")
					_ = conn.Close()
					return
				}
			}
		}).ServeHTTP(w, r)
	}
}
