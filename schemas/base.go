package schemas

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
