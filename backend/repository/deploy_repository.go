package repository

import (
	"errors"
	"fmt"

	"github.com/yiran-1ry/auto-deploy/backend/models"
	"gorm.io/gorm"
)

type DeployRepository interface {
	Create(deploy *models.Deploy) error
	FindByID(id uint) (*models.Deploy, error)
	Update(deploy *models.Deploy) error
	Delete(id uint) error
	List(page, pageSize int, serviceID uint, status string) ([]models.Deploy, int64, error)
	CreateLog(log *models.DeployLog) error
	GetLogs(deployID uint) ([]models.DeployLog, error)
	GetFileByID(id uint) (*models.File, error)
	GetFileByServiceID(serviceID uint) (*models.File, error)
	GetContainerByServiceID(serviceID uint) (*models.Container, error)
	GetServerByID(id uint) (*models.Server, error)
	GetServiceByID(id uint) (*models.Service, error)
	// ResolveServerForAutoDeploy 仅使用 services.server_id 对应的服务器（与创建服务时的绑定一致）
	ResolveServerForAutoDeploy(serviceID uint) (*models.Server, error)
	UpdateStatus(id uint, status string) error
}

type deployRepository struct {
	db *gorm.DB
}

func NewDeployRepository(db *gorm.DB) DeployRepository {
	return &deployRepository{db: db}
}

func (r *deployRepository) Create(deploy *models.Deploy) error {
	return r.db.Create(deploy).Error
}

func (r *deployRepository) FindByID(id uint) (*models.Deploy, error) {
	var deploy models.Deploy
	err := r.db.Where("id = ?", id).First(&deploy).Error
	return &deploy, err
}

func (r *deployRepository) Update(deploy *models.Deploy) error {
	return r.db.Save(deploy).Error
}

func (r *deployRepository) Delete(id uint) error {
	return r.db.Delete(&models.Deploy{}, id).Error
}

func (r *deployRepository) List(page, pageSize int, serviceID uint, status string) ([]models.Deploy, int64, error) {
	var deploys []models.Deploy
	var total int64

	query := r.db.Model(&models.Deploy{})

	if serviceID > 0 {
		query = query.Where("service_id = ?", serviceID)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("create_time DESC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&deploys).Error

	return deploys, total, err
}

func (r *deployRepository) CreateLog(log *models.DeployLog) error {
	return r.db.Create(log).Error
}

func (r *deployRepository) GetLogs(deployID uint) ([]models.DeployLog, error) {
	var logs []models.DeployLog
	err := r.db.Where("deploy_id = ?", deployID).
		Order("create_time DESC").
		Limit(10).
		Find(&logs).Error
	return logs, err
}

func (r *deployRepository) GetFileByID(id uint) (*models.File, error) {
	var file models.File
	err := r.db.First(&file, id).Error
	return &file, err
}

func (r *deployRepository) GetFileByServiceID(serviceID uint) (*models.File, error) {
	var file models.File
	err := r.db.Where("service_id = ?", serviceID).
		Order("updated_at DESC").
		First(&file).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *deployRepository) GetContainerByServiceID(serviceID uint) (*models.Container, error) {
	var container models.Container
	// CI notify 刚插入的记录 id 最大；用 id DESC 比 updated_at 更稳
	err := r.db.Where("service_id = ?", serviceID).
		Order("id DESC").
		First(&container).Error
	if err != nil {
		return nil, err
	}
	return &container, nil
}

func (r *deployRepository) GetServerByID(id uint) (*models.Server, error) {
	var server models.Server
	err := r.db.First(&server, id).Error
	return &server, err
}

func (r *deployRepository) GetServiceByID(id uint) (*models.Service, error) {
	var svc models.Service
	err := r.db.First(&svc, id).Error
	return &svc, err
}

func (r *deployRepository) ResolveServerForAutoDeploy(serviceID uint) (*models.Server, error) {
	svc, err := r.GetServiceByID(serviceID)
	if err != nil {
		return nil, err
	}
	if svc.ServerID == 0 {
		return nil, fmt.Errorf("服务未绑定服务器(server_id=0)，请在服务管理中配置")
	}
	srv, err := r.GetServerByID(svc.ServerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("绑定的服务器不存在或已删除 (server_id=%d)", svc.ServerID)
		}
		return nil, err
	}
	return srv, nil
}

func (r *deployRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&models.Deploy{}).Where("id = ?", id).Update("status", status).Error
}
