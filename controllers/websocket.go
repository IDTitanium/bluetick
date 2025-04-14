package controllers

import (
	"bluetick/database"
	"bluetick/models"
	"bluetick/utils"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins (adjust for production)
	},
}

// Store active connections: userID -> *websocket.Conn
var clients = make(map[uint64]*websocket.Conn)
var usersOnline = make(map[uint64]bool) // Track user online status
var mu sync.Mutex

// HandleWebSocket manages WebSocket connections
func HandleWebSocket(c *gin.Context) {
	tokenString := c.GetHeader("Sec-WebSocket-Protocol") // Read token from header
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		return
	}

	tokenParts := string(tokenString)

	//TODO: Split token into parts to remove the Bearer portion

	splitToken := strings.Split(tokenParts, " ")

	token := splitToken[1]

	fmt.Println("Token received from websocket", token)

	utils.Log("Token received from websocket, " + tokenString)

	// Validate the token
	_, claims, err := utils.ValidateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	userID := uint64(claims["user_id"].(float64))
	fmt.Printf("User %d connected via WebSocket\n", userID)

	// Upgrade to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	// Store the connection
	mu.Lock()
	clients[userID] = conn
	mu.Unlock()

	// Notify users about online status
	broadcastUserStatus(userID, true)

	// Listen for messages
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("User %d disconnected: %v\n", userID, err)
			mu.Lock()
			delete(clients, userID)
			delete(usersOnline, userID)
			mu.Unlock()
			broadcastUserStatus(userID, false)
			break
		}
		fmt.Printf("Message from user %d: %s\n", userID, msg)

		// Example: Extract receiver and forward the message
		// For now, let's expect the message format: "receiverID:message"
		handleIncomingMessage(userID, msg)
	}
}

// handleIncomingMessage forwards messages between users
func handleIncomingMessage(senderID uint64, msg []byte) {
	// Assuming message format: "receiverID:message"
	parts := string(msg)

	// Split at the first colon
	colonIndex := strings.Index(parts, ":")
	if colonIndex == -1 {
		fmt.Println("Invalid message format")
		return
	}

	receiverIDStr := parts[:colonIndex]
	content := parts[colonIndex+1:]

	receiverID, err := strconv.ParseUint(strings.TrimSpace(receiverIDStr), 10, 64)
	if err != nil {
		fmt.Println("Invalid receiver ID")
		return
	}

	// Store message in DB
	message := models.Message{
		SenderID:   senderID,
		ReceiverID: uint64(receiverID),
		Content:    strings.TrimSpace(content),
	}
	database.DB.Create(&message)

	mu.Lock()
	defer mu.Unlock()

	// Check if receiver is connected
	if conn, ok := clients[receiverID]; ok {
		response := fmt.Sprintf("From %d: %s", senderID, content)
		conn.WriteMessage(websocket.TextMessage, []byte(response))
	} else {
		fmt.Printf("User %d is offline. Message not delivered.\n", receiverID)
	}
}

func broadcastUserStatus(userID uint64, online bool) {
	mu.Lock()
	defer mu.Unlock()

	status := "offline"
	if online {
		status = "online"
	}

	message := fmt.Sprintf("User %d is now %s", userID, status)
	for _, conn := range clients {
		conn.WriteMessage(websocket.TextMessage, []byte(message))
	}
}
