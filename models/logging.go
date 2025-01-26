package models

type Log struct {
	BaseModel
	StatusCode int
	Method     string
	Path       string
	IP         string
	Params     *string
	ReqID      string
	Body       []byte   // Raw request body
}
