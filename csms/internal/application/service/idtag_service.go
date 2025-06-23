package service

import (
	"context"
	"errors"
	"time"

	"github.com/malikkhoiri/csms/internal/domain"
)

type IDTagService struct {
	idTagRepo domain.IDTagRepository
}

func NewIDTagService(idTagRepo domain.IDTagRepository) domain.IDTagService {
	return &IDTagService{
		idTagRepo: idTagRepo,
	}
}

func (s *IDTagService) Authorize(ctx context.Context, request *domain.AuthorizeRequest) (*domain.AuthorizeResponse, error) {
	idTag, err := s.idTagRepo.GetByTag(ctx, request.IDTag)
	if err != nil {
		return &domain.AuthorizeResponse{
			IDTagInfo: domain.IDTagInfo{
				Status: domain.AuthorizeStatusInvalid,
			},
		}, nil
	}

	if idTag.Status != domain.AuthorizeStatusAccepted {
		return &domain.AuthorizeResponse{
			IDTagInfo: domain.IDTagInfo{
				Status: idTag.Status,
			},
		}, nil
	}

	expiryDate := idTag.ExpiryDate.UTC()
	if !idTag.ExpiryDate.IsZero() && idTag.ExpiryDate.Before(time.Now()) {
		return &domain.AuthorizeResponse{
			IDTagInfo: domain.IDTagInfo{
				Status:     domain.AuthorizeStatusExpired,
				ExpiryDate: &expiryDate,
			},
		}, nil
	}

	return &domain.AuthorizeResponse{
		IDTagInfo: domain.IDTagInfo{
			Status:     domain.AuthorizeStatusAccepted,
			ExpiryDate: &expiryDate,
		},
	}, nil
}

func (s *IDTagService) CreateIDTag(ctx context.Context, idTag *domain.IDTag) error {
	existingTag, err := s.idTagRepo.GetByTag(ctx, idTag.Tag)
	if err == nil && existingTag != nil {
		return errors.New("IDTag already exists")
	}

	return s.idTagRepo.Create(ctx, idTag)
}

func (s *IDTagService) GetIDTag(ctx context.Context, id uint) (*domain.IDTag, error) {
	return s.idTagRepo.GetByID(ctx, id)
}

func (s *IDTagService) GetByTag(ctx context.Context, tag string) (*domain.IDTag, error) {
	return s.idTagRepo.GetByTag(ctx, tag)
}

func (s *IDTagService) UpdateIDTag(ctx context.Context, idTag *domain.IDTag) error {
	existingTag, err := s.idTagRepo.GetByID(ctx, idTag.ID)
	if err != nil {
		return errors.New("IDTag not found")
	}

	if existingTag.Tag != idTag.Tag {
		conflictingTag, err := s.idTagRepo.GetByTag(ctx, idTag.Tag)
		if err == nil && conflictingTag != nil {
			return errors.New("IDTag with this tag already exists")
		}
	}

	return s.idTagRepo.Update(ctx, idTag)
}

func (s *IDTagService) DeleteIDTag(ctx context.Context, id uint) error {
	// Check if IDTag exists
	_, err := s.idTagRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("IDTag not found")
	}

	return s.idTagRepo.Delete(ctx, id)
}

func (s *IDTagService) ListIDTags(ctx context.Context, limit, offset int) ([]domain.IDTag, error) {
	return s.idTagRepo.List(ctx, limit, offset)
}

func (s *IDTagService) ListByUser(ctx context.Context, userID uint) ([]domain.IDTag, error) {
	return s.idTagRepo.ListByUser(ctx, userID)
}
