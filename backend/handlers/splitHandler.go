package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dakshin-2107/BillSplitterLite/backend/logger"
	"github.com/dakshin-2107/BillSplitterLite/backend/managers"
	"github.com/dakshin-2107/BillSplitterLite/backend/models"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type SplitHandler struct {
	sessionManger *managers.SessionManager
	RouteHandlerBase
}

func (sp *SplitHandler) Init(logger logger.ILogger, modelHelper models.ISplitModelHelper, sessionManger managers.ISessionManager) {
	sp.RouteHandlerBase.BaseInit(logger, modelHelper, nil)

	sp.pattern = `/ping`
	sp.method = "GET"
	sp.group = ""
	sp.handler = sp.Handle
}

func (sp *SplitHandler) Handle(ctx *gin.Context) {

}

func testSplitHandler() IRouteHandlerBase {
	return &SplitHandler{}
}

// handleWebSocket handles incoming HTTP requests for WebSocket upgrade.
func handleWebSocket(w http.ResponseWriter, r *http.Request) {

	// Ideally but we can take it up later
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header is required", http.StatusUnauthorized)
		log.Printf("Rejected connection: Authorization header missing from %s", r.RemoteAddr)
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		http.Error(w, "Invalid Authorization header format. Expected 'Bearer <token>'", http.StatusUnauthorized)
		log.Printf("Rejected connection: Invalid Authorization header format '%s'", authHeader)
		return
	}

	sessionID := parts[1]

	log.Printf("Attempting to upgrade connection for token (sessionID): %s from %s", sessionID, r.RemoteAddr)
	if len(sessionID) < 10 || len(sessionID) > 200 { // Example length validation
		http.Error(w, "Invalid token length", http.StatusUnauthorized)
		log.Printf("Rejected connection: Invalid token length '%s'", sessionID)
		return
	}

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection for sessionID %s: %v", sessionID, err)
		return
	}

	// Ensure the WebSocket connection is closed when the function exits.
	defer func() {
		log.Printf("Closing WebSocket connection for sessionID %s from %s", sessionID, conn.RemoteAddr().String())
		conn.Close()
	}()

	log.Printf("New WebSocket connection established for sessionID: %s from %s", sessionID, conn.RemoteAddr().String())

	// 4. Loop indefinitely to read and write messages over the WebSocket.
	for {
		// ReadMessage reads a message from the WebSocket connection.
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			// Handle WebSocket specific close errors or other read errors.
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Printf("WebSocket client for sessionID %s disconnected gracefully.", sessionID)
			} else {
				log.Printf("Error reading WebSocket message for sessionID %s: %v", sessionID, err)
			}
			return // Exit the loop and close the connection
		}

		// Log the received message, including the sessionID.
		log.Printf("Session %s received (Type: %d): %s", sessionID, messageType, string(p))

		// Prepare a response.
		response := fmt.Sprintf("Server received from session %s: '%s' at %s", sessionID, string(p), time.Now().Format(time.RFC3339))

		// Write the response back to the client.
		err = conn.WriteMessage(messageType, []byte(response))
		if err != nil {
			log.Printf("Error writing WebSocket message to session %s: %v", sessionID, err)
			return // Exit the loop and close the connection
		}
	}
}
