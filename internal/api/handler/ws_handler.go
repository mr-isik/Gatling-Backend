package handler

import (
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/mr-isik/gatling-backend/internal/api/ws"
)

type WSHandler struct {
	hub *ws.Hub
}

func NewWSHandler(hub *ws.Hub) *WSHandler {
	return &WSHandler{hub: hub}
}

// ServeLiveMetrics handles websocket requests for live metrics
func (h *WSHandler) ServeLiveMetrics(c *websocket.Conn) {
	runID := c.Params("id")
	if runID == "" {
		log.Println("Run ID is required")
		return
	}

	client := ws.NewClient(h.hub, c, runID, "live")
	client.Register()

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.WritePump()
	client.ReadPump() // Keep the connection open and read messages in the same goroutine to prevent early closure
}

// ServeAnomalies handles websocket requests for anomaly events
func (h *WSHandler) ServeAnomalies(c *websocket.Conn) {
	runID := c.Params("id")
	if runID == "" {
		log.Println("Run ID is required")
		return
	}

	client := ws.NewClient(h.hub, c, runID, "anomalies")
	client.Register()

	go client.WritePump()
	client.ReadPump() // Keep the connection open and read messages in the same goroutine to prevent early closure
}
