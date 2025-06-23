package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/malikkhoiri/csms/internal/config"
	"github.com/malikkhoiri/csms/internal/domain"
)

type TransactionService struct {
	transactionRepo domain.TransactionRepository
	chargePointRepo domain.ChargePointRepository
	idTagRepo       domain.IDTagRepository
	tariffConfig    config.TariffConfig
}

func NewTransactionService(
	transactionRepo domain.TransactionRepository,
	chargePointRepo domain.ChargePointRepository,
	idTagRepo domain.IDTagRepository,
	tariffConfig config.TariffConfig,
) domain.TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		chargePointRepo: chargePointRepo,
		idTagRepo:       idTagRepo,
		tariffConfig:    tariffConfig,
	}
}

func (s *TransactionService) StartTransaction(ctx context.Context, request *domain.StartTransactionRequest, chargePointID uint) (*domain.StartTransactionResponse, error) {
	activeTransaction, err := s.transactionRepo.GetActiveByConnector(ctx, chargePointID, request.ConnectorId)
	if err == nil && activeTransaction != nil {
		return nil, errors.New("connector is already in use")
	}

	idTag, err := s.idTagRepo.GetByTag(ctx, request.IDTag)
	if err != nil {
		return &domain.StartTransactionResponse{
			IDTagInfo: domain.IDTagInfo{
				Status: domain.AuthorizeStatusInvalid,
			},
			TransactionId: 0,
		}, nil
	}

	if idTag.Status != domain.AuthorizeStatusAccepted {
		return &domain.StartTransactionResponse{
			IDTagInfo: domain.IDTagInfo{
				Status: idTag.Status,
			},
			TransactionId: 0,
		}, nil
	}

	if !idTag.ExpiryDate.IsZero() && idTag.ExpiryDate.Before(time.Now()) {
		return &domain.StartTransactionResponse{
			IDTagInfo: domain.IDTagInfo{
				Status:     domain.AuthorizeStatusExpired,
				ExpiryDate: &idTag.ExpiryDate,
			},
			TransactionId: 0,
		}, nil
	}

	transaction := &domain.Transaction{
		ChargePointID:     chargePointID,
		ConnectorID:       request.ConnectorId,
		TransactionID:     int(time.Now().Unix()),
		IDTagID:           idTag.ID,
		StartMeterValue:   float64(request.MeterStart),
		CurrentMeterValue: float64(request.MeterStart),
		StartTime:         time.Now(),
		Status:            domain.TransactionStatusActive,
	}

	if err := s.transactionRepo.Create(ctx, transaction); err != nil {
		return nil, err
	}

	response := &domain.StartTransactionResponse{
		IDTagInfo: domain.IDTagInfo{
			Status: domain.AuthorizeStatusAccepted,
		},
		TransactionId: transaction.TransactionID,
	}

	return response, nil
}

func (s *TransactionService) StopTransaction(ctx context.Context, request *domain.StopTransactionRequest, chargePointID uint) (*domain.StopTransactionResponse, error) {
	transaction, err := s.transactionRepo.GetByTransactionID(ctx, request.TransactionId)
	if err != nil {
		return nil, errors.New("transaction not found")
	}

	if transaction.ChargePointID != chargePointID {
		return nil, errors.New("transaction does not belong to this charge point")
	}

	now := time.Now()
	transaction.StopMeterValue = float64(request.MeterStop)
	transaction.CurrentMeterValue = float64(request.MeterStop)
	transaction.StopTime = &now
	transaction.Status = domain.TransactionStatusCompleted
	transaction.Reason = request.Reason
	transaction.EnergyConsumed = transaction.StopMeterValue - transaction.StartMeterValue

	transaction.TotalCost = transaction.EnergyConsumed * s.tariffConfig.PricePerKwh

	if err := s.transactionRepo.Update(ctx, transaction); err != nil {
		return nil, err
	}

	response := &domain.StopTransactionResponse{
		Status: "Accepted",
	}

	return response, nil
}

func (s *TransactionService) GetTransaction(ctx context.Context, id uint) (*domain.Transaction, error) {
	return s.transactionRepo.GetByID(ctx, id)
}

func (s *TransactionService) ListTransactions(ctx context.Context, limit, offset int) ([]domain.Transaction, error) {
	return s.transactionRepo.List(ctx, limit, offset)
}

func (s *TransactionService) ListTransactionsByChargePoint(ctx context.Context, chargePointID uint) ([]domain.Transaction, error) {
	return s.transactionRepo.ListByChargePoint(ctx, chargePointID)
}

func (s *TransactionService) ListTransactionsByUser(ctx context.Context, idTag string) ([]domain.Transaction, error) {
	return s.transactionRepo.ListByUser(ctx, idTag)
}

func (s *TransactionService) UpdateMeterValues(ctx context.Context, request *domain.MeterValuesRequest, chargePointID uint) error {
	if request.TransactionId == nil {
		// No transaction ID provided, just log the meter values
		log.Printf("Meter values received for connector %d without transaction ID", request.ConnectorId)
		return nil
	}

	transaction, err := s.transactionRepo.GetByTransactionID(ctx, *request.TransactionId)
	if err != nil {
		log.Printf("Transaction not found for meter values update: %d", *request.TransactionId)
		return err
	}

	// Check if transaction belongs to this charge point
	if transaction.ChargePointID != chargePointID {
		log.Printf("Transaction %d does not belong to charge point %d", *request.TransactionId, chargePointID)
		return errors.New("transaction does not belong to this charge point")
	}

	if len(request.MeterValue) > 0 {
		latestMeterValue := request.MeterValue[len(request.MeterValue)-1] // Get the latest meter value

		if len(latestMeterValue.SampledValue) > 0 {
			// Find the energy meter value (usually measured in Wh or kWh)
			for _, sampledValue := range latestMeterValue.SampledValue {
				if sampledValue.Measurand == "Energy.Active.Import.Register" ||
					sampledValue.Measurand == "Energy.Active.Import.Interval" ||
					sampledValue.Measurand == "Power.Active.Import" {

					// Parse the meter value
					meterValue, err := parseMeterValue(sampledValue.Value, sampledValue.Unit)
					if err != nil {
						log.Printf("Error parsing meter value: %v", err)
						continue
					}

					// Update transaction with current meter value
					transaction.CurrentMeterValue = meterValue

					// Calculate energy consumed (convert to kWh if needed)
					energyConsumed := meterValue - transaction.StartMeterValue
					if sampledValue.Unit == "Wh" {
						energyConsumed = energyConsumed / 1000 // Convert Wh to kWh
					}

					transaction.EnergyConsumed = energyConsumed

					if err := s.transactionRepo.Update(ctx, transaction); err != nil {
						log.Printf("Error updating transaction with meter values: %v", err)
						return err
					}

					log.Printf("Updated transaction %d with meter value: %.2f %s, energy consumed: %.2f kWh",
						*request.TransactionId, meterValue, sampledValue.Unit, energyConsumed)
					break
				}
			}
		}
	}

	return nil
}

func parseMeterValue(value string, unit string) (float64, error) {
	value = strings.TrimSpace(value)

	meterValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid meter value format: %s", value)
	}

	return meterValue, nil
}
