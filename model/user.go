package model

type User struct {
	ID                    string      `json:"id" graphql:"id"`
	Email                 *string      `json:"email" graphql:"email"`
	PhoneNumber           string      `json:"phone_number" graphql:"phone_number"`
	FirstName             string      `json:"first_name" graphql:"first_name"`
	MiddleName            string      `json:"middle_name" graphql:"middle_name"`
	LastName              string      `json:"last_name" graphql:"last_name"`
	Gender string `json:"gender" graphql:"gender"`
	Active                bool      `json:"active" graphql:"active"`
	Password              string      `json:"password" graphql:"password"`
	Enable2FA             bool        `json:"enable_2fa" graphql:"enable_2fa"`
	IsEmailConfirmed      bool        `json:"is_email_confirmed" graphql:"is_email_confirmed"`
	IsPhoneConfirmed      bool        `json:"is_phone_confirmed" graphql:"is_phone_confirmed"`
	PasswordChanged       bool        `json:"password_changed" graphql:"password_changed"`
	AlternatePhoneNumber  string        `json:"alternate_phone_number" graphql:"alternate_phone_number"`
	EmailConfirmedAt      timestamptz `json:"email_confirmed_at" graphql:"email_confirmed_at"`
	PhoneConfirmedAt      timestamptz `json:"phone_confirmed_at" graphql:"phone_confirmed_at"`
	LastPasswordChangedAt timestamptz `json:"last_password_changed_at" graphql:"last_password_changed_at"`
	CreatedAt             timestamptz `json:"created_at" graphql:"created_at"`
	UpdatedAt             timestamptz `json:"updated_at" graphql:"updated_at"`
	SignUpMethod             string `json:"signup_method" graphql:"signup_method"`
	UserRoles []struct {
		Role struct {
			Name string `json:"name" graphql:"name"`
		}`json:"role" graphql:"role"`
	}`json:"user_roles" graphql:"user_roles"`
}

type Signup struct {
	Email                 *string      `json:"email" graphql:"email,omitempty"`
	PhoneNumber           string      `json:"phone_number" graphql:"phone_number"`
	FirstName             string      `json:"first_name" graphql:"first_name"`
	MiddleName            string      `json:"middle_name" graphql:"middle_name"`
	LastName              string      `json:"last_name" graphql:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
	Password string `json:"password"`
	Gender string `json:"gender"`
}

type Login struct {
	Username string `json:"username" graphql:"username"`
	Password string `json:"password" graphql:"password"`
}