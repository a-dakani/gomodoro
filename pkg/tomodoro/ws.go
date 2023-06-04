package tomodoro

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

type MessageType string

const (
	Tick         MessageType = "tick"
	TimerStopped MessageType = "timerStopped"
	TimerStarted MessageType = "timerStarted"
)

type Message struct {
	Type    MessageType `json:"type"`
	Payload struct {
		Name          string `json:"name"`
		RemainingTime int64  `json:"remainingTime"`
		Team          string `json:"team"`
		Timestamp     int64  `json:"timestamp"`
	} `json:"payload"`
}

const pingPeriod = 10 * time.Second
const tickPeriod = 1 * time.Second

type WebSocketClient struct {
	configStr string
	ctx       context.Context
	ctxCancel context.CancelFunc

	mu     sync.RWMutex
	wsconn *websocket.Conn
}

func NewWebSocketClient(address string) (*WebSocketClient, error) {
	c := WebSocketClient{}

	c.ctx, c.ctxCancel = context.WithCancel(context.Background())
	c.configStr = address

	go c.listen()
	go c.ping()
	return &c, nil
}

func (wsc *WebSocketClient) Connect() (*websocket.Conn, error) {
	wsc.mu.Lock()
	defer wsc.mu.Unlock()
	if wsc.wsconn != nil {
		return wsc.wsconn, nil
	}

	ticker := time.NewTicker(tickPeriod)
	defer ticker.Stop()
	for ; ; <-ticker.C {
		select {
		case <-wsc.ctx.Done():
			return nil, errors.New("context done")
		default:
			ws, _, err := websocket.DefaultDialer.Dial(wsc.configStr, nil)
			if err != nil {
				fmt.Sprintln(fmt.Sprintf("Cannot connect to websocket: %s\nError: %s", wsc.configStr, err.Error()))
				continue
			}
			wsc.wsconn = ws
			return wsc.wsconn, nil
		}
	}
}

func (wsc *WebSocketClient) listen() {
	fmt.Println("listen started")
	ticker := time.NewTicker(tickPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-wsc.ctx.Done():
			return
		case <-ticker.C:
			for {
				ws, err := wsc.Connect()
				if err != nil {
					fmt.Println(fmt.Sprintf("Cannot connect to websocket got error %s", err.Error()))
					break
				}
				if ws == nil {
					return
				}
				_, bytMsg, err := ws.ReadMessage()
				if err != nil {
					fmt.Println(fmt.Sprintf("Cannot read websocket message got error %s", err.Error()))
					wsc.closeWs()
					break
				}

				if err != nil {
					fmt.Println(fmt.Sprintf("Cannot unmarshal websocket message got error %s", err.Error()))
					break
				}
				// push messages to channel
				wsc.msgHandler(bytMsg)
			}
		}
	}
}

func (wsc *WebSocketClient) Stop() {
	wsc.ctxCancel()
	wsc.closeWs()
}

func (wsc *WebSocketClient) closeWs() {
	wsc.mu.Lock()
	if wsc.wsconn != nil {
		wsc.wsconn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		wsc.wsconn.Close()
		wsc.wsconn = nil
	}
	wsc.mu.Unlock()
}

func (wsc *WebSocketClient) ping() {
	fmt.Println("ping pong started")
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			ws, err := wsc.Connect()
			if err != nil {
				fmt.Println(fmt.Sprintf("Cannot connect to websocket got error %s", err.Error()))
				break
			}
			if ws == nil {
				continue
			}
			if err := wsc.wsconn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(pingPeriod/2)); err != nil {
				wsc.closeWs()
			}
		case <-wsc.ctx.Done():
			return
		}
	}
}

func (wsc *WebSocketClient) msgHandler(msg []byte) {
	var m Message
	_ = json.Unmarshal(msg, &m)
	fmt.Printf("received Message: %v\n ", m)
}
