package routes

import (
	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/routes/helpers"
	"github.com/LitPad/backend/schemas"
	"github.com/LitPad/backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// @Summary View User Profile
// @Description This endpoint views a user profile
// @Tags Profiles
// @Param username path string true "Username of user"
// @Success 200 {object} schemas.ResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Router /profiles/profile/{username} [get]
func (ep Endpoint) GetProfile(c *fiber.Ctx) error {
	db := ep.DB

	username := c.Params("username")
	if username == "" {
		return c.Status(400).JSON(utils.RequestErr(utils.ERR_INVALID_REQUEST, "Invalid path params"))
	}

	var user models.User
	err := db.Preload("Followers").Preload("Followings").Where("username = ? AND is_email_verified = ?", username, true).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(utils.RequestErr(utils.ERR_NON_EXISTENT, "User does not exist!"))
		}
		return c.Status(500).JSON(utils.RequestErr(utils.ERR_SERVER_ERROR, "Internal server error"))
	}

	response := schemas.UserProfileResponseSchema{
		ResponseSchema: ResponseMessage("Profile fetched successfully"),
		Data:           schemas.UserProfile{}.Init(user),
	}
	return c.Status(200).JSON(response)
}

// @Summary Update User Profile
// @Description This endpoint updates a user's profile
// @Tags Profiles
// @Param profile body schemas.UpdateUserProfileSchema true "Profile object"
// @Success 200 {object} schemas.ResponseSchema
// @Failure 400 {object} utils.ErrorResponse
// @Router /profiles/update [patch]
// @Security BearerAuth
func (ep Endpoint) UpdateProfile(c *fiber.Ctx) error {
	db := ep.DB
	savedUser := RequestUser(c)
	data := schemas.UpdateUserProfileSchema{}
	if errCode, errData := ValidateRequest(c, &data); errData != nil {
		return c.Status(*errCode).JSON(errData)
	}

	username := *data.Username
	if data.Username != nil {
		searchUser := models.User{Username: username}
		db.Not(models.User{BaseModel: models.BaseModel{ID: savedUser.ID}}).Take(&searchUser, searchUser)
		if searchUser.ID != uuid.Nil {
			data := map[string]string{"username": "Username is already taken"}
			return c.Status(400).JSON(utils.RequestErr(utils.ERR_INVALID_REQUEST, "Invalid Entry", data))
		}
		savedUser.Username = username
	}

	// current design supports update of username only hence the code looking like this
	db.Save(&savedUser)

	response := schemas.UserProfileResponseSchema{
		ResponseSchema: ResponseMessage("User details updated successfully"),
		Data:           schemas.UserProfile{}.Init(*savedUser),
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


func (ep Endpoint) FollowUser(c *fiber.Ctx) error {
    db := ep.DB

    toFollowUsername := c.Params("username")
    if toFollowUsername == "" {
        return c.Status(400).JSON(utils.RequestErr(utils.ERR_INVALID_REQUEST, "Invalid path parameter for username"))
    }

	followerUsername := RequestUser(c).Username

    if toFollowUsername == followerUsername {
        return c.Status(400).JSON(utils.RequestErr(utils.ERR_INVALID_REQUEST, "Cannot follow yourself"))
    }

    var toFollowUser, followerUser models.User

    // Retrieve the user to follow
    if err := db.Where("username = ? AND is_email_verified = ?", toFollowUsername, true).First(&toFollowUser).Error; err != nil {
		return helpers.UserNotFoundError(c, "User to follow does not exist", err)
	}
    

    // Retrieve the follower user
	if err := db.Where("username = ? AND is_email_verified = ?", followerUsername, true).First(&followerUser).Error; err !=nil{
		return helpers.UserNotFoundError(c, "Follower user does not exist", err)
	}

	// check if both are readers
	  if toFollowUser.AccountType != "READER" || followerUser.AccountType != "READER" {
        return c.Status(403).JSON(utils.RequestErr(utils.ERR_INVALID_REQUEST, "Only users of type 'Reader' can follow each other"))
    }

	 tx := db.Begin()

	  // Toggle follow
    var existingFollowings []models.User

    tx.Model(&followerUser).Association("Followings").Find(&existingFollowings, "id = ?", toFollowUser.ID)

    alreadyFollowing := len(existingFollowings) > 0

	if alreadyFollowing {
		 // Remove following and followers
        if err := tx.Model(&followerUser).Association("Followings").Delete(&toFollowUser); err != nil {
            tx.Rollback()
            return c.Status(500).JSON(utils.RequestErr(utils.ERR_SERVER_ERROR, "Failed to unfollow user"))
        }
        if err := tx.Model(&toFollowUser).Association("Followers").Delete(&followerUser); err != nil {
            tx.Rollback()
            return c.Status(500).JSON(utils.RequestErr(utils.ERR_SERVER_ERROR, "Failed to remove follower"))
        }
	} else {
		// Add the following relationship
		if err := tx.Model(&followerUser).Association("Followings").Append(&toFollowUser); err != nil {
			tx.Rollback()
			return c.Status(500).JSON(utils.RequestErr(utils.ERR_SERVER_ERROR, "Failed to follow user"))
		}

		// Add the follower relationship
		if err := tx.Model(&toFollowUser).Association("Followers").Append(&followerUser); err != nil {
			tx.Rollback()
			return c.Status(500).JSON(utils.RequestErr(utils.ERR_SERVER_ERROR, "Failed to add follower"))
		}
	}

	 if err := tx.Commit().Error; err != nil {
        return c.Status(500).JSON(utils.RequestErr(utils.ERR_SERVER_ERROR, "Failed to commit changes"))
    }

   if alreadyFollowing {
        return c.Status(200).JSON(fiber.Map{
            "status":  "success",
            "message": "Successfully unfollowed the user",
        })
    } else {
        return c.Status(200).JSON(fiber.Map{
            "status":  "success",
            "message": "Successfully followed the user",
        })
    }
}