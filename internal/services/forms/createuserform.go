package forms

import "strings"

type CreateUserForm struct {
	FormErrors      *FormErrors
	Name            string
	Email           string
	Username        string
	Password        string
	ConfirmPassword string
}

func NewCreateUserForm() *CreateUserForm {
	return &CreateUserForm{
		FormErrors: &FormErrors{},
	}
}

func (f *CreateUserForm) Validate() bool {
	f.FormErrors = &FormErrors{}

	// Username validation
	if strings.TrimSpace(f.Username) == "" {
		f.FormErrors.Add("username", "Username is required")
	} else if len(f.Username) < 3 {
		f.FormErrors.Add("username", "Username must be at least 3 characters")
	}

	// Name validation
	if strings.TrimSpace(f.Name) == "" {
		f.FormErrors.Add("name", "Name field is required")
	} else if len(f.Name) < 3 {
		f.FormErrors.Add("name", "Name field must be at least 3 characters")
	}

	// Email validation
	if strings.TrimSpace(f.Email) == "" {
		f.FormErrors.Add("email", "Email is required")
	} else if !strings.Contains(f.Email, "@") {
		f.FormErrors.Add("email", "Invalid email format")
	}

	// Password validation
	if strings.TrimSpace(f.Password) == "" {
		f.FormErrors.Add("password", "Password field is required")
	} else if len(f.Password) < 8 {
		f.FormErrors.Add("password", "Password must be at least 8 characters")
	}

	// Password confirmation validation
	if strings.TrimSpace(f.ConfirmPassword) == "" {
		f.FormErrors.Add("confirm_password", "Confirm Password field is required")
	} else if f.Password != f.ConfirmPassword {
		f.FormErrors.Add("confirm_password", "Passwords do not match")
	}

	return !f.FormErrors.HasErrors()
}
