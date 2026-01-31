package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

type CategoryService struct {
	repo *repositories.CategoryRepository
}

func NewCategoryService(repo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetAll() ([]models.Category, error) {
	return s.repo.GetAll()
}

func (s *CategoryService) CreateCategory(data *models.Category) error {
	return s.repo.CreateCategory(*data)
}

func (s *CategoryService) GetByID(id int) (*models.Category, error) {
	return s.repo.GetByID(id)
}

func (s *CategoryService) UpdateByID(id int, data *models.Category) error {
	return s.repo.UpdateByID(id, *data)
}

func (s *CategoryService) DeleteByID(id int) error {
	return s.repo.DeleteByID(id)
}
