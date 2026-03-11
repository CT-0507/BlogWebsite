package sse

import (
	"log"
	"sync"
)

type Message struct {
	Event string
	Data  interface{}
}

type Client struct {
	ID      string
	Channel chan Message
	Topics  map[string]bool
}

type Broker struct {
	topics  map[string]map[string]*Client
	clients map[string]*Client
	mu      sync.RWMutex
}

func NewBroker() *Broker {
	return &Broker{
		topics:  make(map[string]map[string]*Client),
		clients: make(map[string]*Client),
	}
}

func (b *Broker) AddClient(id string) *Client {
	b.mu.Lock()
	defer b.mu.Unlock()

	client := &Client{
		ID:      id,
		Channel: make(chan Message, 10),
		Topics:  make(map[string]bool),
	}

	b.clients[id] = client
	return client
}

func (b *Broker) RemoveClient(id string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	client, ok := b.clients[id]
	if !ok {
		return
	}

	for topic := range client.Topics {
		delete(b.topics[topic], id)
	}

	close(client.Channel)
	delete(b.clients, id)
}

func (b *Broker) Subscribe(client *Client, topic string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.topics[topic] == nil {
		b.topics[topic] = make(map[string]*Client)
	}

	b.topics[topic][client.ID] = client
	client.Topics[topic] = true
}

func (b *Broker) Publish(topic, event string, data interface{}) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	subs := b.topics[topic]

	msg := Message{
		Event: event,
		Data:  data,
	}
	log.Println("In publishing")
	log.Println(subs)
	for _, client := range subs {

		select {
		case client.Channel <- msg:
		default:
			// skip slow client
		}

	}
}
