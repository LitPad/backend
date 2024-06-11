package routes

import (
	"encoding/json"
	"net/http"
	"net/url"
	"sync"

	"os"

	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/choices"
	webs "github.com/fasthttp/websocket"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Maintain db & a list of connected clients
var (
	clients      = make(map[*websocket.Conn]bool)
	clientsMutex = &sync.Mutex{}
)

// Function to add a client to the list
func AddClient(c *websocket.Conn) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	clients[c] = true
}

// Function to remove a client from the list
func RemoveClient(c *websocket.Conn) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	delete(clients, c)
}

type ErrorResp struct {
	Status  string             `json:"status"`
	Code    int                `json:"code"`
	Type    string             `json:"type"`
	Message string             `json:"message"`
	Data    *map[string]string `json:"data,omitempty"`
}

func ReturnError(c *websocket.Conn, errType string, message string, code int, dataOpts ...*map[string]string) {
	errorResponse := ErrorResp{Status: "failure", Code: code, Type: errType, Message: message}
	if len(dataOpts) > 0 {
		errorResponse.Data = dataOpts[0]
	}
	jsonResponse, _ := json.Marshal(errorResponse)
	c.WriteMessage(websocket.TextMessage, jsonResponse)
}

func ValidateAuth(db *gorm.DB, token string) (*models.User, *string, *string) {
	var (
		errMsg *string
		secret *string
		user   *models.User
	)
	if len(token) < 1 {
		err := "Auth bearer not set"
		errMsg = &err
	} else if token == cfg.SocketSecret {
		secret = &token
	} else {
		// Get User
		userObj, err := GetUser(token, db)
		if err != nil {
			errMsg = err
		}
		user = userObj
	}
	return user, secret, errMsg
}

func SendNotificationInSocket(fiberCtx *fiber.Ctx, notification models.Notification, statusOpts ...choices.NotificationStatus) error {
	if os.Getenv("ENVIRONMENT") == "TESTING" {
		return nil
	}
	
	// Check if page size is provided as an argument
	status := choices.NS_CREATED
	if len(statusOpts) > 0 {
		status = statusOpts[0]
	}
	webSocketScheme := "ws://"
	if fiberCtx.Secure() {
		webSocketScheme = "wss://"
	}
	uri := webSocketScheme + fiberCtx.Hostname() + "/api/v1/ws/notifications/"
	notificationData := SocketNotificationSchema{Status: status}
	if status == choices.NS_CREATED {
		notificationData = SocketNotificationSchema{
			Status:             status,
		}.Init(notification)
	}

	// Connect to the WebSocket server
	u, err := url.Parse(uri)
	if err != nil {
		return err
	}

	headers := make(http.Header)
	headers.Add("Authorization", cfg.SocketSecret)
	conn, _, err := webs.DefaultDialer.Dial(u.String(), headers)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Marshal the notification data to JSON
	data, err := json.Marshal(notificationData)
	if err != nil {
		return err
	}

	// Send the notification to the WebSocket server
	err = conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		return err
	}

	// Close the WebSocket connection
	return conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}
