package model

type OTP struct {
	ID      string `json:"id" graphql:"id"`
	Code    string `json:"code" graphql:"code"`
	Type    string `json:"type" graphql:"type"`
	Expired bool   `json:"expired" graphql:"expired"`
	UserID  string `json:"user_id" graphql:"user_id"`
}
