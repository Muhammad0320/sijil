package identity

import (
	"context"
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
