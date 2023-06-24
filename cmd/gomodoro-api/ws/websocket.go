package ws

import (
	"github.com/gofiber/contrib/websocket"
	"sync"
)

var (
	hubInstance *hub
)

func Start() {
	hubInstance = newHub()
	go hubInstance.run()
}

func Serve(c *websocket.Conn) {
	newClient := NewClient(hubInstance, c, &sync.WaitGroup{})

	hubInstance.register <- newClient

	newClient.startExchange()
}
