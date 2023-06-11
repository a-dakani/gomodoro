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
	Slug      string
	ctx       context.Context
	ctxCancel context.CancelFunc
	mu        sync.RWMutex
	conn      *websocket.Conn
	OutChan   chan Message
}

func NewWebSocketClient(teamSlug string) *WebSocketClient {
	wsc := WebSocketClient{}
	wsc.Slug = teamSlug
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
				fmt.Println(fmt.Sprintf("Cannot connect to websocket got error %s", err.Error()))
				continue
			}
			wsc.conn = ws
			wsc.eventHandler(Connected, nil)
			return wsc.conn
		}
	}
}

func (wsc *WebSocketClient) listen() {
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
	wsc.eventHandler(Terminating, nil)
	wsc.ctxCancel()
	wsc.closeWs()
	close(wsc.OutChan)
}

func (wsc *WebSocketClient) closeWs() {
	wsc.mu.Lock()
	if wsc.conn != nil {
		err := wsc.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			fmt.Println(fmt.Sprintf("Cannot close websocket got error %s", err.Error()))
			return
		}
		err = wsc.conn.Close()
		if err != nil {
			fmt.Println(fmt.Sprintf("Cannot close websocket got error %s", err.Error()))
			return
		}
		wsc.conn = nil
	}
	wsc.mu.Unlock()
}
