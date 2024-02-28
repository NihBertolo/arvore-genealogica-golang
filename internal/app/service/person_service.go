package service

import (
	"arvore-genealogica-golang/internal/app/domain/models"
	"arvore-genealogica-golang/internal/app/repository"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type PersonService struct {
	repository *repository.PersonRepository
}

func NewPersonService(repo *repository.PersonRepository) *PersonService {
	return &PersonService{repo}
}

func (s *PersonService) Create(ctx context.Context, data *models.Person) error {
	data.ID = uuid.New()
	return s.repository.Create(ctx, data)
}

func (s *PersonService) Delete(ctx context.Context, id uuid.UUID) error {
	stringID := id.String()
	data, err := s.repository.FindByID(stringID)
	if err != nil {
		fmt.Printf("Erro ao deletar o usuário: %s, error: %s", id, err)
		return err
	}
	return s.repository.DeleteNode(data.ID.String())
}

func (s *PersonService) FindByID(ctx context.Context, id uuid.UUID) (*models.Person, error) {
	data, err := s.repository.FindByID(id.String())
	if err != nil {
		fmt.Printf("Erro ao procurar o usuário: %s, error: %s", id, err)
		return nil, err
	}
	return data, nil
}

func (s *PersonService) Update(ctx context.Context, data *models.Person) error {
	return s.repository.UpdatePerson(data)
}
