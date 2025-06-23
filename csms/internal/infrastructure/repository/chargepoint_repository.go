package repository

import (
	"context"
	"time"

	"github.com/malikkhoiri/csms/internal/domain"
	"gorm.io/gorm"
)

type ChargePointRepository struct {
	db *gorm.DB
}

func NewChargePointRepository(db *gorm.DB) domain.ChargePointRepository {
	return &ChargePointRepository{db: db}
}

func (r *ChargePointRepository) Create(ctx context.Context, cp *domain.ChargePoint) error {
	return r.db.WithContext(ctx).Create(cp).Error
}

func (r *ChargePointRepository) GetByID(ctx context.Context, id uint) (*domain.ChargePoint, error) {
	var cp domain.ChargePoint
	err := r.db.WithContext(ctx).Preload("Connectors").First(&cp, id).Error
	if err != nil {
		return nil, err
	}
	return &cp, nil
}

func (r *ChargePointRepository) GetByCode(ctx context.Context, code string) (*domain.ChargePoint, error) {
	var cp domain.ChargePoint
	err := r.db.WithContext(ctx).Preload("Connectors").Where("charge_point_code = ?", code).First(&cp).Error
	if err != nil {
		return nil, err
	}
	return &cp, nil
}

func (r *ChargePointRepository) Update(ctx context.Context, cp *domain.ChargePoint) error {
	return r.db.WithContext(ctx).Save(cp).Error
}

func (r *ChargePointRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.ChargePoint{}, id).Error
}

func (r *ChargePointRepository) List(ctx context.Context, limit, offset int) ([]domain.ChargePoint, error) {
	var cps []domain.ChargePoint
	err := r.db.WithContext(ctx).Preload("Connectors").Limit(limit).Offset(offset).Find(&cps).Error
	return cps, err
}

func (r *ChargePointRepository) UpdateStatus(ctx context.Context, id uint, status string) error {
	return r.db.WithContext(ctx).Model(&domain.ChargePoint{}).Where("id = ?", id).Update("status", status).Error
}

func (r *ChargePointRepository) UpdateHeartbeat(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&domain.ChargePoint{}).Where("id = ?", id).Update("last_heartbeat", time.Now()).Error
}
