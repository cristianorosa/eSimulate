package usecase

import (
	"time"

	"github.com/cristianorosa/eSimulate/backend/domain"
)

// UserExamUsecase implementa as regras de negócio para aplicação de exames
type UserExamUsecase struct {
	Repo            domain.UserExamRepository
	AnswerRepo      domain.UserAnswerRepository
	PerformanceRepo domain.PerformanceRepository
}

// StartExam inicia um exame para um usuário
func (uc *UserExamUsecase) StartExam(userID, examID int) (*domain.UserExam, error) {
	if userID == 0 || examID == 0 {
		return nil, domain.ErrInvalidInput
	}

	// Verifica se já existe um exame ativo para este usuário
	activeExam, err := uc.Repo.FindActiveByUser(userID, examID)
	if err == nil && activeExam != nil {
		return activeExam, nil // Retorna o exame ativo existente
	}

	userExam := &domain.UserExam{
		UserID:    userID,
		ExamID:    examID,
		StartedAt: time.Now(),
	}

	err = uc.Repo.Create(userExam)
	if err != nil {
		return nil, err
	}

	return userExam, nil
}

// SubmitAnswer submete uma resposta do usuário
func (uc *UserExamUsecase) SubmitAnswer(userExamID, questionID, optionID int, markedForReview bool) error {
	if userExamID == 0 || questionID == 0 {
		return domain.ErrInvalidInput
	}

	answer := &domain.UserAnswer{
		UserExamID:        userExamID,
		QuestionID:        questionID,
		OptionID:          &optionID,
		IsMarkedForReview: markedForReview,
		AnsweredAt:        time.Now(),
	}

	return uc.AnswerRepo.Create(answer)
}

// FinishExam finaliza um exame e calcula a pontuação
func (uc *UserExamUsecase) FinishExam(userExamID int) (*domain.UserExam, error) {
	if userExamID == 0 {
		return nil, domain.ErrInvalidInput
	}

	userExam, err := uc.Repo.FindByID(userExamID)
	if err != nil {
		return nil, err
	}

	// Busca todas as respostas do usuário
	answers, err := uc.AnswerRepo.ListByUserExam(userExamID)
	if err != nil {
		return nil, err
	}

	// Calcula pontuação total
	totalQuestions := len(answers)
	correctAnswers := 0

	for _, answer := range answers {
		if answer.IsCorrect != nil && *answer.IsCorrect {
			correctAnswers++
		}
	}

	score := float64(0)
	if totalQuestions > 0 {
		score = (float64(correctAnswers) / float64(totalQuestions)) * 100
	}

	// Calcula tempo gasto
	timeSpent := int(time.Since(userExam.StartedAt).Minutes())

	// Atualiza o exame
	now := time.Now()
	userExam.FinishedAt = &now
	userExam.TotalScore = &score
	userExam.TimeSpentMinutes = &timeSpent

	// Determina se passou (assumindo 70% como padrão)
	passed := score >= 70.0
	userExam.Passed = &passed

	err = uc.Repo.Update(userExam)
	if err != nil {
		return nil, err
	}

	return userExam, nil
}

// GetUserExam busca um exame aplicado por ID
func (uc *UserExamUsecase) GetUserExam(id int) (*domain.UserExam, error) {
	if id == 0 {
		return nil, domain.ErrInvalidInput
	}
	return uc.Repo.FindByID(id)
}

// ListByUser lista todos os exames aplicados por um usuário
func (uc *UserExamUsecase) ListByUser(userID int) ([]*domain.UserExam, error) {
	if userID == 0 {
		return nil, domain.ErrInvalidInput
	}
	return uc.Repo.ListByUser(userID)
}
