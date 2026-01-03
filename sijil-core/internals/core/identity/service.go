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

func (s *Service) Register(ctx context.Context, req RegisterRequest) (string, *User, error) {

	hash, _ := auth.HashPasswod(req.Password)

	rawToken, _ := utils.GenerateRandomString(32)
	hashedToken := utils.Hashtoken(rawToken)
	expires := time.Now().Add(24 * time.Hour)

	u := &User{

		VerificationToken:        &hashedToken,
		VerificationTokenExpires: &expires,
		PlanID:                   1,
		FirstName:                req.FirstName,
		LastName:                 req.LastName,
		Email:                    req.Email,
		PasswordHash:             hash,
	}

	user, err := s.repo.Create(ctx, u)
	if err != nil {
		return "", &User{}, err
	}

	go func() {

		link := fmt.Sprintf("https://sijil.io/verify?token=%s", rawToken)
		html := fmt.Sprintf(`
			<h2>Welcome to Sijil, %s!</h2>
			<p>Please verify your email to activate your account.</p>
			<a href="%s" style="background:#000;color:#fff;padding:10px 20px;text-decoration:none;border-radius:4px;">Verify Email</a>
			<p>Or paste this link: %s</p>
		`, u.FirstName, link, link)

		err = s.mailer(u.Email, "Verify your Sijil Account", html)
		if err != nil {
			fmt.Println("SERVICE: failed tot send account verification email")

		}

	}()

	token, err := auth.CreateJWT(s.jwtSecret, user.ID)

	return token, user, err
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (string, *User, error) {

	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return "", &User{}, errors.New("invalid credentials")
	}

	if !auth.ComparePasswordHash(req.Password, user.PasswordHash) {
		return "", &User{}, errors.New("invalid credentials")
	}

	token, err := auth.CreateJWT(s.jwtSecret, user.ID)

	return token, user, err
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

	go func() {

		link := fmt.Sprintf("https://sijil.dev/reset-password?token=%s", rawToken)
		html := fmt.Sprintf(`
			<h2>Reset your Password</h2>
			<p>Someone requested a password reset for your Sijil account.</p>
			<a href="%s" style="background:#000;color:#fff;padding:10px 20px;text-decoration:none;border-radius:4px;">Reset Password</a>
			<p>If this wasn't you, ignore this email.</p>
		`, link)

		err = s.mailer(email, "Reset Password request", html)
	}()

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

func (s *Service) UpgradePlan(ctx context.Context, userID int, planName string) error {
	plan, err := s.repo.GetPlanByName(ctx, planName)
	if err != nil {
		return err
	}
	return s.repo.UpdateUserPlan(ctx, userID, plan.ID)
}
