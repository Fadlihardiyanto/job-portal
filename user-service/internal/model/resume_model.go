package model

type ResponseResume struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Attachment string `json:"attachment"`
	UserID     string `json:"user_id"`
	Status     string `json:"status"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	Users      User   `json:"users"`
}

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	About     string `json:"about"`
	Photo     string `json:"photo"`
	Role      string `json:"role"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
type RequestResume struct {
	Name       string `json:"name" validate:"required"`
	Attachment string `json:"attachment" validate:"required"`
	UserID     string `json:"user_id" validate:"required"`
}

type UpdateResume struct {
	Name       string `json:"name" validate:"omitempty"`
	Attachment string `json:"attachment" validate:"omitempty"`
	UserID     int    `json:"user_id" validate:"omitempty"`
}

type RequestFindResumeByUser struct {
	UserID string `json:"user_id" validate:"required"`
}

type RequestFindResumeByID struct {
	ID string `json:"id" validate:"required"`
}
