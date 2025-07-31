package infra

import (
	"database/sql"

	"github.com/cristianorosa/eSimulate/backend/domain"
)

// UserExamRepositoryPG implementa UserExamRepository para PostgreSQL
type UserExamRepositoryPG struct {
	DB *sql.DB
}

// Create cria uma nova aplicação de exame
func (r *UserExamRepositoryPG) Create(userExam *domain.UserExam) error {
	query := `
		INSERT INTO user_exams (user_id, exam_id, started_at)
		VALUES ($1, $2, $3)
		RETURNING id`

	err := r.DB.QueryRow(
		query,
		userExam.UserID,
		userExam.ExamID,
		userExam.StartedAt,
	).Scan(&userExam.ID)

	if err != nil {
		return err
	}

	return nil
}

// Update atualiza uma aplicação de exame
func (r *UserExamRepositoryPG) Update(userExam *domain.UserExam) error {
	query := `
		UPDATE user_exams 
		SET finished_at = $1, total_score = $2, passed = $3, time_spent_minutes = $4
		WHERE id = $5`

	result, err := r.DB.Exec(
		query,
		userExam.FinishedAt,
		userExam.TotalScore,
		userExam.Passed,
		userExam.TimeSpentMinutes,
		userExam.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// FindByID busca uma aplicação de exame por ID
func (r *UserExamRepositoryPG) FindByID(id int) (*domain.UserExam, error) {
	query := `
		SELECT id, user_id, exam_id, started_at, finished_at, total_score, passed, time_spent_minutes
		FROM user_exams
		WHERE id = $1`

	userExam := &domain.UserExam{}
	var finishedAt sql.NullTime
	var totalScore sql.NullFloat64
	var passed sql.NullBool
	var timeSpent sql.NullInt32

	err := r.DB.QueryRow(query, id).Scan(
		&userExam.ID,
		&userExam.UserID,
		&userExam.ExamID,
		&userExam.StartedAt,
		&finishedAt,
		&totalScore,
		&passed,
		&timeSpent,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	if finishedAt.Valid {
		userExam.FinishedAt = &finishedAt.Time
	}
	if totalScore.Valid {
		userExam.TotalScore = &totalScore.Float64
	}
	if passed.Valid {
		userExam.Passed = &passed.Bool
	}
	if timeSpent.Valid {
		timeSpentInt := int(timeSpent.Int32)
		userExam.TimeSpentMinutes = &timeSpentInt
	}

	return userExam, nil
}

// ListByUser lista todas as aplicações de exame de um usuário
func (r *UserExamRepositoryPG) ListByUser(userID int) ([]*domain.UserExam, error) {
	query := `
		SELECT id, user_id, exam_id, started_at, finished_at, total_score, passed, time_spent_minutes
		FROM user_exams
		WHERE user_id = $1
		ORDER BY started_at DESC`

	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userExams []*domain.UserExam
	for rows.Next() {
		userExam := &domain.UserExam{}
		var finishedAt sql.NullTime
		var totalScore sql.NullFloat64
		var passed sql.NullBool
		var timeSpent sql.NullInt32

		err := rows.Scan(
			&userExam.ID,
			&userExam.UserID,
			&userExam.ExamID,
			&userExam.StartedAt,
			&finishedAt,
			&totalScore,
			&passed,
			&timeSpent,
		)
		if err != nil {
			return nil, err
		}

		if finishedAt.Valid {
			userExam.FinishedAt = &finishedAt.Time
		}
		if totalScore.Valid {
			userExam.TotalScore = &totalScore.Float64
		}
		if passed.Valid {
			userExam.Passed = &passed.Bool
		}
		if timeSpent.Valid {
			timeSpentInt := int(timeSpent.Int32)
			userExam.TimeSpentMinutes = &timeSpentInt
		}

		userExams = append(userExams, userExam)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return userExams, nil
}

// ListByExam lista todas as aplicações de um exame
func (r *UserExamRepositoryPG) ListByExam(examID int) ([]*domain.UserExam, error) {
	query := `
		SELECT id, user_id, exam_id, started_at, finished_at, total_score, passed, time_spent_minutes
		FROM user_exams
		WHERE exam_id = $1
		ORDER BY started_at DESC`

	rows, err := r.DB.Query(query, examID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userExams []*domain.UserExam
	for rows.Next() {
		userExam := &domain.UserExam{}
		var finishedAt sql.NullTime
		var totalScore sql.NullFloat64
		var passed sql.NullBool
		var timeSpent sql.NullInt32

		err := rows.Scan(
			&userExam.ID,
			&userExam.UserID,
			&userExam.ExamID,
			&userExam.StartedAt,
			&finishedAt,
			&totalScore,
			&passed,
			&timeSpent,
		)
		if err != nil {
			return nil, err
		}

		if finishedAt.Valid {
			userExam.FinishedAt = &finishedAt.Time
		}
		if totalScore.Valid {
			userExam.TotalScore = &totalScore.Float64
		}
		if passed.Valid {
			userExam.Passed = &passed.Bool
		}
		if timeSpent.Valid {
			timeSpentInt := int(timeSpent.Int32)
			userExam.TimeSpentMinutes = &timeSpentInt
		}

		userExams = append(userExams, userExam)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return userExams, nil
}

// FindActiveByUser busca um exame ativo (não finalizado) de um usuário
func (r *UserExamRepositoryPG) FindActiveByUser(userID int, examID int) (*domain.UserExam, error) {
	query := `
		SELECT id, user_id, exam_id, started_at, finished_at, total_score, passed, time_spent_minutes
		FROM user_exams
		WHERE user_id = $1 AND exam_id = $2 AND finished_at IS NULL
		ORDER BY started_at DESC
		LIMIT 1`

	userExam := &domain.UserExam{}
	var finishedAt sql.NullTime
	var totalScore sql.NullFloat64
	var passed sql.NullBool
	var timeSpent sql.NullInt32

	err := r.DB.QueryRow(query, userID, examID).Scan(
		&userExam.ID,
		&userExam.UserID,
		&userExam.ExamID,
		&userExam.StartedAt,
		&finishedAt,
		&totalScore,
		&passed,
		&timeSpent,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Não há exame ativo
		}
		return nil, err
	}

	if finishedAt.Valid {
		userExam.FinishedAt = &finishedAt.Time
	}
	if totalScore.Valid {
		userExam.TotalScore = &totalScore.Float64
	}
	if passed.Valid {
		userExam.Passed = &passed.Bool
	}
	if timeSpent.Valid {
		timeSpentInt := int(timeSpent.Int32)
		userExam.TimeSpentMinutes = &timeSpentInt
	}

	return userExam, nil
}
