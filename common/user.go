package common

type UserCreationInput struct {
	FullName          string `json:"full_name"`
	// LastName          string `json:"last_name"`
	Email             string `json:"email"`
	// Phone             string `json:"phone"`
	// Password          string `json:"password"`
	// RegistrationToken string `json:"registration_token"`
}

func NewUserCreationInput() *UserCreationInput {
	return &UserCreationInput{}
}
