package stream_chat

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type ConnectRequest struct {
	ServerDeterminesID bool `json:"server_determines_connection_id"`
	UserDetails        User `json:"user_details" validate:"required,dive"`
}

type WebsocketConn struct {
	id      string // from hello message
	handler EventHandler
	url     string

	context.Context
	context.CancelFunc

	net.Conn
}

func NewWebsocketConn(url string, handler EventHandler) *WebsocketConn {
	ctx, cancel := context.WithCancel(context.Background())
	return &WebsocketConn{
		Context:    ctx,
		CancelFunc: cancel,
		url:        url,
		Conn:       nil,
		handler:    handler,
	}
}

func (wsConn *WebsocketConn) Dial() error {
	conn, _, _, err := ws.DefaultDialer.Dial(wsConn.Context, wsConn.url)
	if err != nil {
		return err
	}
	wsConn.Conn = conn

	var buf bytes.Buffer
	if err := wsConn.readEvent(&buf); err != nil {
		return err
	}

	var event Event
	if err := json.NewDecoder(&buf).Decode(&event); err != nil {
		return err
	}
	wsConn.id = event.ConnectionID
	go wsConn.readLoop()
	go wsConn.writeLoop()
	return nil
}

func (wsConn *WebsocketConn) ID() string {
	return wsConn.id
}

func (wsConn *WebsocketConn) readLoop() error {
	var buf bytes.Buffer
	for {
		select {
		case <-wsConn.Done():
			return wsConn.Err()
		default:
		}
		buf.Reset()
		if err := wsConn.readEvent(&buf); err != nil {
			return err
		}
		var event Event
		if err := json.NewDecoder(&buf).Decode(&event); err != nil {
			log.Println(err)
			return err
		}
		switch event.Type {
		case EventHealthCheck:

		default:
			wsConn.handler(&event)
		}
	}
}

var ping = []byte("ping")

func (wsConn *WebsocketConn) writeLoop() error {
	for {
		select {
		case <-wsConn.Done():
			return wsConn.Err()
		default:
		}

		if err := wsConn.SetWriteDeadline(time.Now().Add(time.Second * 8)); err != nil {
			return err
		}
		if err := wsutil.WriteClientBinary(wsConn, ping); err != nil {
			return err
		}
		time.Sleep(time.Second * 28)
	}
}

func (wsConn *WebsocketConn) readEvent(buffer *bytes.Buffer) error {
	if err := wsConn.SetReadDeadline(time.Now().Add(time.Second * 35)); err != nil {
		return err
	}
	controlHandler := wsutil.ControlFrameHandler(wsConn, ws.StateClientSide)

	rd := wsutil.Reader{
		Source:          wsConn,
		State:           ws.StateClientSide,
		CheckUTF8:       true,
		SkipHeaderCheck: false,
		OnIntermediate:  controlHandler,
	}

	for {
		hdr, err := rd.NextFrame()
		if err != nil {
			return err
		}
		if hdr.OpCode.IsControl() {
			if err := controlHandler(hdr, &rd); err != nil {
				return err
			}
			continue
		}
		if hdr.OpCode&ws.OpText == 0 {
			if err := rd.Discard(); err != nil {
				return err
			}
			continue
		}
		_, err = buffer.ReadFrom(&rd)
		return err
	}
}
