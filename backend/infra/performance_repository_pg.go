package infra

import (
	"database/sql"
	"log"

	"github.com/cristianorosa/eSimulate/backend/domain"
)

// PerformanceRepositoryPG implementa PerformanceRepository usando PostgreSQL
type PerformanceRepositoryPG struct {
	DB *sql.DB
}

// Create cria um novo registro de performance
func (r *PerformanceRepositoryPG) Create(performance *domain.Performance) error {
	query := `
		INSERT INTO performance (user_exam_id, domain_id, questions_answered, correct_answers, score_percentage, needs_improvement, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP)
		RETURNING id, created_at`

	err := r.DB.QueryRow(
		query,
		performance.UserExamID,
		performance.DomainID,
		performance.QuestionsAnswered,
		performance.CorrectAnswers,
		performance.ScorePercentage,
		performance.NeedsImprovement,
	).Scan(&performance.ID, &performance.CreatedAt)

	if err != nil {
		return err
	}
	return nil
}

// Update atualiza um registro de performance
func (r *PerformanceRepositoryPG) Update(performance *domain.Performance) error {
	query := `
		UPDATE performance 
		SET user_exam_id = $1, domain_id = $2, questions_answered = $3, correct_answers = $4, score_percentage = $5, needs_improvement = $6
		WHERE id = $7`

	result, err := r.DB.Exec(
		query,
		performance.UserExamID,
		performance.DomainID,
		performance.QuestionsAnswered,
		performance.CorrectAnswers,
		performance.ScorePercentage,
		performance.NeedsImprovement,
		performance.ID,
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

// ListByUserExam lista todas as performances de um user_exam
func (r *PerformanceRepositoryPG) ListByUserExam(userExamID int) ([]*domain.Performance, error) {
	query := `
		SELECT id, user_exam_id, domain_id, questions_answered, correct_answers, score_percentage, needs_improvement, created_at
		FROM performance 
		WHERE user_exam_id = $1
		ORDER BY created_at DESC`

	rows, err := r.DB.Query(query, userExamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var performances []*domain.Performance
	for rows.Next() {
		var performance domain.Performance
		err := rows.Scan(
			&performance.ID,
			&performance.UserExamID,
			&performance.DomainID,
			&performance.QuestionsAnswered,
			&performance.CorrectAnswers,
			&performance.ScorePercentage,
			&performance.NeedsImprovement,
			&performance.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		performances = append(performances, &performance)
	}

	return performances, nil
}

// GetReport obtém o relatório de desempenho de um usuário.
func (r *PerformanceRepositoryPG) GetReport(userID int) (*domain.PerformanceReport, error) {
	log.Printf("PerformanceRepositoryPG: Iniciando relatório para usuário %d", userID)

	var report domain.PerformanceReport
	report.UserID = userID

	// Verificar se as tabelas existem antes de fazer as consultas
	var tableExists bool

	// Verificar se a tabela user_quiz existe
	log.Printf("PerformanceRepositoryPG: Verificando existência da tabela user_quiz...")
	err := r.DB.QueryRow(`SELECT EXISTS (
		SELECT FROM information_schema.tables 
		WHERE table_schema = 'public' 
		AND table_name = 'user_quiz'
	)`).Scan(&tableExists)
	if err != nil {
		log.Printf("PerformanceRepositoryPG: Erro ao verificar tabela user_quiz: %v", err)
		return nil, err
	}

	log.Printf("PerformanceRepositoryPG: Tabela user_quiz existe: %v", tableExists)

	if !tableExists {
		// Se as tabelas não existem, retornar relatório vazio
		log.Printf("PerformanceRepositoryPG: Tabela user_quiz não existe, retornando relatório vazio")
		report.TotalQuizzes = 0
		report.TotalQuestions = 0
		report.CorrectAnswers = 0
		report.Accuracy = 0.0
		return &report, nil
	}

	// Total de simulados realizados (contando tanto user_quiz quanto user_exams)
	log.Printf("PerformanceRepositoryPG: Contando simulados para usuário %d", userID)

	// Verificar se a tabela user_exams existe
	var userExamsExists bool
	err = r.DB.QueryRow(`SELECT EXISTS (
		SELECT FROM information_schema.tables 
		WHERE table_schema = 'public' 
		AND table_name = 'user_exams'
	)`).Scan(&userExamsExists)
	if err != nil {
		log.Printf("PerformanceRepositoryPG: Erro ao verificar tabela user_exams: %v", err)
		return nil, err
	}

	var quizCount, examCount int
	// Contar user_quiz
	err = r.DB.QueryRow(`SELECT COUNT(*) FROM user_quiz WHERE user_id = $1`, userID).Scan(&quizCount)
	if err != nil {
		log.Printf("PerformanceRepositoryPG: Erro ao contar user_quiz: %v", err)
		quizCount = 0
	}

	// Contar user_exams se a tabela existir
	if userExamsExists {
		err = r.DB.QueryRow(`SELECT COUNT(*) FROM user_exams WHERE user_id = $1`, userID).Scan(&examCount)
		if err != nil {
			log.Printf("PerformanceRepositoryPG: Erro ao contar user_exams: %v", err)
			examCount = 0
		}
	}

	report.TotalQuizzes = quizCount + examCount
	log.Printf("PerformanceRepositoryPG: Total de simulados (quiz: %d, exam: %d, total: %d)", quizCount, examCount, report.TotalQuizzes)

	// Verificar se a tabela user_answers existe
	log.Printf("PerformanceRepositoryPG: Verificando existência da tabela user_answers...")
	err = r.DB.QueryRow(`SELECT EXISTS (
		SELECT FROM information_schema.tables 
		WHERE table_schema = 'public' 
		AND table_name = 'user_answers'
	)`).Scan(&tableExists)
	if err != nil {
		log.Printf("PerformanceRepositoryPG: Erro ao verificar tabela user_answers: %v", err)
		return nil, err
	}

	log.Printf("PerformanceRepositoryPG: Tabela user_answers existe: %v", tableExists)

	if !tableExists {
		log.Printf("PerformanceRepositoryPG: Tabela user_answers não existe, definindo valores vazios")
		report.TotalQuestions = 0
		report.CorrectAnswers = 0
		report.Accuracy = 0.0
		return &report, nil
	}

	// Total de questões respondidas (tentando ambos os sistemas)
	log.Printf("PerformanceRepositoryPG: Contando questões respondidas para usuário %d", userID)

	var questionsFromQuiz, questionsFromExam int

	// Tentar contar questões do sistema de quiz (user_quiz)
	err = r.DB.QueryRow(`SELECT COUNT(*) FROM user_answers ua JOIN user_quiz uq ON ua.user_quiz_id = uq.id WHERE uq.user_id = $1`, userID).Scan(&questionsFromQuiz)
	if err != nil {
		log.Printf("PerformanceRepositoryPG: Erro ao contar questões do quiz (pode ser normal se não existir): %v", err)
		questionsFromQuiz = 0
	}

	// Tentar contar questões do sistema de exam (user_exams)
	if userExamsExists {
		err = r.DB.QueryRow(`SELECT COUNT(*) FROM user_answers ua JOIN user_exams ue ON ua.user_exam_id = ue.id WHERE ue.user_id = $1`, userID).Scan(&questionsFromExam)
		if err != nil {
			log.Printf("PerformanceRepositoryPG: Erro ao contar questões do exam: %v", err)
			questionsFromExam = 0
		}
	}

	report.TotalQuestions = questionsFromQuiz + questionsFromExam
	log.Printf("PerformanceRepositoryPG: Total de questões respondidas (quiz: %d, exam: %d, total: %d)", questionsFromQuiz, questionsFromExam, report.TotalQuestions)

	// Verificar se a tabela options existe
	log.Printf("PerformanceRepositoryPG: Verificando existência da tabela options...")
	err = r.DB.QueryRow(`SELECT EXISTS (
		SELECT FROM information_schema.tables 
		WHERE table_schema = 'public' 
		AND table_name = 'options'
	)`).Scan(&tableExists)
	if err != nil {
		log.Printf("PerformanceRepositoryPG: Erro ao verificar tabela options: %v", err)
		return nil, err
	}

	log.Printf("PerformanceRepositoryPG: Tabela options existe: %v", tableExists)

	if !tableExists {
		log.Printf("PerformanceRepositoryPG: Tabela options não existe, definindo acertos como 0")
		report.CorrectAnswers = 0
		report.Accuracy = 0.0
		return &report, nil
	}

	// Total de acertos (tentando ambos os sistemas)
	log.Printf("PerformanceRepositoryPG: Contando acertos para usuário %d", userID)

	var correctFromQuiz, correctFromExam int

	// Tentar contar acertos do sistema de quiz
	err = r.DB.QueryRow(`SELECT COUNT(*) FROM user_answers ua JOIN options o ON ua.option_id = o.id JOIN user_quiz uq ON ua.user_quiz_id = uq.id WHERE uq.user_id = $1 AND o.is_correct = true`, userID).Scan(&correctFromQuiz)
	if err != nil {
		log.Printf("PerformanceRepositoryPG: Erro ao contar acertos do quiz (pode ser normal se não existir): %v", err)
		correctFromQuiz = 0
	}

	// Tentar contar acertos do sistema de exam
	if userExamsExists {
		err = r.DB.QueryRow(`SELECT COUNT(*) FROM user_answers ua JOIN options o ON ua.option_id = o.id JOIN user_exams ue ON ua.user_exam_id = ue.id WHERE ue.user_id = $1 AND o.is_correct = true`, userID).Scan(&correctFromExam)
		if err != nil {
			log.Printf("PerformanceRepositoryPG: Erro ao contar acertos do exam: %v", err)
			correctFromExam = 0
		}
	}

	report.CorrectAnswers = correctFromQuiz + correctFromExam
	log.Printf("PerformanceRepositoryPG: Total de acertos (quiz: %d, exam: %d, total: %d)", correctFromQuiz, correctFromExam, report.CorrectAnswers)

	if report.TotalQuestions > 0 {
		report.Accuracy = float64(report.CorrectAnswers) / float64(report.TotalQuestions) * 100.0
		log.Printf("PerformanceRepositoryPG: Precisão calculada: %.2f%%", report.Accuracy)
	}

	log.Printf("PerformanceRepositoryPG: Relatório final - Simulados: %d, Questões: %d, Acertos: %d, Precisão: %.2f%%",
		report.TotalQuizzes, report.TotalQuestions, report.CorrectAnswers, report.Accuracy)

	return &report, nil
}
