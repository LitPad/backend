package routes

import (
	"fmt"

	"github.com/LitPad/backend/managers"
	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/choices"
	"github.com/LitPad/backend/models/scopes"
	"github.com/LitPad/backend/routes/helpers"
	"github.com/LitPad/backend/schemas"
	"github.com/LitPad/backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var notificationManager = managers.NotificationManager{}

// @Summary View User Profile
// @Description This endpoint views a user profile
// @Tags Profiles
// @Param username path string true "Username of user"
// @Success 200 {object} schemas.UserProfileResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Router /profiles/profile/{username} [get]
func (ep Endpoint) GetProfile(c *fiber.Ctx) error {
	db := ep.DB

	username := c.Params("username")
	if username == "" {
		return c.Status(400).JSON(utils.RequestErr(utils.ERR_INVALID_REQUEST, "Invalid path params"))
	}

	user := userManager.GetByUsername(db, username)
	if user == nil {
		return c.Status(404).JSON(utils.RequestErr(utils.ERR_NON_EXISTENT, "User does not exist!"))
	}

	response := schemas.UserProfileResponseSchema{
		ResponseSchema: ResponseMessage("Profile fetched successfully"),
		Data:           schemas.UserProfile{}.Init(*user),
	}
	return c.Status(200).JSON(response)
}

// @Summary Update User Profile
// @Description This endpoint updates a user's profile
// @Tags Profiles
// @Param profile formData schemas.UpdateUserProfileSchema true "Profile object"
// @Param avatar formData file false "Avatar Image to upload"
// @Success 200 {object} schemas.UserProfileResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Router /profiles/update [patch]
// @Security BearerAuth
func (ep Endpoint) UpdateProfile(c *fiber.Ctx) error {
	db := ep.DB
	user := RequestUser(c)
	data := schemas.UpdateUserProfileSchema{}
	if errCode, errData := ValidateFormRequest(c, &data); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}

	username := *data.Username
	if data.Username != nil {
		existingUser := models.User{Username: username}
		db.Not(models.User{BaseModel: models.BaseModel{ID: user.ID}}).Take(&existingUser, existingUser)
		if existingUser.ID != uuid.Nil {
			return c.Status(400).JSON(utils.ValidationErr("username", "Username is already taken"))
		}
		user.Username = username
	}
	if data.Name != nil {
		user.Name = data.Name
	}
	if data.Bio != nil {
		user.Bio = data.Bio
	}

	// Check and validate image
	file, err := ValidateImage(c, "avatar", false)
	if err != nil {
		return c.Status(422).JSON(err)
	}
	// Upload File
	if file != nil {
		avatar := UploadFile(file, string(choices.IF_AVATAR))
		user.Avatar = avatar
	}

	db.Save(&user)

	response := schemas.UserProfileResponseSchema{
		ResponseSchema: ResponseMessage("User details updated successfully"),
		Data:           schemas.UserProfile{}.Init(*user),
	}
	return c.Status(200).JSON(response)
}

// @Summary Update User Password
// @Description This endpoint updates a user's password
// @Tags Profiles
// @Param profile body schemas.UpdatePasswordSchema true "Password object"
// @Success 200 {object} schemas.ResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Router /profiles/update-password [put]
// @Security BearerAuth
func (ep Endpoint) UpdatePassword(c *fiber.Ctx) error {
	db := ep.DB

	data := schemas.UpdatePasswordSchema{}

	if errCode, errData := ValidateRequest(c, &data); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}

	user := RequestUser(c)

	if !utils.CheckPasswordHash(data.OldPassword, user.Password) {
		data := map[string]string{"old_password": "Password Mismatch"}

		return c.Status(400).JSON(utils.RequestErr(utils.ERR_PASSWORD_MISMATCH, "Invalid Entry", data))
	}

	if data.NewPassword == data.OldPassword {
		data := map[string]string{"new_password": "New password is same as old password"}
		return c.Status(400).JSON(utils.RequestErr(utils.ERR_PASSWORD_SAME, "Invalid Entry", data))
	}

	user.Password = utils.HashPassword(data.NewPassword)
	// Clear tokens to logout user
	user.Access = nil
	user.Refresh = nil
	db.Save(&user)
	return c.Status(200).JSON(ResponseMessage("Password updated successfully"))
}

// @Summary Toggle Follow Status
// @Description `This endpoint allows a user to follow or unfollow a writer`.
// @Tags Profiles
// @Param username path string true "Username of the user to follow or unfollow"
// @Success 200 {object} schemas.ResponseSchema
// @Failure 400 {object} utils.ErrorResponse "Returns an error for invalid request parameters."
// @Failure 403 {object} utils.ErrorResponse "Returns an error when trying to follow a user that isn't a writer"
// @Failure 404 {object} utils.ErrorResponse "Returns an error when either the user to follow or the follower user does not exist."
// @Failure 500 {object} utils.ErrorResponse "Returns an error when there is an internal server error or a transaction fails."
// @Router /profiles/profile/{username}/follow [get]
// @Security BearerAuth
func (ep Endpoint) FollowUser(c *fiber.Ctx) error {
	db := ep.DB
	user := RequestUser(c)
	toFollowUsername := c.Params("username")
	if toFollowUsername == "" {
		return c.Status(400).JSON(utils.RequestErr(utils.ERR_INVALID_REQUEST, "Invalid path parameter for username"))
	}

	followerUsername := user.Username

	if toFollowUsername == followerUsername {
		return c.Status(400).JSON(utils.RequestErr(utils.ERR_INVALID_REQUEST, "Cannot follow yourself"))
	}

	toFollowUser := models.User{}

	// Retrieve the user to follow
	if err := db.Scopes(scopes.VerifiedUserScope).Where(models.User{Username: toFollowUsername}).Take(&toFollowUser).Error; err != nil {
		return helpers.UserNotFoundError(c, "User to follow does not exist", err)
	}

	// check if both are readers
	if toFollowUser.AccountType == choices.ACCTYPE_READER {
		return c.Status(403).JSON(utils.RequestErr(utils.ERR_INVALID_REQUEST, "Readers cannot be followed"))
	}
	tx := db.Begin()

	// Toggle follow
	count := tx.Model(&user).Where("id = ?", toFollowUser.ID).Association("Followings").Count()
	alreadyFollowing := count > 0

	if alreadyFollowing {
		// Remove following and followers
		if err := tx.Model(&user).Association("Followings").Delete(&toFollowUser); err != nil {
			tx.Rollback()
			return c.Status(500).JSON(utils.RequestErr(utils.ERR_SERVER_ERROR, "Failed to unfollow user"))
		}
	} else {
		// Add the following relationship
		if err := tx.Model(&user).Omit("Followings.*").Association("Followings").Append(&toFollowUser); err != nil {
			tx.Rollback()
			return c.Status(500).JSON(utils.RequestErr(utils.ERR_SERVER_ERROR, "Failed to follow user"))
		}
	}

	if err := tx.Commit().Error; err != nil {
		return c.Status(500).JSON(utils.RequestErr(utils.ERR_SERVER_ERROR, "Failed to commit changes"))
	}

	message := "User followed successfully"
	if alreadyFollowing {
		message = "User unfollowed successfully"
	} else {
		// Create notification and send in socket
		notification := notificationManager.Create(
			db, user, toFollowUser, choices.NT_FOLLOWING,
			fmt.Sprintf("%s started following you.", user.Username),
			nil, nil, nil, nil,
		)
		SendNotificationInSocket(c, notification)
	}
	return c.Status(200).JSON(ResponseMessage(message))
}

// @Summary View Notifications
// @Description This endpoint allows a user to view his/her notificatios
// @Tags Profiles
// @Success 200 {object} schemas.NotificationsResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Router /profiles/notifications [get]
// @Security BearerAuth
func (ep Endpoint) GetNotifications(c *fiber.Ctx) error {
	db := ep.DB
	user := RequestUser(c)
	notifications := notificationManager.GetAllByUser(db, user)

	// Paginate and return notifications
	paginatedData, paginatedNotifications, err := PaginateQueryset(notifications, c, 50)
	if err != nil {
		return c.Status(400).JSON(err)
	}
	notifications = paginatedNotifications.([]models.Notification)
	response := schemas.NotificationsResponseSchema{
		ResponseSchema: ResponseMessage("Notifications fetched successfully"),
		Data: schemas.NotificationsResponseDataSchema{
			PaginatedResponseDataSchema: *paginatedData,
		}.Init(notifications),
	}
	return c.Status(200).JSON(response)
}

// @Summary Read Notification
// @Description This endpoint allows a user to read his/her notification.
// @Tags Profiles
// @Param notification body schemas.ReadNotificationSchema true "Notification Read object"
// @Success 200 {object} schemas.ResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Router /profiles/notifications/read [post]
// @Security BearerAuth
func (ep Endpoint) ReadNotification(c *fiber.Ctx) error {
	db := ep.DB

	data := schemas.ReadNotificationSchema{}
	if errCode, errData := ValidateRequest(c, &data); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}
	user := RequestUser(c)
	notificationID := data.ID
	markAllAsRead := data.MarkAllAsRead

	respMessage := "Notifications read"
	if markAllAsRead {
		// Mark all notifications as read
		notificationManager.MarkAsRead(db, user)
	} else if notificationID != nil {
		// Mark single notification as read
		err := notificationManager.ReadOne(db, user, *notificationID)
		if err != nil {
			return c.Status(404).JSON(err)
		}
		respMessage = "Notification read"
	}
	return c.Status(200).JSON(ResponseMessage(respMessage))
}
