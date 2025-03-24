package routes

import (
	"log"
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
        "success": sess.Get("success"),
    }
	// Remove the error from session after it's been rendered
    sess.Set("error", nil)
    sess.Set("success", nil)
    sess.Save()
	return c.Render("login", sessionData)
}

// Handles login form submission
func (ep Endpoint) HandleLogsLogin(c *fiber.Ctx) error {
	db := ep.DB
	session := Session(c, ep.Store)
	email := c.FormValue("email")
	password := c.FormValue("password")
	user := userManager.GetByEmail(db, email)
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
	
	accessToken := GenerateAccessToken(*user)
	user.Access = &accessToken
	// Hash password
	if user.Password == password {
		user.Password = utils.HashPassword(password)
	}
	db.Save(&user)
	session.Set("access", accessToken)
	session.Set("success", "Logged in successfully!")
	session.Set("error", nil)
	session.Save()
	return c.Redirect("/logs")
}

// Handles login form submission
func (ep Endpoint) HandleLogsLogout(c *fiber.Ctx) error {
	// Verify view access
	db := ep.DB
	session := Session(c, ep.Store)
	user := VerifyAccess(session, db)
	if user == nil {
		return c.Redirect("/logs/login")
	}
	user.Access = nil
	user.Refresh = nil
	db.Save(&user)
	session.Set("access", nil)
	session.Set("success", "Logged out successfully!")
	session.Save()
	return c.Redirect("/logs/login")
}

// Renders the logs page with optional filters
func (ep Endpoint) RenderLogs(c *fiber.Ctx) error {
	db := ep.DB
	// Verify page access
	session := Session(c, ep.Store)
	user := VerifyAccess(session, db)
	if user == nil {
		return c.Redirect("/logs/login")
	}
	logs := []models.Log{}


	// Get query parameters for filtering
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")
	method := c.Query("method")
	path := c.Query("path")
	status := c.QueryInt("status", 0)

	// Pagination parameters
	page := c.QueryInt("page", 1) // Default to page 1
	limit := c.QueryInt("limit", 100) // Default to 100 logs per page
	offset := (page - 1) * limit

	query := db.Model(&models.Log{}).Order("created_at DESC")
	
	// Apply filters based on query parameters
	if startTime != "" {
		// Check if startTime is missing seconds or timezone and add them
		if len(startTime) == 16 { // "2025-01-27T22:21" format
			startTime = startTime + ":00Z" // Add missing seconds and assume UTC timezone
		}
		
		start, err := time.Parse(time.RFC3339, startTime)
		if err == nil {
			query = query.Where("created_at >= ?", start)
		} else {
			log.Println("Error parsing startTime:", err)
		}
	}
	
	if endTime != "" {
		// Check if endTime is missing seconds or timezone and add them
		if len(endTime) == 16 { // "2025-01-27T22:21" format
			endTime = endTime + ":00Z" // Add missing seconds and assume UTC timezone
		}
	
		end, err := time.Parse(time.RFC3339, endTime)
		if err == nil {
			query = query.Where("created_at <= ?", end)
		} else {
			log.Println("Error parsing endTime:", err)
		}
	}
	
	if method != "" {
		query = query.Where("method = ?", method)
	}
	if path != "" {
		query = query.Where("path LIKE ?", "%"+path+"%")
	}
	if status != 0 {
		query = query.Where("status_code = ?", status)
	}

	// Fetch total count of logs for pagination
	var totalCount int64
	query.Count(&totalCount)

	// Fetch logs with pagination
	query.Limit(limit).Offset(offset).Find(&logs)

	totalPages := int((totalCount + int64(limit) - 1) / int64(limit))

	successMessage := session.Get("success")
	if successMessage != nil {
		session.Set("success", nil)
		session.Save()
	}
	return c.Render("logs", fiber.Map{
		"logs": logs,
		"success": successMessage,
		"page":       page,
		"limit":      limit,
		"totalPages": totalPages,
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