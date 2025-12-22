package projects

import (
	"context"
	"errors"
	"sijil-core/internals/auth"
	"sijil-core/internals/database"
	"sijil-core/internals/utils"
)

var ErrForbidden = errors.New("you do not have permission to perform this action")
var ErrLimitReached = errors.New("plan limit reached")

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

type CreateProjectResponse struct {
	ID        int    `json:"id"`
	APIKey    string `json:"api_key"`
	APISecret string `json:"api_secret"`
}

func (s *Service) CreateProject(ctx context.Context, userID int, req CreateProjectRequest) (*CreateProjectResponse, error) {
	// 1. Check Plan Limits
	plan, err := s.repo.GetUserPlan(ctx, userID)
	if err != nil {
		return nil, err
	}

	count, _ := s.repo.CountProjects(ctx, userID)
	limits := database.GetPlanLimits(plan) // Assuming you have this helper in shared/database
	if count >= limits.MaxProject {
		return nil, ErrLimitReached
	}

	// 2. Generate Credentials
	randKey, _ := utils.GenerateRandomString(16)
	apiKey := "pk_live_" + randKey

	randSecret, _ := utils.GenerateRandomString(32)
	apiSecret := "sk_live_" + randSecret
	secretHash, _ := auth.HashPasswod(apiSecret)

	// 3. Save
	p := &Project{
		UserID:        userID,
		Name:          req.Name,
		APIKey:        apiKey,
		APISecretHash: secretHash,
	}

	id, err := s.repo.Create(ctx, p)
	if err != nil {
		return nil, err
	}

	return &CreateProjectResponse{
		ID:        id,
		APIKey:    apiKey,
		APISecret: apiSecret,
	}, nil
}

func (s *Service) ListProjects(ctx context.Context, userID int) ([]Project, error) {
	return s.repo.ListByUserID(ctx, userID)
}

func (s *Service) AddMember(ctx context.Context, userID int, projectID int, req AddMemberRequest) error {
	// 1. Check Permission (Must be Owner or Admin)
	role, err := s.repo.GetRole(ctx, projectID, userID)
	if err != nil {
		return err
	}
	if role != "owner" && role != "admin" {
		return ErrForbidden
	}

	project, err := s.repo.GetByID(ctx, projectID)
	if err != nil {
		return err
	}

	plan, _ := s.repo.GetUserPlan(ctx, project.UserID)
	count, _ := s.repo.CountMembers(ctx, projectID)
	limits := database.GetPlanLimits(plan)

	if count >= limits.MaxMemebers {
		return ErrLimitReached
	}

	// 3. Add
	return s.repo.AddMember(ctx, projectID, req.Email, req.Role)
}

func (s *Service) GetMembers(ctx context.Context, userID, projectID int) ([]ProjectMember, error) {

	role, err := s.repo.GetRole(ctx, projectID, userID)
	if err != nil {
		return nil, err
	}
	if role == "" {
		return nil, ErrForbidden
	}

	return s.repo.ListMembers(ctx, projectID)
}
