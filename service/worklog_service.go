// service/worklog_service.go
package service

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jinwoole/worklog-backend/models"
	"github.com/jinwoole/worklog-backend/repository"
)

// WorkLogService defines business logic for work logs
type WorkLogService interface {
	CreateWorkLog(userID int, content string) (*models.WorkLog, error)
	UpdateWorkLog(userID int, content string) error
	GetAllWorkLogs(userID int) ([]models.WorkLog, error)
}

// workLogService is a concrete implementation of WorkLogService
type workLogService struct {
	workLogRepo repository.WorkLogRepository
}

// NewWorkLogService constructs a WorkLogService with given repository
func NewWorkLogService(workLogRepo repository.WorkLogRepository) WorkLogService {
	return &workLogService{workLogRepo: workLogRepo}
}

// CreateWorkLog adds a new log if none exists for today
func (s *workLogService) CreateWorkLog(userID int, content string) (*models.WorkLog, error) {
	today := time.Now()
	_, err := s.workLogRepo.GetByUserAndDate(userID, today)
	if err == nil {
		return nil, errors.New("today's work log already exists")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	logEntry := &models.WorkLog{
		UserID:    userID,
		Content:   content,
		CreatedAt: today,
	}
	if err := s.workLogRepo.Create(logEntry); err != nil {
		return nil, err
	}
	return logEntry, nil
}

// UpdateWorkLog modifies today's log content
func (s *workLogService) UpdateWorkLog(userID int, content string) error {
	today := time.Now()
	logEntry, err := s.workLogRepo.GetByUserAndDate(userID, today)
	if err != nil {
		return err
	}
	logEntry.Content = content
	return s.workLogRepo.Update(logEntry)
}

// GetAllWorkLogs retrieves all logs for a user
func (s *workLogService) GetAllWorkLogs(userID int) ([]models.WorkLog, error) {
	return s.workLogRepo.GetAllByUser(userID)
}
