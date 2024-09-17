package schemas

import "github.com/LitPad/backend/models"

type UserProfilesResponseDataSchema struct {
	PaginatedResponseDataSchema
	Items []UserProfile `json:"users"`
}

func (u UserProfilesResponseDataSchema) Init(users []models.User) UserProfilesResponseDataSchema {
	// Set Initial Data
	userItems := make([]UserProfile, 0)

	for _, user := range users {
		userItems = append(userItems, UserProfile{}.Init(user))
	}
	u.Items = userItems
	return u
}


type UserProfilesResponseSchema struct {
	ResponseSchema
	Data UserProfilesResponseDataSchema `json:"data"`
}