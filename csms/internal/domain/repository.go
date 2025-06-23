package domain

import (
	"context"
)

type ChargePointRepository interface {
	Create(ctx context.Context, cp *ChargePoint) error
	GetByID(ctx context.Context, id uint) (*ChargePoint, error)
	GetByCode(ctx context.Context, code string) (*ChargePoint, error)
	Update(ctx context.Context, cp *ChargePoint) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, limit, offset int) ([]ChargePoint, error)
	UpdateStatus(ctx context.Context, id uint, status string) error
	UpdateHeartbeat(ctx context.Context, id uint) error
}

type ConnectorRepository interface {
	Create(ctx context.Context, connector *Connector) error
	GetByID(ctx context.Context, id uint) (*Connector, error)
	GetByChargePointAndConnectorID(ctx context.Context, chargePointID uint, connectorID int) (*Connector, error)
	Update(ctx context.Context, connector *Connector) error
	Delete(ctx context.Context, id uint) error
	ListByChargePoint(ctx context.Context, chargePointID uint) ([]Connector, error)
	UpdateStatus(ctx context.Context, id uint, status string) error
}

type TransactionRepository interface {
	Create(ctx context.Context, transaction *Transaction) error
	GetByID(ctx context.Context, id uint) (*Transaction, error)
	GetByTransactionID(ctx context.Context, transactionID int) (*Transaction, error)
	Update(ctx context.Context, transaction *Transaction) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, limit, offset int) ([]Transaction, error)
	ListByChargePoint(ctx context.Context, chargePointID uint) ([]Transaction, error)
	ListByUser(ctx context.Context, idTag string) ([]Transaction, error)
	GetActiveByConnector(ctx context.Context, chargePointID uint, connectorID int) (*Transaction, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uint) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, limit, offset int) ([]User, error)
	UpdateStatus(ctx context.Context, id uint, status string) error
}

type OCPPMessageRepository interface {
	Create(ctx context.Context, message *OCPPMessage) error
	GetByID(ctx context.Context, id uint) (*OCPPMessage, error)
	ListByChargePoint(ctx context.Context, chargePointID uint, limit, offset int) ([]OCPPMessage, error)
	ListByAction(ctx context.Context, action string, limit, offset int) ([]OCPPMessage, error)
}

type IDTagRepository interface {
	Create(ctx context.Context, idTag *IDTag) error
	GetByID(ctx context.Context, id uint) (*IDTag, error)
	GetByTag(ctx context.Context, tag string) (*IDTag, error)
	Update(ctx context.Context, idTag *IDTag) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, limit, offset int) ([]IDTag, error)
	ListByUser(ctx context.Context, userID uint) ([]IDTag, error)
}
