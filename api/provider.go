package api

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
)

type Handler interface {
	Configure(Config) error
	Get() (Metrics, error)
	Put(Metrics) error
}

func RunProvider(handler Handler) {
	encoder := json.NewEncoder(os.Stdout)
	decoder := json.NewDecoder(os.Stdin)

	for {
		message := Message{}
		if err := decoder.Decode(&message); err != nil {
			if err != io.EOF {
				log.Printf("failed to decode request: %s\n", err)
			}
			break
		}

		respond := func(msg string, payload interface{}) {
			encoder.Encode(&Message{Type: msg, Payload: payload})
		}

		if message.Type == MSG_CFG_REQ {
			if err := handler.Configure(message.Payload.(Config)); err == nil {
				respond(MSG_CFG_RES, nil)
			} else {
				respond(MSG_ERROR, err)
			}
		} else if message.Type == MSG_GET_REQ {
			if payload, err := handler.Get(); err == nil {
				respond(MSG_GET_RES, payload)
			} else {
				respond(MSG_ERROR, err)
			}
		} else if message.Type == MSG_PUT_REQ {
			if err := handler.Put(message.Payload.(Metrics)); err == nil {
				respond(MSG_PUT_RES, nil)
			} else {
				respond(MSG_ERROR, err)
			}
		} else {
			err := errors.New("invalid message type")
			respond(MSG_ERROR, err)
		}
	}
}
