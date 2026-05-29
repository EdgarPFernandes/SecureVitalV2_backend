package users

type RegisterUserPayload struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
	Role        string `json:"role"`
	Photo       string `json:"photo"`
}

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
