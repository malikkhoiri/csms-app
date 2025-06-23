package service

import (
	"context"

	"github.com/malikkhoiri/csms/internal/domain"
)

// ConnectorService implements domain.ConnectorService
type ConnectorService struct {
	connectorRepo domain.ConnectorRepository
}

// NewConnectorService creates a new connector service
func NewConnectorService(connectorRepo domain.ConnectorRepository) domain.ConnectorService {
	return &ConnectorService{
		connectorRepo: connectorRepo,
	}
}

// UpdateConnectorStatus updates the status of a connector
func (s *ConnectorService) UpdateConnectorStatus(ctx context.Context, request *domain.StatusNotificationRequest, chargePointID uint) error {
	connector, err := s.connectorRepo.GetByChargePointAndConnectorID(ctx, chargePointID, request.ConnectorId)
	if err != nil {
		// Create new connector if it doesn't exist
		connector = &domain.Connector{
			ChargePointID:   chargePointID,
			ConnectorID:     request.ConnectorId,
			Status:          request.Status,
			ErrorCode:       request.ErrorCode,
			Info:            request.Info,
			VendorID:        request.VendorId,
			VendorErrorCode: request.VendorErrorCode,
		}
		return s.connectorRepo.Create(ctx, connector)
	}

	// Update existing connector
	connector.Status = request.Status
	connector.ErrorCode = request.ErrorCode
	connector.Info = request.Info
	connector.VendorID = request.VendorId
	connector.VendorErrorCode = request.VendorErrorCode

	return s.connectorRepo.Update(ctx, connector)
}

// GetConnector gets a connector by ID
func (s *ConnectorService) GetConnector(ctx context.Context, id uint) (*domain.Connector, error) {
	return s.connectorRepo.GetByID(ctx, id)
}

// ListConnectorsByChargePoint lists connectors by charge point ID
func (s *ConnectorService) ListConnectorsByChargePoint(ctx context.Context, chargePointID uint) ([]domain.Connector, error) {
	return s.connectorRepo.ListByChargePoint(ctx, chargePointID)
}

// GetConnectorByChargePointAndID gets a connector by charge point ID and connector ID
func (s *ConnectorService) GetConnectorByChargePointAndID(ctx context.Context, chargePointID uint, connectorID int) (*domain.Connector, error) {
	return s.connectorRepo.GetByChargePointAndConnectorID(ctx, chargePointID, connectorID)
}
