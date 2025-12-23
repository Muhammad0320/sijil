package observability

import (
	"context"
	"errors"
	"sijil-core/internals/database"
	"sijil-core/internals/ingest"
	"sijil-core/internals/projects"
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
func (s *Service) Ingest(ctx context.Context, projectID int, logs []LogEntry) error {
	// 1. Enrich & Validate
	for i := range logs {
		logs[i].ProjectID = projectID
		if logs[i].Timestamp.IsZero() {
			logs[i].Timestamp = time.Now()
		}
		// Clamp huge messages?
		if len(logs[i].Message) > 10000 {
			logs[i].Message = logs[i].Message[:10000] + "...(truncated)"
		}
	}

	// 2. Durability (WAL)

	var dbLogs []database.LogEntry
	for _, l := range logs {
		dbLogs = append(dbLogs, database.LogEntry(l)) // conversion
	}

	if err := s.engine.Wal.WriteBatch(dbLogs); err != nil {
		ingest.RecordError()
		ingest.RecordDropped(1)
		return errors.New("durability failure")
	}

	// 3. Processing (Queue)
	for _, l := range dbLogs {
		s.engine.LogQueue <- l
		ingest.RecordQueued(1)
	}

	return nil
}

// Search: The read path
func (s *Service) Search(ctx context.Context, userID, projectID int, query string, limit, offset int) ([]LogEntry, error) {
	// 1. Permission Check
	if err := s.checkAccess(ctx, userID, projectID); err != nil {
		return nil, err
	}

	// 2. Defaults
	to := time.Now()
	from := to.Add(-24 * time.Hour)

	// A HUGEEEEEEEEEEEEEEEEE   TRASH
	return s.repo.SearchLogs(ctx, projectID, limit, offset, query, from, to, 10)
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
