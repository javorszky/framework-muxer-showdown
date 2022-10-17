package handlers

import (
	"encoding/json"

	"github.com/dgrr/fastws"
	"github.com/rs/zerolog"
)

// This will have a websocket handler.

func WSStd(l zerolog.Logger) func(*fastws.Conn) {
	return func(conn *fastws.Conn) {
		// io.Copy(conn, conn)
		message := messageResponse{}
		var rawMessage []byte
		for {
			_, b, err := conn.ReadMessage(rawMessage)
			if err != nil {
				l.Err(err).Msg("reading message failed")
				_ = conn.Close()
			}
			err = json.Unmarshal(b, &message)
			if err != nil {
				l.Err(err).Msg("converting message to message struct from json failed")
				_ = conn.Close()
			}

			message = messageResponse{Message: "pong"}
			enc, err := json.Marshal(message)
			if err != nil {
				l.Err(err).Msg("converting struct to json failed")
			}
			if _, err := conn.Write(enc); err != nil {
				l.Err(err).Msgf("could not send message")
				_ = conn.Close()
				return
			}
		}
	}
}
