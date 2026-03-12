package sse

import (
	"io"
	"log"
	"strings"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SSEHandler struct {
	broker *Broker
}

func NewSSEHandler(b *Broker) *SSEHandler {
	return &SSEHandler{broker: b}
}

func (h *SSEHandler) StreamPublic(c *gin.Context) {

	log.Println("In public")
	clientID := uuid.NewString()
	client := h.broker.AddClient(clientID)
	defer h.broker.RemoveClient(clientID)

	topics := strings.SplitSeq(c.Query("topics"), ",")

	for topic := range topics {

		if topic == "" {
			continue
		}

		h.broker.Subscribe(client, topic)
	}

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	c.Stream(func(w io.Writer) bool {

		select {

		case msg, ok := <-client.Channel:

			if !ok {
				return false
			}

			c.SSEvent(msg.Topic, msg.Content)

			return true

		case <-c.Request.Context().Done():
			return false
		}

	})
}

func (h *SSEHandler) StreamAuth(c *gin.Context) {

	log.Println("In auth")

	clientID := uuid.NewString()
	client := h.broker.AddClient(clientID)
	defer h.broker.RemoveClient(clientID)
	log.Println("opics")
	roles := c.MustGet("roles").([]string)
	userID := c.MustGet("userID").(string)

	topics := strings.SplitSeq(c.Query("topics"), ",")
	log.Println(userID)
	log.Println(topics)
	log.Println(roles)
	log.Println(utils.CanSubscribe(userID, roles, "blog_created_admin"))

	for topic := range topics {
		log.Println(topic)
		if topic == "" {
			continue
		}

		// if !utils.CanSubscribe(userID, roles, topic) {
		// 	c.SSEvent("error", "unauthorized_topic")
		// 	log.Println("Failed")
		// 	return
		// }
		log.Println("In auth sub")
		h.broker.Subscribe(client, topic)
	}

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	c.Stream(func(w io.Writer) bool {
		log.Println("In auth sub1")
		select {
		case msg, ok := <-client.Channel:
			log.Println("In auth sub2")

			if !ok {
				return false
			}

			c.SSEvent(msg.Topic, msg.Content)

			return true

		case <-c.Request.Context().Done():
			return false
		}

	})
}
