package broker

import (
	"encoding/json"
	"io"
	"sync"

	"github.com/gin-gonic/gin"
)

type Broker struct {
	clients map[chan string]bool
	mutex   sync.Mutex
}

func NewBroker() *Broker {
	return &Broker{
		clients: make(map[chan string]bool),
	}
}

func (b *Broker) AddClient(ch chan string) {
	b.mutex.Lock()
	b.clients[ch] = true
	b.mutex.Unlock()
}

func (b *Broker) RemoveClient(ch chan string) {
	b.mutex.Lock()
	delete(b.clients, ch)
	close(ch)
	b.mutex.Unlock()
}

func (b *Broker) Broadcast(data interface{}) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	bytes, _ := json.Marshal(data)

	for client := range b.clients {
		select {
		case client <- string(bytes):
		default:
			// Skip blocked clients
		}
	}
}

func (b *Broker) ServeSSE(c *gin.Context) {
	clientChan := make(chan string)

	b.AddClient(clientChan)
	defer b.RemoveClient(clientChan)

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	c.Stream(func(w io.Writer) bool {
		if msg, ok := <-clientChan; ok {
			c.SSEvent("admin-notification", msg)
			return true
		}
		return false
	})
}
