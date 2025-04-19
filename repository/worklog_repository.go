// repository/worklog_repository.go
package repository

import (
	"time"

	"github.com/jinwoole/worklog-backend/models"
	"github.com/jmoiron/sqlx"
)

type WorkLogRepository interface {
	GetByUserAndDate(userID int, date time.Time) (*models.WorkLog, error)
	Create(log *models.WorkLog) error
	Update(log *models.WorkLog) error
	GetAllByUser(userID int) ([]models.WorkLog, error)
}

type workLogRepository struct {
	db *sqlx.DB
}

// NewWorkLogRepository returns a WorkLogRepository backed by sqlx.DB
func NewWorkLogRepository(db *sqlx.DB) WorkLogRepository {
	return &workLogRepository{db: db}
}

func (r *workLogRepository) GetByUserAndDate(userID int, date time.Time) (*models.WorkLog, error) {
	var log models.WorkLog
	err := r.db.Get(
		&log,
		`SELECT id, user_id, content, created_at FROM work_logs WHERE user_id = $1 AND created_at::date = $2`,
		userID,
		date.Format("2006-01-02"),
	)
	if err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *workLogRepository) Create(log *models.WorkLog) error {
	// Insert new work log and return generated ID
	return r.db.QueryRowx(
		`INSERT INTO work_logs (user_id, content, created_at) VALUES ($1, $2, $3) RETURNING id`,
		log.UserID,
		log.Content,
		log.CreatedAt,
	).Scan(&log.ID)
}

func (r *workLogRepository) Update(log *models.WorkLog) error {
	_, err := r.db.Exec(
		`UPDATE work_logs SET content = $1 WHERE user_id = $2 AND created_at::date = $3`,
		log.Content,
		log.UserID,
		log.CreatedAt.Format("2006-01-02"),
	)
	return err
}

func (r *workLogRepository) GetAllByUser(userID int) ([]models.WorkLog, error) {
	var logs []models.WorkLog
	err := r.db.Select(
		&logs,
		`SELECT id, user_id, content, created_at FROM work_logs WHERE user_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	return logs, nil
}
