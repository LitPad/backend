package tests

import (
	"fmt"
	"testing"

	"github.com/LitPad/backend/schemas"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func getProfile(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	t.Run("Reject Profile Fetch Due to Non Existent User", func(t *testing.T) {
		url := fmt.Sprintf("%s/profile/invalid-username", baseUrl)
		res := ProcessTestGetOrDelete(app, url, "GET")
		// Assert Status code
		assert.Equal(t, 404, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "User does not exist!", body["message"])
	})

	t.Run("Accept User Profile Fetch", func(t *testing.T) {
		user := TestVerifiedUser(db)
		url := fmt.Sprintf("%s/profile/%s", baseUrl, user.Username)
		res := ProcessTestGetOrDelete(app, url, "GET")
		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Profile fetched successfully", body["message"])
	})
}

func updateProfile(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	user := TestVerifiedUser(db)
	anotherUser := TestUser(db)
	token := AccessToken(db, user)
	name := "Test Name Updated"
	bio := "Test Bio Updated"
	username := anotherUser.Username
	profileData := schemas.UpdateUserProfileSchema{Name: &name, Bio: &bio, Username: &username}

	t.Run("Reject Profile Update Due To Already Used Username", func(t *testing.T) {
		url := fmt.Sprintf("%s/update", baseUrl)
		res := ProcessMultipartTestBody(t, app, url, "PATCH", profileData, []string{}, []string{}, token)
		assert.Equal(t, 422, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Invalid Entry", body["message"])
		assert.Equal(t, "Username is already taken", body["data"].(map[string]interface{})["username"])
	})

	t.Run("Accept Profile Update Due To Valid Data", func(t *testing.T) {
		profileData.Username = &user.Username
		url := fmt.Sprintf("%s/update", baseUrl)
		res := ProcessMultipartTestBody(t, app, url, "PATCH", profileData, []string{}, []string{}, token)
		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "User details updated successfully", body["message"])
	})
}

func updatePassword(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	user := TestVerifiedUser(db)
	token := AccessToken(db, user)
	password := "newpassword"
	passwordData := schemas.UpdatePasswordSchema{OldPassword: password, NewPassword: password}

	t.Run("Reject Password Update Due To New Password Same As Old", func(t *testing.T) {
		url := fmt.Sprintf("%s/update-password", baseUrl)
		res := ProcessJsonTestBody(t, app, url, "PUT", passwordData, token)
		assert.Equal(t, 422, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Invalid Entry", body["message"])
		assert.Equal(t, "New password is same as old password", body["data"].(map[string]interface{})["new_password"])
	})

	t.Run("Reject Password Update Due To Old Password Mismatch", func(t *testing.T) {
		passwordData.OldPassword = "invalidpassword"
		url := fmt.Sprintf("%s/update-password", baseUrl)
		res := ProcessJsonTestBody(t, app, url, "PUT", passwordData, token)
		assert.Equal(t, 422, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Invalid Entry", body["message"])
		assert.Equal(t, "Password Mismatch", body["data"].(map[string]interface{})["old_password"])
	})

	t.Run("Accept Password Update Due To Valid Data", func(t *testing.T) {
		passwordData.OldPassword = MASTER_PASSWORD
		url := fmt.Sprintf("%s/update-password", baseUrl)
		res := ProcessJsonTestBody(t, app, url, "PUT", passwordData, token)
		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Password updated successfully", body["message"])
	})
}

func followUser(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	user := TestVerifiedUser(db)
	user2 := TestVerifiedUser(db, true)
	author := TestAuthor(db)
	token := AccessToken(db, user)

	t.Run("Reject User Follow If User Tries To Follow Himself", func(t *testing.T) {
		url := fmt.Sprintf("%s/profile/%s/follow", baseUrl, user.Username)
		res := ProcessTestGetOrDelete(app, url, "GET", token)
		assert.Equal(t, 400, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Cannot follow yourself", body["message"])
	})

	t.Run("Reject Follow If User To Follow Does Not Exist", func(t *testing.T) {
		url := fmt.Sprintf("%s/profile/%s/follow", baseUrl, "invalid-username")
		res := ProcessTestGetOrDelete(app, url, "GET", token)
		assert.Equal(t, 404, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "User to follow does not exist", body["message"])
	})

	t.Run("Reject Follow If You're trying To Follow A Reader", func(t *testing.T) {
		url := fmt.Sprintf("%s/profile/%s/follow", baseUrl, user2.Username)
		res := ProcessTestGetOrDelete(app, url, "GET", token)
		assert.Equal(t, 403, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Readers cannot be followed", body["message"])
	})

	t.Run("Accept User Follow Due To Valid Conditions", func(t *testing.T) {
		url := fmt.Sprintf("%s/profile/%s/follow", baseUrl, author.Username)
		res := ProcessTestGetOrDelete(app, url, "GET", token)
		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "User followed successfully", body["message"])
	})
}

func getNotifications(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	user := TestVerifiedUser(db)
	author := TestAuthor(db)
	token := AccessToken(db, author)
	NotificationData(db, user, author)
	t.Run("Accept User Notifications Fetch", func(t *testing.T) {
		url := fmt.Sprintf("%s/notifications", baseUrl)
		res := ProcessTestGetOrDelete(app, url, "GET", token)
		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Notifications fetched successfully", body["message"])
	})
}

func readNotification(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	sender := TestVerifiedUser(db)
	receiver := TestAuthor(db)
	token := AccessToken(db, receiver)
	notification := NotificationData(db, sender, receiver)
	id := uuid.New()
	notificationReadData := schemas.ReadNotificationSchema{MarkAllAsRead: false, ID: &id}
	t.Run("Reject Notification Read Due To Invalid ID", func(t *testing.T) {
		url := fmt.Sprintf("%s/notifications/read", baseUrl)
		res := ProcessJsonTestBody(t, app, url, "POST", notificationReadData, token)
		assert.Equal(t, 404, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "User has no notification with that ID", body["message"])
	})

	t.Run("Accept Notification Read Due To Valid Data", func(t *testing.T) {
		notificationReadData.ID = &notification.ID
		url := fmt.Sprintf("%s/notifications/read", baseUrl)
		res := ProcessJsonTestBody(t, app, url, "POST", notificationReadData, token)
		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Notification read successfully", body["message"])
	})
}

func TestProfiles(t *testing.T) {
	app := fiber.New()
	db := Setup(t, app)
	baseUrl := "/api/v1/profiles"

	// Run Profiles Endpoint Tests
	getProfile(t, app, db, baseUrl)
	updateProfile(t, app, db, baseUrl)
	updatePassword(t, app, db, baseUrl)
	followUser(t, app, db, baseUrl)
	getNotifications(t, app, db, baseUrl)
	readNotification(t, app, db, baseUrl)
}