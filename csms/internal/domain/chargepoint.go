package domain

import (
	"time"
)

type ChargePoint struct {
	ID                      uint      `json:"id" gorm:"primaryKey"`
	ChargePointCode         string    `json:"chargePointCode" gorm:"uniqueIndex;not null"`
	ChargeBoxSerialNumber   string    `json:"chargeBoxSerialNumber"`
	ChargePointModel        string    `json:"chargePointModel" gorm:"not null"`
	ChargePointVendor       string    `json:"chargePointVendor"`
	ChargePointSerialNumber string    `json:"chargePointSerialNumber"`
	FirmwareVersion         string    `json:"firmwareVersion"`
	Iccid                   string    `json:"iccid"`
	Imsi                    string    `json:"imsi"`
	MeterType               string    `json:"meterType"`
	MeterSerialNumber       string    `json:"meterSerialNumber"`
	Status                  string    `json:"status" gorm:"default:'Available'"`
	LastHeartbeat           time.Time `json:"lastHeartbeat"`
	LastBootNotification    time.Time `json:"lastBootNotification"`
	CreatedAt               time.Time `json:"createdAt"`
	UpdatedAt               time.Time `json:"updatedAt"`

	Connectors   []Connector   `json:"connectors" gorm:"foreignKey:ChargePointID"`
	Transactions []Transaction `json:"transactions" gorm:"foreignKey:ChargePointID"`
}

type Connector struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	ChargePointID   uint      `json:"chargePointId" gorm:"not null"`
	ConnectorID     int       `json:"connectorId" gorm:"not null"`
	Status          string    `json:"status" gorm:"default:'Available'"`
	ErrorCode       string    `json:"errorCode"`
	Info            string    `json:"info"`
	VendorID        string    `json:"vendorId"`
	VendorErrorCode string    `json:"vendorErrorCode"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`

	ChargePoint ChargePoint `json:"chargePoint" gorm:"foreignKey:ChargePointID"`
}

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	Email     string    `json:"email" gorm:"uniqueIndex"`
	Password  string    `json:"-" gorm:"not null"`
	Phone     string    `json:"phone"`
	Role      string    `json:"role" gorm:"default:'customer'"`
	Status    string    `json:"status" gorm:"default:'active'"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	IDTags []IDTag `json:"idTags" gorm:"foreignKey:UserID"`
}

type IDTag struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Tag        string    `json:"tag" gorm:"uniqueIndex;not null"`
	Status     string    `json:"status" gorm:"default:'Accepted'"`
	ExpiryDate time.Time `json:"expiryDate"`
	UserID     uint      `json:"userId" gorm:"not null"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`

	User         User          `json:"user" gorm:"foreignKey:UserID"`
	Transactions []Transaction `json:"transactions" gorm:"foreignKey:IDTagID"`
}

type Transaction struct {
	ID                uint       `json:"id" gorm:"primaryKey"`
	ChargePointID     uint       `json:"chargePointId" gorm:"not null"`
	ConnectorID       int        `json:"connectorId" gorm:"not null"`
	TransactionID     int        `json:"transactionId" gorm:"not null"`
	IDTagID           uint       `json:"idTagId" gorm:"not null"`
	StartMeterValue   float64    `json:"startMeterValue"`
	StopMeterValue    float64    `json:"stopMeterValue"`
	CurrentMeterValue float64    `json:"currentMeterValue"`
	EnergyConsumed    float64    `json:"energyConsumed"`
	TotalCost         float64    `json:"totalCost"`
	StartTime         time.Time  `json:"startTime"`
	StopTime          *time.Time `json:"stopTime"`
	Status            string     `json:"status" gorm:"default:'Active'"`
	Reason            string     `json:"reason"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`

	IDTag       IDTag       `json:"idTag" gorm:"foreignKey:IDTagID"`
	ChargePoint ChargePoint `json:"chargePoint" gorm:"foreignKey:ChargePointID"`
}

type OCPPMessage struct {
	MessageType   int       `json:"messageType"`
	MessageID     string    `json:"messageId"`
	Action        string    `json:"action"`
	Payload       string    `json:"payload" gorm:"type:jsonb"`
	ChargePointID uint      `json:"chargePointId"`
	Timestamp     time.Time `json:"timestamp"`
}

type BootNotificationRequest struct {
	ChargePointVendor       string `json:"chargePointVendor"`
	ChargePointModel        string `json:"chargePointModel"`
	ChargePointSerialNumber string `json:"chargePointSerialNumber"`
	ChargeBoxSerialNumber   string `json:"chargeBoxSerialNumber"`
	FirmwareVersion         string `json:"firmwareVersion"`
	Iccid                   string `json:"iccid"`
	Imsi                    string `json:"imsi"`
	MeterType               string `json:"meterType"`
	MeterSerialNumber       string `json:"meterSerialNumber"`
}

type BootNotificationResponse struct {
	Status      string `json:"status"`
	CurrentTime string `json:"currentTime"`
	Interval    int    `json:"interval"`
}

type AuthorizeRequest struct {
	IDTag string `json:"idTag"`
}

type AuthorizeResponse struct {
	IDTagInfo IDTagInfo `json:"idTagInfo"`
}

type IDTagInfo struct {
	Status      string     `json:"status"`
	ExpiryDate  *time.Time `json:"expiryDate,omitempty"`
	ParentIDTag string     `json:"parentIdTag,omitempty"`
}

type StartTransactionRequest struct {
	ConnectorId   int    `json:"connectorId"`
	IDTag         string `json:"idTag"`
	MeterStart    int    `json:"meterStart"`
	ReservationId *int   `json:"reservationId,omitempty"`
}

type StartTransactionResponse struct {
	IDTagInfo     IDTagInfo `json:"idTagInfo"`
	TransactionId int       `json:"transactionId"`
}

type StopTransactionRequest struct {
	TransactionId int    `json:"transactionId"`
	IDTag         string `json:"idTag"`
	MeterStop     int    `json:"meterStop"`
	Reason        string `json:"reason"`
}

type StopTransactionResponse struct {
	Status string `json:"status"`
}

type StatusNotificationRequest struct {
	ConnectorId     int    `json:"connectorId"`
	Status          string `json:"status"`
	ErrorCode       string `json:"errorCode"`
	Info            string `json:"info"`
	VendorId        string `json:"vendorId"`
	VendorErrorCode string `json:"vendorErrorCode"`
}

type MeterValuesRequest struct {
	ConnectorId   int          `json:"connectorId"`
	TransactionId *int         `json:"transactionId,omitempty"`
	MeterValue    []MeterValue `json:"meterValue"`
}

type MeterValue struct {
	Timestamp    string         `json:"timestamp"`
	SampledValue []SampledValue `json:"sampledValue"`
}

type SampledValue struct {
	Value     string `json:"value"`
	Context   string `json:"context,omitempty"`
	Format    string `json:"format,omitempty"`
	Measurand string `json:"measurand,omitempty"`
	Phase     string `json:"phase,omitempty"`
	Location  string `json:"location,omitempty"`
	Unit      string `json:"unit,omitempty"`
}

type AuthRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	AccessToken string `json:"accessToken"`
	User        User   `json:"user"`
}
