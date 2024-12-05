package schemas

import "github.com/LitPad/backend/models"

type ResponseSchema struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Data fetched/created/updated/deleted"`
}

func (obj ResponseSchema) Init() ResponseSchema {
	if obj.Status == "" {
		obj.Status = "success"
	}
	return obj
}

type PaginatedResponseDataSchema struct {
	PerPage     uint `json:"per_page" example:"100"`
	CurrentPage uint `json:"current_page" example:"1"`
	LastPage    uint `json:"last_page" example:"100"`
}

type UserDataSchema struct {
	// For short user data
	FullName string `json:"full_name"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

func (u UserDataSchema) Init(user models.User) UserDataSchema {
	u.FullName = user.FullName()
	u.Username = user.Username
	u.Avatar = user.Avatar
	return u
}