package ws

import (
	"bytes"
	"github.com/gofiber/contrib/websocket"
	"log"
	"sync"
	"time"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// TODO: Implement channels/groups for gomodoros

type client struct {
	hub  *hub
	conn *websocket.Conn

	send chan []byte

	wg *sync.WaitGroup
}

func NewClient(hub *hub, conn *websocket.Conn, wg *sync.WaitGroup) *client {
	return &client{
		wg:   wg,
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}
}

func (c *client) startExchange() {
	c.wg.Add(2)
	go c.read()
	go c.write()
	c.wg.Wait()
}

func (c *client) read() {
	defer func() {
		c.hub.unregister <- c
		c.wg.Done()

		err := c.conn.Close()
		if err != nil {
			log.Printf("Error Closing Websocket Connection: %v", err)
			return
		}
	}()
	c.conn.SetReadLimit(maxMessageSize)

	err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		log.Printf("Error Setting Read Deadline : %v", err)
		return
	}

	c.conn.SetPongHandler(func(string) error {
		err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			log.Printf("Error Setting Read Deadline : %v", err)
			return err
		}
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			// if the error is not a normal closure, log it
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
				log.Printf("Unexpected Websocket Close Error: %v", err)
			}

			break
		}

		message = bytes.TrimSpace(bytes.ReplaceAll(message, newline, space))

		c.hub.broadcast <- message
	}
}

func (c *client) write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.wg.Done()

		err := c.conn.Close()
		if err != nil {
			log.Printf("Error Closing Websocket Connection: %v", err)
			return
		}
	}()

	for {
		select {
		case message, ok := <-c.send:
			err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Printf("Error Setting Write Deadline: %v", err)
				return
			}

			if !ok {
				err := c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					log.Printf("Error Setting Write Deadline: %v", err)
					return
				}

				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			// Todo extract writer to a function
			_, err = w.Write(message)
			if err != nil {
				log.Printf("Error Writing Message: %v", err)
				return
			}

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				_, err := w.Write(newline)
				if err != nil {
					log.Printf("Error Writing Message: %v", err)
					return
				}

				_, err = w.Write(<-c.send)
				if err != nil {
					log.Printf("Error Writing Message: %v", err)
					return
				}
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Printf("Error Setting Write Deadline: %v", err)
				return
			}

			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
