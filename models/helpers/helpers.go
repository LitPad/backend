package helpers

// import (
// 	"github.com/LitPad/backend/models"
// 	"gorm.io/gorm"
// )

// func GetFollowingsCount (db *gorm.DB, user models.User) {
// 	count := db.Model("user_followers").Where("follower_id = ?", user.ID).Association("Followings").Count()
// }

// func GetFollowings (db *gorm.DB, user models.User) {
// 	count := db.Model("user_followers").Where("follower_id = ?", user.ID).Association("Followings").Count()
// }