package infra

import (
	"database/sql"

	"github.com/cristianorosa/eSimulate/backend/domain"
)

// QuestionRepositoryPG implementa QuestionRepository usando PostgreSQL

// QuestionRepositoryPG implementa QuestionRepository usando PostgreSQL.
type QuestionRepositoryPG struct {
	DB *sql.DB
}

// Create cria uma nova questão no banco de dados.
func (r *QuestionRepositoryPG) Create(q *domain.Question) error {
	query := `INSERT INTO questions (domain_id, statement, explanation, created_by) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.DB.QueryRow(query, q.DomainID, q.Statement, q.Explanation, q.CreatedBy).Scan(&q.ID)
	if err != nil {
		return err
	}
	for _, opt := range q.Options {
		opt.QuestionID = q.ID
		if err := r.createOption(opt); err != nil {
			return err
		}
	}
	return nil
}

func (r *QuestionRepositoryPG) createOption(opt *domain.Option) error {
	query := `INSERT INTO options (question_id, text, is_correct, explanation) VALUES ($1, $2, $3, $4) RETURNING id`
	return r.DB.QueryRow(query, opt.QuestionID, opt.Text, opt.IsCorrect, opt.Explanation).Scan(&opt.ID)
}

// Update atualiza uma questão existente no banco de dados.
func (r *QuestionRepositoryPG) Update(q *domain.Question) error {
	query := `UPDATE questions SET domain_id = $1, statement = $2, explanation = $3 WHERE id = $4`
	_, err := r.DB.Exec(query, q.DomainID, q.Statement, q.Explanation, q.ID)
	return err
}

// Delete remove uma questão do banco de dados.
func (r *QuestionRepositoryPG) Delete(id int) error {
	query := `DELETE FROM questions WHERE id = $1`
	_, err := r.DB.Exec(query, id)
	return err
}

// FindByID busca uma questão pelo seu ID.
func (r *QuestionRepositoryPG) FindByID(id int) (*domain.Question, error) {
	query := `SELECT id, domain_id, statement, explanation, created_by FROM questions WHERE id = $1`
	row := r.DB.QueryRow(query, id)
	var q domain.Question
	if err := row.Scan(&q.ID, &q.DomainID, &q.Statement, &q.Explanation, &q.CreatedBy); err != nil {
		return nil, err
	}
	q.Options, _ = r.listOptions(q.ID)
	return &q, nil
}

// ListByExam lista todas as questões de um exame
func (r *QuestionRepositoryPG) ListByExam(examID int) ([]*domain.Question, error) {
	query := `SELECT id, domain_id, statement, explanation, created_by FROM questions WHERE exam_id = $1`
	rows, err := r.DB.Query(query, examID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var questions []*domain.Question
	for rows.Next() {
		var q domain.Question
		if err := rows.Scan(&q.ID, &q.DomainID, &q.Statement, &q.Explanation, &q.CreatedBy); err != nil {
			return nil, err
		}
		q.Options, _ = r.listOptions(q.ID)
		questions = append(questions, &q)
	}
	return questions, nil
}

// ListByDomain lista todas as questões de um domínio
func (r *QuestionRepositoryPG) ListByDomain(domainID int) ([]*domain.Question, error) {
	query := `SELECT id, domain_id, statement, explanation, created_by FROM questions WHERE domain_id = $1`
	rows, err := r.DB.Query(query, domainID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var questions []*domain.Question
	for rows.Next() {
		var q domain.Question
		if err := rows.Scan(&q.ID, &q.DomainID, &q.Statement, &q.Explanation, &q.CreatedBy); err != nil {
			return nil, err
		}
		q.Options, _ = r.listOptions(q.ID)
		questions = append(questions, &q)
	}
	return questions, nil
}

// ListActive lista todas as questões ativas
func (r *QuestionRepositoryPG) ListActive() ([]*domain.Question, error) {
	query := `SELECT id, domain_id, statement, explanation, created_by FROM questions WHERE is_active = true`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var questions []*domain.Question
	for rows.Next() {
		var q domain.Question
		if err := rows.Scan(&q.ID, &q.DomainID, &q.Statement, &q.Explanation, &q.CreatedBy); err != nil {
			return nil, err
		}
		q.Options, _ = r.listOptions(q.ID)
		questions = append(questions, &q)
	}
	return questions, nil
}

func (r *QuestionRepositoryPG) listOptions(questionID int) ([]*domain.Option, error) {
	query := `SELECT id, question_id, text, is_correct, explanation FROM options WHERE question_id = $1`
	rows, err := r.DB.Query(query, questionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var opts []*domain.Option
	for rows.Next() {
		var o domain.Option
		if err := rows.Scan(&o.ID, &o.QuestionID, &o.Text, &o.IsCorrect, &o.Explanation); err != nil {
			return nil, err
		}
		opts = append(opts, &o)
	}
	return opts, nil
}
