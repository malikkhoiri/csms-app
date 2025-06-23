package repository

import (
	"context"

	"github.com/malikkhoiri/csms/internal/domain"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) domain.TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(ctx context.Context, transaction *domain.Transaction) error {
	return r.db.WithContext(ctx).Create(transaction).Error
}

func (r *TransactionRepository) GetByID(ctx context.Context, id uint) (*domain.Transaction, error) {
	var transaction domain.Transaction
	err := r.db.WithContext(ctx).Preload("ChargePoint").Preload("IDTag").First(&transaction, id).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *TransactionRepository) GetByTransactionID(ctx context.Context, transactionID int) (*domain.Transaction, error) {
	var transaction domain.Transaction
	err := r.db.WithContext(ctx).Preload("ChargePoint").Preload("IDTag").Where("transaction_id = ?", transactionID).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *TransactionRepository) Update(ctx context.Context, transaction *domain.Transaction) error {
	return r.db.WithContext(ctx).Save(transaction).Error
}

func (r *TransactionRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Transaction{}, id).Error
}

func (r *TransactionRepository) List(ctx context.Context, limit, offset int) ([]domain.Transaction, error) {
	var transactions []domain.Transaction
	err := r.db.WithContext(ctx).Preload("ChargePoint").Preload("IDTag").Limit(limit).Offset(offset).Find(&transactions).Error
	return transactions, err
}

func (r *TransactionRepository) ListByChargePoint(ctx context.Context, chargePointID uint) ([]domain.Transaction, error) {
	var transactions []domain.Transaction
	err := r.db.WithContext(ctx).Preload("ChargePoint").Preload("IDTag").Where("charge_point_id = ?", chargePointID).Find(&transactions).Error
	return transactions, err
}

func (r *TransactionRepository) ListByUser(ctx context.Context, idTag string) ([]domain.Transaction, error) {
	var transactions []domain.Transaction
	err := r.db.WithContext(ctx).Preload("ChargePoint").Preload("IDTag").Where("id_tag_id = ?", idTag).Find(&transactions).Error
	return transactions, err
}

func (r *TransactionRepository) GetActiveByConnector(ctx context.Context, chargePointID uint, connectorID int) (*domain.Transaction, error) {
	var transaction domain.Transaction
	err := r.db.WithContext(ctx).Preload("ChargePoint").Preload("IDTag").
		Where("charge_point_id = ? AND connector_id = ? AND status = ?", chargePointID, connectorID, domain.TransactionStatusActive).
		First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}
