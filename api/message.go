package api

import (
	"encoding/json"
	"errors"
)

const (
	MSG_CFG_REQ = "cfgreq" // config payload
	MSG_CFG_RES = "cfgres" // no payload
	MSG_GET_REQ = "getreq" // no payload
	MSG_GET_RES = "getres" // metrics payload
	MSG_PUT_REQ = "putreq" // metrics payload
	MSG_PUT_RES = "putres" // no payload
	MSG_ERROR   = "error"  // error payload
)

type Config map[string]interface{}

// Basic message type.
type Message struct {
	Type    string
	Payload interface{}
}

func (m *Message) MarshalJSON() (raw []byte, err error) {
	switch m.Type {
	case MSG_CFG_REQ:
		msg := struct {
			Type    string
			Payload Config
		} {
			m.Type,
			Config(m.Payload.(map[string]interface{})),
		}
		raw, err = json.Marshal(&msg)
	case MSG_GET_RES:
		fallthrough
	case MSG_PUT_REQ:
		msg := struct {
			Type    string
			Payload Metrics
		} {
			m.Type,
			m.Payload.(Metrics),
		}
		raw, err = json.Marshal(&msg)
	case MSG_ERROR:
		msg := struct {
			Type    string
			Payload string
		} {
			m.Type,
			(m.Payload.(error)).Error(),
		}
		raw, err = json.Marshal(&msg)
	default:
		msg := struct {
			Type    string
			Payload interface{}
		} {
			m.Type,
			m.Payload,
		}
		raw, err = json.Marshal(&msg)
	}
	return
}

func (m *Message) UnmarshalJSON(raw []byte) (err error) {
	var msgType struct {
		Type string
	}
	if err = json.Unmarshal(raw, &msgType); err != nil {
		return
	}

	switch msgType.Type {
	case MSG_CFG_REQ:
		var msg struct {
			Type    string
			Payload Config
		}
		if err = json.Unmarshal(raw, &msg); err == nil {
			m.Type = msg.Type
			m.Payload = msg.Payload
		}
	case MSG_GET_RES:
		fallthrough
	case MSG_PUT_REQ:
		var msg struct {
			Type    string
			Payload Metrics
		}
		if err = json.Unmarshal(raw, &msg); err == nil {
			m.Type = msg.Type
			m.Payload = msg.Payload
		}
	case MSG_ERROR:
		var msg struct {
			Type    string
			Payload string
		}
		if err = json.Unmarshal(raw, &msg); err == nil {
			m.Type = msg.Type
			m.Payload = errors.New(msg.Payload)
		}
	default:
		var msg struct {
			Type    string
			Payload interface{}
		}
		if err = json.Unmarshal(raw, &msg); err == nil {
			m.Type = msg.Type
			m.Payload = msg.Payload
		}
	}
	return
}

// Return the payload as an error.
func (m *Message) Error() string {
	if m.Type != MSG_ERROR {
		return ""
	}
	return (m.Payload.(error)).Error()
}
