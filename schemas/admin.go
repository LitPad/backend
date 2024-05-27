package schemas


type UserProfilesResponseSchema struct {
	ResponseSchema
	Data []UserProfile `json:"data"`
}