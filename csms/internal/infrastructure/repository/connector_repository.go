package repository

import (
	"context"

	"github.com/malikkhoiri/csms/internal/domain"
	"gorm.io/gorm"
)

type ConnectorRepository struct {
	db *gorm.DB
}

func NewConnectorRepository(db *gorm.DB) domain.ConnectorRepository {
	return &ConnectorRepository{db: db}
}

func (r *ConnectorRepository) Create(ctx context.Context, connector *domain.Connector) error {
	return r.db.WithContext(ctx).Create(connector).Error
}

func (r *ConnectorRepository) GetByID(ctx context.Context, id uint) (*domain.Connector, error) {
	var connector domain.Connector
	err := r.db.WithContext(ctx).First(&connector, id).Error
	if err != nil {
		return nil, err
	}
	return &connector, nil
}

func (r *ConnectorRepository) GetByChargePointAndConnectorID(ctx context.Context, chargePointID uint, connectorID int) (*domain.Connector, error) {
	var connector domain.Connector
	err := r.db.WithContext(ctx).Where("charge_point_id = ? AND connector_id = ?", chargePointID, connectorID).First(&connector).Error
	if err != nil {
		return nil, err
	}
	return &connector, nil
}

func (r *ConnectorRepository) Update(ctx context.Context, connector *domain.Connector) error {
	return r.db.WithContext(ctx).Save(connector).Error
}

func (r *ConnectorRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Connector{}, id).Error
}

func (r *ConnectorRepository) ListByChargePoint(ctx context.Context, chargePointID uint) ([]domain.Connector, error) {
	var connectors []domain.Connector
	err := r.db.WithContext(ctx).Where("charge_point_id = ?", chargePointID).Find(&connectors).Error
	return connectors, err
}

func (r *ConnectorRepository) UpdateStatus(ctx context.Context, id uint, status string) error {
	return r.db.WithContext(ctx).Model(&domain.Connector{}).Where("id = ?", id).Update("status", status).Error
}
