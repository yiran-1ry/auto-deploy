package service

import (
	"fmt"

	"github.com/yiran-1ry/auto-deploy/backend/models"
	"github.com/yiran-1ry/auto-deploy/backend/repository"
)

type AppService interface {
	ListServices(page, pageSize int, query map[string]interface{}) ([]models.Service, int64, error)
	GetServiceByID(id uint) (*models.Service, error)
	CreateService(service *models.Service) error
	UpdateService(id uint, service *models.Service) error
	DeleteService(id uint) error
}

type appService struct {
	serviceRepo repository.ServiceRepository
	serverRepo  repository.ServerRepository
}

func NewAppService(serviceRepo repository.ServiceRepository, serverRepo repository.ServerRepository) AppService {
	return &appService{serviceRepo: serviceRepo, serverRepo: serverRepo}
}

func (s *appService) ListServices(page, pageSize int, query map[string]interface{}) ([]models.Service, int64, error) {
	return s.serviceRepo.List(page, pageSize, query)
}

func (s *appService) GetServiceByID(id uint) (*models.Service, error) {
	return s.serviceRepo.FindByID(id)
}

func (s *appService) CreateService(service *models.Service) error {
	if service.ServerID == 0 {
		return fmt.Errorf("必须关联服务器")
	}
	if _, err := s.serverRepo.FindByID(service.ServerID); err != nil {
		return fmt.Errorf("服务器不存在: %w", err)
	}
	return s.serviceRepo.Create(service)
}

func (s *appService) UpdateService(id uint, service *models.Service) error {
	_, err := s.serviceRepo.FindByID(id)
	if err != nil {
		return err
	}
	if service.ServerID == 0 {
		return fmt.Errorf("必须关联服务器")
	}
	if _, err := s.serverRepo.FindByID(service.ServerID); err != nil {
		return fmt.Errorf("服务器不存在: %w", err)
	}
	service.ID = uint64(id)
	return s.serviceRepo.Update(service)
}

func (s *appService) DeleteService(id uint) error {
	return s.serviceRepo.Delete(id)
}
