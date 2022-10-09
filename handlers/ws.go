package handlers

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

// This will have a websocket handler.

type msg struct {
	Message string `json:"message"`
}

func Ping() echo.HandlerFunc {
	return func(c echo.Context) error {
		websocket.Handler(func(conn *websocket.Conn) {
			message := msg{}
			for {
				if err := websocket.JSON.Receive(conn, &message); err != nil {
					c.Logger().Errorf("ws receive message: got error: %s", err)
					_ = conn.Close()
					return
				}

				message = msg{Message: "pong"}
				if err := websocket.JSON.Send(conn, message); err != nil {
					c.Logger().Errorf("ws send message: got error: %s", err)
					_ = conn.Close()
					return
				}
			}
		}).ServeHTTP(c.Response(), c.Request())
		return nil
	}
}
