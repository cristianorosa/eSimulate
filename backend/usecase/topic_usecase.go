package usecase

import (
	"context"
	"errors"

	"github.com/cristianorosa/eSimulate/backend/domain"
)

// TopicUsecase implementa as regras de negócio para tópicos
type TopicUsecase struct {
	Repo domain.TopicRepository
}

// NewTopicUsecase cria uma nova instância de TopicUsecase
func NewTopicUsecase(repo domain.TopicRepository) *TopicUsecase {
	return &TopicUsecase{Repo: repo}
}

// CreateTopic cria um novo tópico
func (uc *TopicUsecase) CreateTopic(ctx context.Context, topic *domain.Topic) error {
	if topic.Name == "" {
		return errors.New("nome do tópico é obrigatório")
	}

	if topic.WeightPercentage < 0 || topic.WeightPercentage > 100 {
		return errors.New("peso do tópico deve estar entre 0 e 100")
	}

	if topic.QuestionsCount < 0 {
		return errors.New("quantidade de questões deve ser maior ou igual a 0")
	}

	return uc.Repo.Create(topic)
}

// UpdateTopic atualiza um tópico existente
func (uc *TopicUsecase) UpdateTopic(ctx context.Context, topic *domain.Topic) error {
	if topic.Name == "" {
		return errors.New("nome do tópico é obrigatório")
	}

	if topic.WeightPercentage < 0 || topic.WeightPercentage > 100 {
		return errors.New("peso do tópico deve estar entre 0 e 100")
	}

	if topic.QuestionsCount < 0 {
		return errors.New("quantidade de questões deve ser maior ou igual a 0")
	}

	return uc.Repo.Update(topic)
}

// DeleteTopic remove um tópico
func (uc *TopicUsecase) DeleteTopic(ctx context.Context, id int) error {
	return uc.Repo.Delete(id)
}

// GetTopic obtém um tópico por ID
func (uc *TopicUsecase) GetTopic(ctx context.Context, id int) (*domain.Topic, error) {
	return uc.Repo.FindByID(id)
}

// ListTopicsByExam lista todos os tópicos de um exame
func (uc *TopicUsecase) ListTopicsByExam(ctx context.Context, examID int) ([]*domain.Topic, error) {
	return uc.Repo.ListByExam(examID)
}

// ListAllTopics lista todos os tópicos
func (uc *TopicUsecase) ListAllTopics(ctx context.Context) ([]*domain.Topic, error) {
	return uc.Repo.ListAll()
}

// ListTopicsPaginated lista tópicos com paginação
func (uc *TopicUsecase) ListTopicsPaginated(ctx context.Context, page, pageSize int, examID *int) ([]*domain.Topic, *domain.Pagination, error) {
	topics, pagination, err := uc.Repo.ListPaginated(page, pageSize, examID)
	if err != nil {
		return nil, nil, err
	}
	return topics, pagination, nil
}
