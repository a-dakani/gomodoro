package tomodoro

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/url"
	"sync"
	"time"
)

const pingPeriod = 10 * time.Second
const tickPeriod = 1 * time.Second

type WebSocketClient struct {
	configStr string
	ctx       context.Context
	ctxCancel context.CancelFunc
	mu        sync.RWMutex
	conn      *websocket.Conn
	OutChan   chan Message
}

func NewWebSocketClient(teamSlug string) *WebSocketClient {
	wsc := WebSocketClient{}

	wsc.configStr, _ = url.JoinPath(BaseWSURLV1, URLTeamSlug, teamSlug, "ws")

	wsc.ctx, wsc.ctxCancel = context.WithCancel(context.Background())

	wsc.OutChan = make(chan Message, 100)

	return &wsc
}

func (wsc *WebSocketClient) Start() {
	go wsc.listen()
	go wsc.ping()
}

func (wsc *WebSocketClient) Connect() *websocket.Conn {
	wsc.mu.Lock()
	defer wsc.mu.Unlock()
	if wsc.conn != nil {
		return wsc.conn
	}

	ticker := time.NewTicker(tickPeriod)
	defer ticker.Stop()
	for ; ; <-ticker.C {
		select {
		case <-wsc.ctx.Done():
			return nil
		default:
			wsc.eventHandler(Connecting, nil)
			ws, _, err := websocket.DefaultDialer.Dial(wsc.configStr, nil)
			if err != nil {
				continue
			}
			wsc.conn = ws
			return wsc.conn
		}
	}
}

func (wsc *WebSocketClient) listen() {
	wsc.eventHandler(Listening, nil)
	ticker := time.NewTicker(tickPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-wsc.ctx.Done():
			return
		case <-ticker.C:
			for {
				ws := wsc.Connect()
				if ws == nil {
					return
				}
				_, bytMsg, err := ws.ReadMessage()
				if err != nil {
					wsc.closeWs()
					break
				}

				if err != nil {
					fmt.Println(fmt.Sprintf("Cannot unmarshal websocket message got error %s", err.Error()))
					break
				}
				// push messages to handler
				wsc.msgHandler(bytMsg)
			}
		}
	}
}

func (wsc *WebSocketClient) ping() {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			ws := wsc.Connect()
			if ws == nil {
				continue
			}
			if err := wsc.conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(pingPeriod/2)); err != nil {
				wsc.closeWs()
			}
		case <-wsc.ctx.Done():
			return
		}
	}
}

func (wsc *WebSocketClient) eventHandler(messageType MessageType, err error) {
	var m Message
	m.Type = messageType
	if err != nil {
		m.Error = err.Error()
	}
	wsc.OutChan <- m

}

func (wsc *WebSocketClient) msgHandler(msg []byte) {
	var m Message
	_ = json.Unmarshal(msg, &m)
	wsc.OutChan <- m
}

func (wsc *WebSocketClient) Stop() {
	wsc.ctxCancel()
	wsc.closeWs()
	wsc.eventHandler(Terminated, nil)
}

func (wsc *WebSocketClient) closeWs() {
	wsc.mu.Lock()
	if wsc.conn != nil {
		wsc.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		wsc.conn.Close()
		wsc.conn = nil
	}
	wsc.mu.Unlock()
}
