package observability

import (
	"context"
	"errors"
	"sijil-core/internals/core/projects"
	"sijil-core/internals/database"
	"sijil-core/internals/ingest"
	"time"
)

var ErrForbidden = errors.New("access denied to project")

type Service struct {
	repo         Repository
	projectsRepo projects.Repository
	engine       *ingest.IngestionEngine
}

func NewService(repo Repository, projectRepo projects.Repository, engine *ingest.IngestionEngine) *Service {

	return &Service{
		repo:         repo,
		projectsRepo: projectRepo,
		engine:       engine,
	}

}

// Ingest: The write path
func (s *Service) ProcessAndQueue(ctx context.Context, projectID int, log *LogEntry) {

	log.ProjectID = projectID
	if log.Timestamp.IsZero() {
		log.Timestamp = time.Now()
	}

	if len(log.Message) > 1000 {
		log.Message = log.Message[:10000] + "..."
	}

	s.engine.LogQueue <- database.LogEntry(*log)
	ingest.RecordQueued(1)
}

// Search: The read path
func (s *Service) Search(ctx context.Context, userID, projectID int, query string, limit, offset, retentionDays int) ([]LogEntry, error) {
	// 1. Permission Check
	if err := s.checkAccess(ctx, userID, projectID); err != nil {
		return nil, err
	}

	// 2. Defaults
	to := time.Now()
	from := to.Add(-24 * time.Hour)

	return s.repo.SearchLogs(ctx, projectID, limit, offset, query, from, to, retentionDays)
}

func (s *Service) checkAccess(ctx context.Context, userID, projectID int) error {
	role, err := s.projectsRepo.GetRole(ctx, projectID, userID)
	if err != nil {
		return err
	}
	if role == "" {
		return ErrForbidden
	}
	return nil
}
