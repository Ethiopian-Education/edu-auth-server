package model

type timestamptz string

type Response struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type TwilioBody struct {
	To      string `json:"to"`
	Message string `json:"message"`
	From    string `json:"from,omitemty"`
}

type ConfirmPhone struct {
	PhoneNumber string `json:"phone_number"`
	Code        string `json:"code"`
}
type ConfirmEmail struct {
	Email string `json:"email"`
	Code        string `json:"code"`
}
