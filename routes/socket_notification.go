package routes

import (
	"encoding/json"
	"log"

	"github.com/LitPad/backend/config"
	"github.com/LitPad/backend/database"
	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/choices"
	"github.com/LitPad/backend/schemas"
	"github.com/LitPad/backend/utils"
	"github.com/gofiber/contrib/websocket"
	"gorm.io/gorm"
)

type SocketNotificationSchema struct {
	schemas.NotificationSchema
	Status choices.NotificationStatus `json:"status"`
}

func (s SocketNotificationSchema) Init(notification models.Notification) SocketNotificationSchema {
	s.NotificationSchema = s.NotificationSchema.Init(notification, true)
	return s
}

var notificationObj SocketNotificationSchema
// Function to broadcast a notification data to all connected clients
func broadcastNotificationMessage(db *gorm.DB, mt int, msg []byte) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	for client := range clients {
		user := client.Locals("user").(*models.User)
		if user == nil {
			continue
		}
		json.Unmarshal(msg, &notificationObj)
		// Ensure user is a valid recipient of this notification
		if *notificationObj.ReceiverID == user.ID {
			if err := client.WriteMessage(mt, msg); err != nil {
				log.Println("write:", err)
			}
		}
	}
	if notificationObj.Status == "DELETED" {
		db.Delete(&models.Notification{}, notificationObj.ID)
	}
}

func (ep Endpoint) NotificationSocket(c *websocket.Conn) {
	cfg := config.GetConfig()
	db := database.ConnectDb(cfg, true)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	token := c.Headers("Authorization")

	var (
		mt     int
		msg    []byte
		err    error
		user   *models.User
		secret *string
		errM   *string
	)

	// Validate Auth
	if user, secret, errM = ValidateAuth(db, token); errM != nil {
		ReturnError(c, utils.ERR_INVALID_TOKEN, *errM, 4001)
		return
	}
	// Add the client to the list of connected clients
	c.Locals("user", user)
	AddClient(c)

	// Remove the client from the list when the handler exits
	defer RemoveClient(c)

	for {
		if mt, msg, err = c.ReadMessage(); err != nil {
			ReturnError(c, utils.ERR_INVALID_ENTRY, "Invalid Entry", 4220)
			break
		}

		// Notifications can only be broadcasted from the app using the socket secret
		if secret != nil {
			broadcastNotificationMessage(db, mt, msg)
		} else {
			ReturnError(c, utils.ERR_UNAUTHORIZED_USER, "Not authorized to send data", 4001)
			break
		}
	}
}