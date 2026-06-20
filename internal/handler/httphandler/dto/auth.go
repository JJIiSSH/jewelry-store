package dto

import "github.com/JJIiSSH/jewelry-store/internal/domain"

// RegisterRequest is the public registration payload. There is deliberately no
// Role field — the role is assigned server-side; letting a client choose it
// would be a privilege-escalation hole.
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
	Name     string `json:"name" binding:"omitempty,max=255"`
	Phone    string `json:"phone" binding:"omitempty,max=50"`
}

// RegisterRequestToUser maps the request to a domain.User. It deliberately does
// not set ID, Role, PasswordHash or timestamps — those belong to the service.
func RegisterRequestToUser(r RegisterRequest) domain.User {

	var user domain.User

	user.Email = r.Email
	user.Name = r.Name
	user.Phone = r.Phone

	return user
}
