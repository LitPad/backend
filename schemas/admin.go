package schemas


type UserProfilesResponseSchema struct {
	ResponseSchema
	Data []UserProfile `json:"data"`
}

type BookResponseSchema struct{
	ResponseSchema
	Data []BookSchema `json:data`
}