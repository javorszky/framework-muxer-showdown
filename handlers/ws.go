package handlers

import (
	"net/http"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/rs/zerolog"
	"golang.org/x/net/websocket"
)

// This will have a websocket handler.
type msg struct {
	Message string
}

func WS() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		go func() {
			defer conn.Close()

			for {
				msg, op, err := wsutil.ReadClientData(conn)
				if err != nil {
					// handle error
				}
				err = wsutil.WriteServerMessage(conn, op, msg)
				if err != nil {
					// handle error
				}
			}
		}()
	}
}

func WSStd(l zerolog.Logger) websocket.Handler {
	return func(conn *websocket.Conn) {
		// io.Copy(conn, conn)
		message := msg{}
		for {
			if err := websocket.JSON.Receive(conn, &message); err != nil {
				l.Err(err).Msgf("could not receive message")
				_ = conn.Close()
				return
			}

			message = msg{Message: "pong"}
			if err := websocket.JSON.Send(conn, message); err != nil {
				l.Err(err).Msgf("could not send message")
				_ = conn.Close()
				return
			}
		}
	}
}
