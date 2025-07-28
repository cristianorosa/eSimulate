package infra

import (
	"database/sql"
	"github.com/cristianorosa/eSimulate/backend/domain"
)

// QuizRepositoryPG implementa QuizRepository usando PostgreSQL

type QuizRepositoryPG struct {
	DB *sql.DB
}

func (r *QuizRepositoryPG) Create(q *domain.Quiz) error {
	query := `INSERT INTO quizzes (title, description, theme_id, created_by) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.DB.QueryRow(query, q.Title, q.Description, q.ThemeID, q.CreatedBy).Scan(&q.ID)
	if err != nil {
		return err
	}
	for _, quest := range q.Questions {
		if err := r.addQuizQuestion(q.ID, quest.ID); err != nil {
			return err
		}
	}
	return nil
}

func (r *QuizRepositoryPG) addQuizQuestion(quizID, questionID int) error {
	query := `INSERT INTO quiz_questions (quiz_id, question_id) VALUES ($1, $2)`
	_, err := r.DB.Exec(query, quizID, questionID)
	return err
}

func (r *QuizRepositoryPG) Update(q *domain.Quiz) error {
	query := `UPDATE quizzes SET title = $1, description = $2, theme_id = $3 WHERE id = $4`
	_, err := r.DB.Exec(query, q.Title, q.Description, q.ThemeID, q.ID)
	return err
}

func (r *QuizRepositoryPG) Delete(id int) error {
	query := `DELETE FROM quizzes WHERE id = $1`
	_, err := r.DB.Exec(query, id)
	return err
}

func (r *QuizRepositoryPG) FindByID(id int) (*domain.Quiz, error) {
	query := `SELECT id, title, description, theme_id, created_by FROM quizzes WHERE id = $1`
	row := r.DB.QueryRow(query, id)
	var q domain.Quiz
	if err := row.Scan(&q.ID, &q.Title, &q.Description, &q.ThemeID, &q.CreatedBy); err != nil {
		return nil, err
	}
	q.Questions, _ = r.listQuizQuestions(q.ID)
	return &q, nil
}

func (r *QuizRepositoryPG) ListAll(themeID *int) ([]*domain.Quiz, error) {
	query := `SELECT id, title, description, theme_id, created_by FROM quizzes`
	var rows *sql.Rows
	var err error
	if themeID != nil {
		query += " WHERE theme_id = $1"
		rows, err = r.DB.Query(query, *themeID)
	} else {
		rows, err = r.DB.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var quizzes []*domain.Quiz
	for rows.Next() {
		var q domain.Quiz
		if err := rows.Scan(&q.ID, &q.Title, &q.Description, &q.ThemeID, &q.CreatedBy); err != nil {
			return nil, err
		}
		q.Questions, _ = r.listQuizQuestions(q.ID)
		quizzes = append(quizzes, &q)
	}
	return quizzes, nil
}

func (r *QuizRepositoryPG) listQuizQuestions(quizID int) ([]*domain.Question, error) {
	query := `SELECT q.id, q.theme_id, q.statement, q.explanation, q.created_by FROM quiz_questions qq JOIN questions q ON qq.question_id = q.id WHERE qq.quiz_id = $1`
	rows, err := r.DB.Query(query, quizID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var questions []*domain.Question
	for rows.Next() {
		var q domain.Question
		if err := rows.Scan(&q.ID, &q.ThemeID, &q.Statement, &q.Explanation, &q.CreatedBy); err != nil {
			return nil, err
		}
		questions = append(questions, &q)
	}
	return questions, nil
}
