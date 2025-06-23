package domain

const (
	AuthorizeStatusAccepted     = "Accepted"
	AuthorizeStatusBlocked      = "Blocked"
	AuthorizeStatusExpired      = "Expired"
	AuthorizeStatusInvalid      = "Invalid"
	AuthorizeStatusConcurrentTx = "ConcurrentTx"
)

const (
	UserStatusActive   = "active"
	UserStatusInactive = "inactive"
	UserStatusBlocked  = "blocked"
)

const (
	ChargePointStatusAvailable   = "Available"
	ChargePointStatusOccupied    = "Occupied"
	ChargePointStatusFaulted     = "Faulted"
	ChargePointStatusUnavailable = "Unavailable"
	ChargePointStatusReserved    = "Reserved"
)

const (
	TransactionStatusActive    = "Active"
	TransactionStatusCompleted = "Completed"
	TransactionStatusCancelled = "Cancelled"
	TransactionStatusFailed    = "Failed"
	TransactionStatusPending   = "Pending"
)
