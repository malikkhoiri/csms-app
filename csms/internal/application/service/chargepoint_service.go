package service

import (
	"context"
	"time"

	"github.com/malikkhoiri/csms/internal/domain"
)

type ChargePointService struct {
	chargePointRepo domain.ChargePointRepository
	connectorRepo   domain.ConnectorRepository
}

func NewChargePointService(
	chargePointRepo domain.ChargePointRepository,
	connectorRepo domain.ConnectorRepository,
) domain.ChargePointService {
	return &ChargePointService{
		chargePointRepo: chargePointRepo,
		connectorRepo:   connectorRepo,
	}
}

func (s *ChargePointService) RegisterChargePoint(ctx context.Context, request *domain.BootNotificationRequest, chargePointCode string) (*domain.BootNotificationResponse, error) {
	existingCP, err := s.chargePointRepo.GetByCode(ctx, chargePointCode)
	if err != nil && err.Error() != "record not found" {
		return nil, err
	}

	now := time.Now()

	if existingCP != nil {
		existingCP.ChargeBoxSerialNumber = request.ChargeBoxSerialNumber
		existingCP.ChargePointVendor = request.ChargePointVendor
		existingCP.ChargePointModel = request.ChargePointModel
		existingCP.ChargePointSerialNumber = request.ChargePointSerialNumber
		existingCP.FirmwareVersion = request.FirmwareVersion
		existingCP.Iccid = request.Iccid
		existingCP.Imsi = request.Imsi
		existingCP.MeterType = request.MeterType
		existingCP.MeterSerialNumber = request.MeterSerialNumber
		existingCP.LastBootNotification = now
		existingCP.Status = "Available"

		if err := s.chargePointRepo.Update(ctx, existingCP); err != nil {
			return nil, err
		}
	} else {
		cp := &domain.ChargePoint{
			ChargePointCode:         chargePointCode,
			ChargeBoxSerialNumber:   request.ChargeBoxSerialNumber,
			ChargePointVendor:       request.ChargePointVendor,
			ChargePointModel:        request.ChargePointModel,
			ChargePointSerialNumber: request.ChargePointSerialNumber,
			FirmwareVersion:         request.FirmwareVersion,
			Iccid:                   request.Iccid,
			Imsi:                    request.Imsi,
			MeterType:               request.MeterType,
			MeterSerialNumber:       request.MeterSerialNumber,
			Status:                  "Available",
			LastBootNotification:    now,
			LastHeartbeat:           now,
		}

		if err := s.chargePointRepo.Create(ctx, cp); err != nil {
			return nil, err
		}

		connector := &domain.Connector{
			ChargePointID: cp.ID,
			ConnectorID:   1,
			Status:        "Available",
		}

		if err := s.connectorRepo.Create(ctx, connector); err != nil {
			return nil, err
		}
	}

	response := &domain.BootNotificationResponse{
		Status:      "Accepted",
		CurrentTime: now.UTC().Format(time.RFC3339),
		Interval:    300, // seconds
	}

	return response, nil
}

func (s *ChargePointService) UpdateChargePointStatus(ctx context.Context, chargePointID uint, status string) error {
	return s.chargePointRepo.UpdateStatus(ctx, chargePointID, status)
}

func (s *ChargePointService) GetChargePoint(ctx context.Context, id uint) (*domain.ChargePoint, error) {
	return s.chargePointRepo.GetByID(ctx, id)
}

func (s *ChargePointService) GetChargePointByCode(ctx context.Context, code string) (*domain.ChargePoint, error) {
	return s.chargePointRepo.GetByCode(ctx, code)
}

func (s *ChargePointService) ListChargePoints(ctx context.Context, limit, offset int) ([]domain.ChargePoint, error) {
	return s.chargePointRepo.List(ctx, limit, offset)
}

func (s *ChargePointService) UpdateHeartbeat(ctx context.Context, chargePointID uint) error {
	return s.chargePointRepo.UpdateHeartbeat(ctx, chargePointID)
}

func (s *ChargePointService) DeleteChargePoint(ctx context.Context, id uint) error {
	_, err := s.chargePointRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.chargePointRepo.Delete(ctx, id)
}
