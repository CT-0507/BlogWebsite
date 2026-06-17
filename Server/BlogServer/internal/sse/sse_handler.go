package sse

import (
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	SUBSCRIBER_LIMIT = 20
)

type SSEHandler struct {
	broker  *Broker
	counter atomic.Int32
}

func NewSSEHandler(b *Broker) *SSEHandler {
	return &SSEHandler{broker: b}
}

func (h *SSEHandler) tryAcquireConnection() bool {
	for {
		current := h.counter.Load()

		if current >= SUBSCRIBER_LIMIT {
			return false
		}

		if h.counter.CompareAndSwap(current, current+1) {
			return true
		}
	}
}

func (h *SSEHandler) StreamPublic(c *gin.Context) {

	if !h.tryAcquireConnection() {
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
			"error": "SSE connection limit reached",
		})
		return
	}
	defer h.counter.Add(-1)

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

			c.SSEvent(msg.Event, msg.Data)

			return true

		case <-c.Request.Context().Done():
			return false
		}

	})
}

func (h *SSEHandler) StreamAuth(c *gin.Context) {

	if !h.tryAcquireConnection() {
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
			"error": "SSE connection limit reached",
		})
		return
	}
	defer h.counter.Add(-1)

	clientID := uuid.NewString()
	client := h.broker.AddClient(clientID)
	defer h.broker.RemoveClient(clientID)

	roles := c.MustGet("roles").([]string)
	userID := c.MustGet("userID").(string)

	topics := strings.SplitSeq(c.Query("topics"), ",")
	for topic := range topics {
		if topic == "" {
			continue
		}

		if !utils.CanSubscribe(userID, roles, topic) {
			c.SSEvent("error", "unauthorized_topic")
			return
		}
		h.broker.Subscribe(client, topic)
	}

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	ticker := time.NewTicker(30 * time.Second)
	c.Stream(func(w io.Writer) bool {
		select {
		case msg, ok := <-client.Channel:
			if !ok {
				return false
			}

			c.SSEvent(msg.Event, msg.Data)

			return true

		case <-c.Request.Context().Done():
			return false
		case <-ticker.C:
			c.SSEvent("ping", "")
			return true
		}

	})
}
