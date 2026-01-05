package projects

import (
	"context"
	"errors"
	"fmt"
	"sijil-core/internals/auth"
	"sijil-core/internals/core/domain"
	"sijil-core/internals/utils"
)

var ErrForbidden = errors.New("you do not have permission to perform this action")
var ErrLimitReached = errors.New("plan limit reached")

type EmailSender func(email, subject, body string) error

type Service struct {
	repo   Repository
	mailer EmailSender
}

func NewService(repo Repository, mailer EmailSender) *Service {
	return &Service{
		repo:   repo,
		mailer: mailer,
	}
}

type CreateProjectResponse struct {
	ID        int    `json:"id"`
	APIKey    string `json:"api_key"`
	APISecret string `json:"api_secret"`
}

func (s *Service) CreateProject(ctx context.Context, userID int, req CreateProjectRequest, plan *domain.Plan) (*CreateProjectResponse, error) {
	// 1. Check Plan Limits
	if plan.MaxProjects != -1 {
		count, _ := s.repo.CountProjects(ctx, userID)

		if count >= plan.MaxProjects {
			return nil, ErrLimitReached
		}

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

func (s *Service) AddMember(ctx context.Context, userID int, projectID int, req AddMemberRequest, plan *domain.Plan) error {

	project, err := s.repo.GetByID(ctx, projectID)
	if err != nil {
		return err
	}

	isOwner := project.UserID == userID

	isAuth := isOwner
	if !isAuth {
		role, err := s.repo.GetRole(ctx, projectID, userID)
		if err == nil && role == "admin" {
			isAuth = true
		}
	}

	if !isAuth {
		return ErrForbidden
	}

	if plan.MaxMembers != -1 {
		currentMembers, _ := s.repo.CountMembers(ctx, projectID)
		if currentMembers >= plan.MaxMembers {
			return ErrLimitReached
		}
	}

	go func() {
		subject := fmt.Sprintf("Invitation to join %s on Sijil", project.Name)

		// Simple HTML template (We can upgrade this to React Email later)
		baseURL := utils.GetAppURL()
		body := fmt.Sprintf(`
            <div style="font-family: sans-serif; padding: 20px;">
                <h2>You've been invited!</h2>
                <p>You have been added to the project <strong>%s</strong> as a <strong>%s</strong>.</p>
                <p>Click below to access the dashboard:</p>
                <a href="%s/dashboard" style="background: #000; color: #fff; padding: 10px 20px; text-decoration: none; border-radius: 5px;">Open Dashboard</a>
            </div>
        `, project.Name, req.Role, baseURL)

		err = s.mailer(req.Email, subject, body)
		if err != nil {
			fmt.Println("SERVICE: failed to send email")
		}

	}()

	return s.repo.AddMember(ctx, projectID, req.Email, req.Role)
}

func (s *Service) GetMembers(ctx context.Context, userID, projectID int) ([]ProjectMember, error) {
	// 1. Fetch Project
	project, err := s.repo.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// 2. Check Permissions (Owner OR Admin ONLY)
	isOwner := project.UserID == userID
	isAdmin := false

	if !isOwner {
		role, err := s.repo.GetRole(ctx, projectID, userID)
		if err == nil && role == "admin" {
			isAdmin = true
		}
	}

	if !isOwner && !isAdmin {
		return nil, ErrForbidden // Viewers get blocked here.
	}

	return s.repo.ListMembers(ctx, projectID)
}
