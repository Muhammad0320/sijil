package identity

import (
	"context"
	"errors"
	"fmt"
	"sijil-core/internals/auth"
	"sijil-core/internals/utils"
	"time"
)

type EmailSender func(email, subject, body string) error

type Service struct {
	repo      Repository
	jwtSecret string
	mailer    EmailSender
}

func NewService(repo Repository, jwtSecret string, mailer EmailSender) *Service {
	return &Service{repo: repo, jwtSecret: jwtSecret, mailer: mailer}
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) (string, error) {

	hash, _ := auth.HashPasswod(req.Password)

	rawToken, _ := utils.GenerateRandomString(32)
	hashedToken := utils.Hashtoken(rawToken)
	expires := time.Now().Add(24 * time.Hour)

	u := &User{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		PasswordHash: hash,

		Plan:                     "free",
		VerificationToken:        &hashedToken,
		VerificationTokenExpires: &expires,
	}

	id, err := s.repo.Create(ctx, u)
	if err != nil {
		return "", err
	}

	go s.mailer(u.Email, "Verify Account", rawToken)

	return auth.CreateJWT(s.jwtSecret, id)
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (string, error) {

	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if !auth.ComparePasswordHash(req.Password, user.PasswordHash) {
		return "", errors.New("invalid credentials")
	}

	return auth.CreateJWT(s.jwtSecret, user.ID)
}

func (s *Service) VerifyEmail(ctx context.Context, rawToken string) error {
	hashedToken := utils.Hashtoken(rawToken)
	return s.repo.VerifyUserAccount(ctx, hashedToken)
}

func (s *Service) ForgotPassword(ctx context.Context, email string) error {

	_, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil
	}

	rawToken, _ := utils.GenerateRandomString(32)
	hashToken := utils.Hashtoken(rawToken)
	expiry := time.Now().Add(15 * time.Minute)

	err = s.repo.SetPasswordResetToken(ctx, email, hashToken, expiry)
	if err != nil {
		return err
	}

	go func(email, token string) {
		fmt.Printf("📧 [Email Mock] To: %s | Subject: Reset Password | Link: https://sijil.dev/reset-password?token=%s\n", email, token)
	}(email, rawToken)

	return nil
}

func (s *Service) ResetPassword(ctx context.Context, rawToken, password string) error {

	hashedToken := utils.Hashtoken(rawToken)
	hashPassword, err := auth.HashPasswod(password)
	if err != nil {
		return err
	}

	return s.repo.ResetPasswordByToken(ctx, hashedToken, hashPassword)
}
