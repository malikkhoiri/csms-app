package repository

import (
	"context"

	"github.com/malikkhoiri/csms/internal/domain"
	"gorm.io/gorm"
)

type IDTagRepository struct {
	db *gorm.DB
}

func NewIDTagRepository(db *gorm.DB) domain.IDTagRepository {
	return &IDTagRepository{db: db}
}

func (r *IDTagRepository) Create(ctx context.Context, idTag *domain.IDTag) error {
	return r.db.WithContext(ctx).Create(idTag).Error
}

func (r *IDTagRepository) GetByID(ctx context.Context, id uint) (*domain.IDTag, error) {
	var idTag domain.IDTag
	err := r.db.WithContext(ctx).Preload("User").First(&idTag, id).Error
	if err != nil {
		return nil, err
	}
	return &idTag, nil
}

func (r *IDTagRepository) GetByTag(ctx context.Context, tag string) (*domain.IDTag, error) {
	var idTag domain.IDTag
	err := r.db.WithContext(ctx).Preload("User").Where("tag = ?", tag).First(&idTag).Error
	if err != nil {
		return nil, err
	}
	return &idTag, nil
}

func (r *IDTagRepository) Update(ctx context.Context, idTag *domain.IDTag) error {
	return r.db.WithContext(ctx).Save(idTag).Error
}

func (r *IDTagRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.IDTag{}, id).Error
}

func (r *IDTagRepository) List(ctx context.Context, limit, offset int) ([]domain.IDTag, error) {
	var idTags []domain.IDTag
	err := r.db.WithContext(ctx).Preload("User").Limit(limit).Offset(offset).Find(&idTags).Error
	return idTags, err
}

func (r *IDTagRepository) ListByUser(ctx context.Context, userID uint) ([]domain.IDTag, error) {
	var idTags []domain.IDTag
	err := r.db.WithContext(ctx).Preload("User").Where("user_id = ?", userID).Find(&idTags).Error
	return idTags, err
}
