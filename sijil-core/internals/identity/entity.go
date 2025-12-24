package identity

import "time"

type User struct {
	ID                       int        `json:"id"`
	FirstName                string     `json:"firstname"`
	LastName                 string     `json:"lastname"`
	Email                    string     `json:"email"`
	PasswordHash             string     `json:"-"`
	IsVerified               bool       `json:"is_verified"`
	VerificationToken        *string    `json:"-"`
	VerificationTokenExpires *time.Time `json:"-"`
	ResetToken               *string    `json:"-"`
	ResetTokenExpires        *time.Time `json:"-"`
	AvatarURL                string     `json:"avatar_url"`
	Plan                     string     `json:"plan"`
	CreatedAt                time.Time  `json:"created_at"`
}

type RegisterRequest struct {
	FirstName string `json:"firstname" binding:"required,min=2"`
	LastName  string `json:"lastname" binding:"required,min=2"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type Plan struct {
	Name          string
	MaxDailyLogs  int
	RetentionDays int
	MaxMemebers   int
	MaxProjects   int
}
