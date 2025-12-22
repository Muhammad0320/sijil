package projects

import "time"

type Project struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`
	Name          string    `json:"name"`
	APIKey        string    `json:"api_key"`
	APISecretHash string    `json:"-"`
	CreatedAt     time.Time `json:"created_at"`
}

type ProjectMember struct {
	ProjectID int       `json:"project_id"`
	UserID    int       `json:"user_id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	JoinedAt  time.Time `json:"joined_at"`
}

// CreateProjectRequest capture the input
type CreateProjectRequest struct {
	Name string `json:"name" binding:"required,min=3,max=50"`
}

// AddMemberRequest captures the invitation
type AddMemberRequest struct {
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role" binding:"required,oneof=admin viewer"`
}
