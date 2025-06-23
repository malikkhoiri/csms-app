package domain

import (
	"context"
)

type ChargePointService interface {
	RegisterChargePoint(ctx context.Context, request *BootNotificationRequest, chargePointCode string) (*BootNotificationResponse, error)
	UpdateChargePointStatus(ctx context.Context, chargePointID uint, status string) error
	GetChargePoint(ctx context.Context, id uint) (*ChargePoint, error)
	GetChargePointByCode(ctx context.Context, code string) (*ChargePoint, error)
	ListChargePoints(ctx context.Context, limit, offset int) ([]ChargePoint, error)
	UpdateHeartbeat(ctx context.Context, chargePointID uint) error
	DeleteChargePoint(ctx context.Context, id uint) error
}

type TransactionService interface {
	StartTransaction(ctx context.Context, request *StartTransactionRequest, chargePointID uint) (*StartTransactionResponse, error)
	StopTransaction(ctx context.Context, request *StopTransactionRequest, chargePointID uint) (*StopTransactionResponse, error)
	GetTransaction(ctx context.Context, id uint) (*Transaction, error)
	ListTransactions(ctx context.Context, limit, offset int) ([]Transaction, error)
	ListTransactionsByChargePoint(ctx context.Context, chargePointID uint) ([]Transaction, error)
	ListTransactionsByUser(ctx context.Context, idTag string) ([]Transaction, error)
	UpdateMeterValues(ctx context.Context, request *MeterValuesRequest, chargePointID uint) error
}

type UserService interface {
	CreateUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, id uint) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id uint) error
	ListUsers(ctx context.Context, limit, offset int) ([]User, error)
	UpdateUserStatus(ctx context.Context, id uint, status string) error
}

type AuthService interface {
	Login(ctx context.Context, request *AuthRequest) (*AuthResponse, error)
	ValidateToken(ctx context.Context, token string) (*User, error)
	Logout(ctx context.Context, token string) error
}

type ConnectorService interface {
	UpdateConnectorStatus(ctx context.Context, request *StatusNotificationRequest, chargePointID uint) error
	GetConnector(ctx context.Context, id uint) (*Connector, error)
	ListConnectorsByChargePoint(ctx context.Context, chargePointID uint) ([]Connector, error)
	GetConnectorByChargePointAndID(ctx context.Context, chargePointID uint, connectorID int) (*Connector, error)
}

type NotificationService interface {
	SendTransactionNotification(ctx context.Context, transaction *Transaction) error
	SendErrorNotification(ctx context.Context, chargePointID uint, error string) error
	SendStatusNotification(ctx context.Context, chargePointID uint, status string) error
}

type MonitoringService interface {
	GetSystemStatus(ctx context.Context) (map[string]interface{}, error)
	GetChargePointMetrics(ctx context.Context, chargePointID uint) (map[string]interface{}, error)
	GetTransactionMetrics(ctx context.Context) (map[string]interface{}, error)
	GetUserMetrics(ctx context.Context) (map[string]interface{}, error)
}

type IDTagService interface {
	Authorize(ctx context.Context, request *AuthorizeRequest) (*AuthorizeResponse, error)
	CreateIDTag(ctx context.Context, idTag *IDTag) error
	GetIDTag(ctx context.Context, id uint) (*IDTag, error)
	GetByTag(ctx context.Context, tag string) (*IDTag, error)
	UpdateIDTag(ctx context.Context, idTag *IDTag) error
	DeleteIDTag(ctx context.Context, id uint) error
	ListIDTags(ctx context.Context, limit, offset int) ([]IDTag, error)
	ListByUser(ctx context.Context, userID uint) ([]IDTag, error)
}
