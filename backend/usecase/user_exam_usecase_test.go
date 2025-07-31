package usecase

import (
	"testing"
	"time"

	"github.com/cristianorosa/eSimulate/backend/domain"
)

// MockUserExamRepository implementa UserExamRepository para testes
type MockUserExamRepository struct {
	userExams map[int]*domain.UserExam
	nextID    int
}

func NewMockUserExamRepository() *MockUserExamRepository {
	return &MockUserExamRepository{
		userExams: make(map[int]*domain.UserExam),
		nextID:    1,
	}
}

func (m *MockUserExamRepository) Create(userExam *domain.UserExam) error {
	userExam.ID = m.nextID
	m.userExams[userExam.ID] = userExam
	m.nextID++
	return nil
}

func (m *MockUserExamRepository) Update(userExam *domain.UserExam) error {
	if _, exists := m.userExams[userExam.ID]; !exists {
		return domain.ErrNotFound
	}
	m.userExams[userExam.ID] = userExam
	return nil
}

func (m *MockUserExamRepository) FindByID(id int) (*domain.UserExam, error) {
	if userExam, exists := m.userExams[id]; exists {
		return userExam, nil
	}
	return nil, domain.ErrNotFound
}

func (m *MockUserExamRepository) ListByUser(userID int) ([]*domain.UserExam, error) {
	var userExams []*domain.UserExam
	for _, userExam := range m.userExams {
		if userExam.UserID == userID {
			userExams = append(userExams, userExam)
		}
	}
	return userExams, nil
}

func (m *MockUserExamRepository) ListByExam(examID int) ([]*domain.UserExam, error) {
	var userExams []*domain.UserExam
	for _, userExam := range m.userExams {
		if userExam.ExamID == examID {
			userExams = append(userExams, userExam)
		}
	}
	return userExams, nil
}

func (m *MockUserExamRepository) FindActiveByUser(userID int, examID int) (*domain.UserExam, error) {
	for _, userExam := range m.userExams {
		if userExam.UserID == userID && userExam.ExamID == examID && userExam.FinishedAt == nil {
			return userExam, nil
		}
	}
	return nil, nil // Não há exame ativo
}

// MockUserAnswerRepository implementa UserAnswerRepository para testes
type MockUserAnswerRepository struct {
	userAnswers map[int]*domain.UserAnswer
	nextID      int
}

func NewMockUserAnswerRepository() *MockUserAnswerRepository {
	return &MockUserAnswerRepository{
		userAnswers: make(map[int]*domain.UserAnswer),
		nextID:      1,
	}
}

func (m *MockUserAnswerRepository) Create(userAnswer *domain.UserAnswer) error {
	userAnswer.ID = m.nextID
	m.userAnswers[userAnswer.ID] = userAnswer
	m.nextID++
	return nil
}

func (m *MockUserAnswerRepository) Update(userAnswer *domain.UserAnswer) error {
	if _, exists := m.userAnswers[userAnswer.ID]; !exists {
		return domain.ErrNotFound
	}
	m.userAnswers[userAnswer.ID] = userAnswer
	return nil
}

func (m *MockUserAnswerRepository) ListByUserExam(userExamID int) ([]*domain.UserAnswer, error) {
	var userAnswers []*domain.UserAnswer
	for _, userAnswer := range m.userAnswers {
		if userAnswer.UserExamID == userExamID {
			userAnswers = append(userAnswers, userAnswer)
		}
	}
	return userAnswers, nil
}

func (m *MockUserAnswerRepository) FindByQuestion(userExamID int, questionID int) (*domain.UserAnswer, error) {
	for _, userAnswer := range m.userAnswers {
		if userAnswer.UserExamID == userExamID && userAnswer.QuestionID == questionID {
			return userAnswer, nil
		}
	}
	return nil, domain.ErrNotFound
}

// MockDomainPerformanceRepository implementa DomainPerformanceRepository para testes
type MockDomainPerformanceRepository struct {
	performances map[int]*domain.Performance
	nextID       int
}

func NewMockDomainPerformanceRepository() *MockDomainPerformanceRepository {
	return &MockDomainPerformanceRepository{
		performances: make(map[int]*domain.Performance),
		nextID:       1,
	}
}

func (m *MockDomainPerformanceRepository) Create(performance *domain.Performance) error {
	performance.ID = m.nextID
	m.performances[performance.ID] = performance
	m.nextID++
	return nil
}

func (m *MockDomainPerformanceRepository) Update(performance *domain.Performance) error {
	if _, exists := m.performances[performance.ID]; !exists {
		return domain.ErrNotFound
	}
	m.performances[performance.ID] = performance
	return nil
}

func (m *MockDomainPerformanceRepository) ListByUserExam(userExamID int) ([]*domain.Performance, error) {
	var performances []*domain.Performance
	for _, performance := range m.performances {
		if performance.UserExamID == userExamID {
			performances = append(performances, performance)
		}
	}
	return performances, nil
}

func TestUserExamUsecase_StartExam(t *testing.T) {
	mockUserExamRepo := NewMockUserExamRepository()
	mockUserAnswerRepo := NewMockUserAnswerRepository()
	mockPerformanceRepo := NewMockDomainPerformanceRepository()

	uc := &UserExamUsecase{
		Repo:            mockUserExamRepo,
		AnswerRepo:      mockUserAnswerRepo,
		PerformanceRepo: mockPerformanceRepo,
	}

	tests := []struct {
		name        string
		userID      int
		examID      int
		expectError bool
	}{
		{
			name:        "Iniciar exame válido",
			userID:      1,
			examID:      1,
			expectError: false,
		},
		{
			name:        "Iniciar exame com userID inválido",
			userID:      0,
			examID:      1,
			expectError: true,
		},
		{
			name:        "Iniciar exame com examID inválido",
			userID:      1,
			examID:      0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := uc.StartExam(tt.userID, tt.examID)

			if tt.expectError {
				if err == nil {
					t.Errorf("Esperava erro, mas não houve erro")
				}
			} else {
				if err != nil {
					t.Errorf("Não esperava erro, mas houve: %v", err)
				}
				if result == nil {
					t.Errorf("Esperava UserExam, mas recebeu nil")
				}
				if result.UserID != tt.userID {
					t.Errorf("Esperava UserID %d, mas recebeu %d", tt.userID, result.UserID)
				}
				if result.ExamID != tt.examID {
					t.Errorf("Esperava ExamID %d, mas recebeu %d", tt.examID, result.ExamID)
				}
			}
		})
	}
}

func TestUserExamUsecase_SubmitAnswer(t *testing.T) {
	mockUserExamRepo := NewMockUserExamRepository()
	mockUserAnswerRepo := NewMockUserAnswerRepository()
	mockPerformanceRepo := NewMockDomainPerformanceRepository()

	uc := &UserExamUsecase{
		Repo:            mockUserExamRepo,
		AnswerRepo:      mockUserAnswerRepo,
		PerformanceRepo: mockPerformanceRepo,
	}

	tests := []struct {
		name            string
		userExamID      int
		questionID      int
		optionID        int
		markedForReview bool
		expectError     bool
	}{
		{
			name:            "Submeter resposta válida",
			userExamID:      1,
			questionID:      1,
			optionID:        1,
			markedForReview: false,
			expectError:     false,
		},
		{
			name:            "Submeter resposta marcada para revisão",
			userExamID:      1,
			questionID:      2,
			optionID:        2,
			markedForReview: true,
			expectError:     false,
		},
		{
			name:            "Submeter resposta com userExamID inválido",
			userExamID:      0,
			questionID:      1,
			optionID:        1,
			markedForReview: false,
			expectError:     true,
		},
		{
			name:            "Submeter resposta com questionID inválido",
			userExamID:      1,
			questionID:      0,
			optionID:        1,
			markedForReview: false,
			expectError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := uc.SubmitAnswer(tt.userExamID, tt.questionID, tt.optionID, tt.markedForReview)

			if tt.expectError {
				if err == nil {
					t.Errorf("Esperava erro, mas não houve erro")
				}
			} else {
				if err != nil {
					t.Errorf("Não esperava erro, mas houve: %v", err)
				}
			}
		})
	}
}

func TestUserExamUsecase_FinishExam(t *testing.T) {
	mockUserExamRepo := NewMockUserExamRepository()
	mockUserAnswerRepo := NewMockUserAnswerRepository()
	mockPerformanceRepo := NewMockDomainPerformanceRepository()

	uc := &UserExamUsecase{
		Repo:            mockUserExamRepo,
		AnswerRepo:      mockUserAnswerRepo,
		PerformanceRepo: mockPerformanceRepo,
	}

	// Criar um exame aplicado primeiro
	userExam := &domain.UserExam{
		UserID:    1,
		ExamID:    1,
		StartedAt: time.Now().Add(-30 * time.Minute), // 30 minutos atrás
	}
	mockUserExamRepo.Create(userExam)

	// Adicionar algumas respostas
	answers := []*domain.UserAnswer{
		{
			UserExamID: userExam.ID,
			QuestionID: 1,
			OptionID:   &[]int{1}[0],
			IsCorrect:  &[]bool{true}[0],
		},
		{
			UserExamID: userExam.ID,
			QuestionID: 2,
			OptionID:   &[]int{2}[0],
			IsCorrect:  &[]bool{false}[0],
		},
		{
			UserExamID: userExam.ID,
			QuestionID: 3,
			OptionID:   &[]int{1}[0],
			IsCorrect:  &[]bool{true}[0],
		},
	}

	for _, answer := range answers {
		mockUserAnswerRepo.Create(answer)
	}

	result, err := uc.FinishExam(userExam.ID)

	if err != nil {
		t.Errorf("Não esperava erro, mas houve: %v", err)
	}

	if result == nil {
		t.Errorf("Esperava UserExam, mas recebeu nil")
	}

	if result.FinishedAt == nil {
		t.Errorf("Esperava que FinishedAt fosse definido")
	}

	if result.TotalScore == nil {
		t.Errorf("Esperava que TotalScore fosse definido")
	}

	if result.Passed == nil {
		t.Errorf("Esperava que Passed fosse definido")
	}

	// Verificar se a pontuação está correta (2 de 3 = 66.67%)
	expectedScore := 66.67
	if *result.TotalScore < expectedScore-1 || *result.TotalScore > expectedScore+1 {
		t.Errorf("Esperava pontuação próxima a %.2f, mas recebeu %.2f", expectedScore, *result.TotalScore)
	}
}

func TestUserExamUsecase_GetUserExam(t *testing.T) {
	mockUserExamRepo := NewMockUserExamRepository()
	mockUserAnswerRepo := NewMockUserAnswerRepository()
	mockPerformanceRepo := NewMockDomainPerformanceRepository()

	uc := &UserExamUsecase{
		Repo:            mockUserExamRepo,
		AnswerRepo:      mockUserAnswerRepo,
		PerformanceRepo: mockPerformanceRepo,
	}

	// Criar um exame aplicado primeiro
	userExam := &domain.UserExam{
		UserID:    1,
		ExamID:    1,
		StartedAt: time.Now(),
	}
	mockUserExamRepo.Create(userExam)

	tests := []struct {
		name        string
		userExamID  int
		expectError bool
	}{
		{
			name:        "Buscar exame aplicado existente",
			userExamID:  userExam.ID,
			expectError: false,
		},
		{
			name:        "Buscar exame aplicado inexistente",
			userExamID:  999,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := uc.GetUserExam(tt.userExamID)

			if tt.expectError {
				if err == nil {
					t.Errorf("Esperava erro, mas não houve erro")
				}
			} else {
				if err != nil {
					t.Errorf("Não esperava erro, mas houve: %v", err)
				}
				if result == nil {
					t.Errorf("Esperava UserExam, mas recebeu nil")
				}
				if result.ID != tt.userExamID {
					t.Errorf("Esperava ID %d, mas recebeu %d", tt.userExamID, result.ID)
				}
			}
		})
	}
}

func TestUserExamUsecase_ListByUser(t *testing.T) {
	mockUserExamRepo := NewMockUserExamRepository()
	mockUserAnswerRepo := NewMockUserAnswerRepository()
	mockPerformanceRepo := NewMockDomainPerformanceRepository()

	uc := &UserExamUsecase{
		Repo:            mockUserExamRepo,
		AnswerRepo:      mockUserAnswerRepo,
		PerformanceRepo: mockPerformanceRepo,
	}

	// Criar alguns exames aplicados
	userExams := []*domain.UserExam{
		{UserID: 1, ExamID: 1, StartedAt: time.Now()},
		{UserID: 1, ExamID: 2, StartedAt: time.Now()},
		{UserID: 2, ExamID: 1, StartedAt: time.Now()},
	}

	for _, userExam := range userExams {
		mockUserExamRepo.Create(userExam)
	}

	tests := []struct {
		name          string
		userID        int
		expectedCount int
		expectError   bool
	}{
		{
			name:          "Listar exames do usuário 1",
			userID:        1,
			expectedCount: 2,
			expectError:   false,
		},
		{
			name:          "Listar exames do usuário 2",
			userID:        2,
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "Listar exames do usuário inexistente",
			userID:        999,
			expectedCount: 0,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := uc.ListByUser(tt.userID)

			if tt.expectError {
				if err == nil {
					t.Errorf("Esperava erro, mas não houve erro")
				}
			} else {
				if err != nil {
					t.Errorf("Não esperava erro, mas houve: %v", err)
				}
				if len(result) != tt.expectedCount {
					t.Errorf("Esperava %d exames, mas recebeu %d", tt.expectedCount, len(result))
				}
			}
		})
	}
}
