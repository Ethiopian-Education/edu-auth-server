package model

type timestamptz string

type Response struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}
