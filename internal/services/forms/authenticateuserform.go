package forms

import "strings"

type AuthenticateUserForm struct {
	FormErrors *FormErrors
	Email      string
	Password   string
}

func NewAuthenticateUserForm() *AuthenticateUserForm {
	return &AuthenticateUserForm{
		FormErrors: &FormErrors{},
	}
}

func (f *AuthenticateUserForm) Validate() bool {
	f.FormErrors = &FormErrors{}

	// Email validation
	if strings.TrimSpace(f.Email) == "" {
		f.FormErrors.Add("email", "Email is required")
	} else if !strings.Contains(f.Email, "@") {
		f.FormErrors.Add("email", "Invalid email format")
	}

	// Password validation
	if strings.TrimSpace(f.Password) == "" {
		f.FormErrors.Add("password", "Password is required")
	}

	return !f.FormErrors.HasErrors()
}
