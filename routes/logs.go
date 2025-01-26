package routes

import (
	"time"

	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"
)

// Renders the login page
func (ep Endpoint) RenderLogsLogin(c *fiber.Ctx) error {
	sess := Session(c, ep.Store)
	sessionData := map[string]interface{}{
        "error": sess.Get("error"),
    }
	return c.Render("templates/logs/login.html", sessionData)
}

// Handles login form submission
func (ep Endpoint) HandleLogsLogin(c *fiber.Ctx) error {
	db := ep.DB
	session := Session(c, ep.Store)
	email := c.FormValue("email")
	password := c.FormValue("password")
	user := userManager.GetByUsername(db, email)
	if user == nil {
		session.Set("error", "Invalid email or password!")
		session.Save()
		return c.Redirect("/logs/login")
	}	
	if user.Password != password && !utils.CheckPasswordHash(password, user.Password) {
		session.Set("error", "Invalid email or password!")
		session.Save()
		return c.Redirect("/logs/login")
	}
	if !user.IsStaff {
		session.Set("error", "Unauthorized user!")
		session.Save()
		return c.Redirect("/logs/login")
	}
	
	accessToken := GenerateAccessToken(user.ID)
	user.Access = &accessToken
	// Hash password
	if user.Password == password {
		user.Password = utils.HashPassword(password)
		db.Save(&user)
	}
	session.Set("access", accessToken)
	session.Set("success", "Logged in successfully!")
	session.Set("error", nil)
	session.Save()
	return c.Redirect("/logs")
}

// Renders the logs page with optional filters
func (ep Endpoint) RenderLogs(c *fiber.Ctx) error {
	// Verify page access
	session := Session(c, ep.Store)
	user := VerifyAccess(session, ep.DB)
	if user == nil {
		return c.Redirect("/logs/login")
	}

	db := ep.DB
	logs := []models.Log{}

	// Get query parameters for filtering
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")
	method := c.Query("method")
	path := c.Query("path")
	query := db.Order("created_at DESC")

	// Apply filters based on query parameters
	if startTime != "" {
		start, err := time.Parse(time.RFC3339, startTime)
		if err == nil {
			query = query.Where("created_at >= ?", start)
		}
	}
	if endTime != "" {
		end, err := time.Parse(time.RFC3339, endTime)
		if err == nil {
			query = query.Where("created_at <= ?", end)
		}
	}
	if method != "" {
		query = query.Where("method = ?", method)
	}
	if path != "" {
		query = query.Where("path LIKE ?", "%"+path+"%")
	}

	// Fetch logs from the database
	query.Find(&logs)

	// Pass logs to the template
	return c.Render("templates/logs/logs.html", fiber.Map{
		"Logs": logs,
	})
}

func VerifyAccess (session *session.Session, db *gorm.DB) *models.User {
	// Verify access
	access := session.Get("access")
	if access == nil {
		return nil
	}
	user, _ := DecodeAccessToken(access.(string), db)
	if (user == nil || !user.IsStaff) {
		session.Set("error", "Login first!")
		session.Save()
	}
	return user
}