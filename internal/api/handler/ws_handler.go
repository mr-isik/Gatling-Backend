package handler

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/mr-isik/gatling-backend/internal/api/ws"
)

type WSHandler struct {
	hub *ws.Hub
}

func NewWSHandler(hub *ws.Hub) *WSHandler {
	return &WSHandler{hub: hub}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow all origins for now. In production, this should be restricted.
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// ServeLiveMetrics handles websocket requests for live metrics
func (h *WSHandler) ServeLiveMetrics(w http.ResponseWriter, r *http.Request) {
	runID := r.PathValue("id")
	if runID == "" {
		http.Error(w, "Run ID is required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	client := ws.NewClient(h.hub, conn, runID, "live")
	client.Register()

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.WritePump()
	go client.ReadPump()
}

// ServeAnomalies handles websocket requests for anomaly events
func (h *WSHandler) ServeAnomalies(w http.ResponseWriter, r *http.Request) {
	runID := r.PathValue("id")
	if runID == "" {
		http.Error(w, "Run ID is required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	client := ws.NewClient(h.hub, conn, runID, "anomalies")
	client.Register()

	go client.WritePump()
	go client.ReadPump()
}
