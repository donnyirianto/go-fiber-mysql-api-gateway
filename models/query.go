package models

type QueryRequest struct {
	Query string `validate:"required" json:"query"`
}
