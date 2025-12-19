package identity

import (
	"context"
	"errors"
	"sijil-core/internals/auth"
)

type Service struct {
	repo      Repository
	jwtSecret string
}

func NewService(repo Repository, jwtSecret string) *Service {
	return &Service{repo: repo, jwtSecret: jwtSecret}
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) (string, error) {

	hash, _ := auth.HashPasswod(req.Password)

	u := &User{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		PasswordHash: hash,
		Plan:         "free",
	}

	id, err := s.repo.Create(ctx, u)
	if err != nil {
		return "", err
	}

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
