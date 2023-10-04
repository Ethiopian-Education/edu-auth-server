package model

type OTP struct {
	ID      string `json:"id" graphql:"id"`
	Code    string `json:"code" graphql:"code"`
	Type    string `json:"type" graphql:"type"`
	UserID  string `json:"user_id" graphql:"user_id"`
	IsValid bool   `json:"is_valid" graphql:"is_valid"`
}
