package users

import repo "github.com/EdgarPFernandes/SecureVitalV2_backend/internal/adapters/postgresql/sqlc"

type UserResponse struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	PhoneNumber *string `json:"phone_number"`
	Role        string  `json:"role"`
	Photo       *string `json:"photo"`
}

func toUserResponse(user repo.User) UserResponse {

	var phone *string
	if user.PhoneNumber.Valid {
		phone = &user.PhoneNumber.String
	}

	var photo *string
	if user.Photo.Valid {
		photo = &user.Photo.String
	}

	return UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: phone,
		Role:        user.Role,
		Photo:       photo,
	}
}
