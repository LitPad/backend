package models

type Log struct {
	BaseModel
	StatusCode  int
	Method      string
	Path        string
	IP          string
	PathParams  string
	QueryParams string
	Body        string
}
