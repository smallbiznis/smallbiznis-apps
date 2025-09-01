package domain

import (
	"crypto/rand"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type TransactionType string

var (
	EARNING      TransactionType = "earning"
	REDEEM       TransactionType = "redeem"
	EXPIRE       TransactionType = "expire"
	ADJUSTMENT   TransactionType = "adjustment"
	TRANSFER_IN  TransactionType = "transfer_in"
	TRANSFER_OUT TransactionType = "transfer_out"
)

func (t TransactionType) String() string {
	switch t {
	case EARNING, REDEEM, EXPIRE, ADJUSTMENT, TRANSFER_IN, TRANSFER_OUT:
		return string(t)
	default:
		return ""
	}
}

func (t TransactionType) IsValid() bool {
	switch t {
	case EARNING, REDEEM, EXPIRE, ADJUSTMENT, TRANSFER_IN, TRANSFER_OUT:
		return true
	default:
		return false
	}
}

type Transaction struct {
	ID                  string          `gorm:"column:id"`
	CreatedAt           time.Time       `gorm:"column:created_at"`
	UpdatedAt           time.Time       `gorm:"column:updated_at"`
	OrganizationID      string          `gorm:"column:organization_id"`
	UserID              string          `gorm:"column:user_id"`
	Type                TransactionType `gorm:"column:type"`
	Amount              int64           `gorm:"column:amount"`
	TransactionID       string          `gorm:"column:transaction_id"` // Generate by system
	ReferenceID         string          `gorm:"column:reference_id"`   // Generate by client
	Description         string          `gorm:"column:description"`
	Metadata            datatypes.JSON  `gorm:"column:metadata"`
	WorkflowID          string          `gorm:"column:workflow_id"`           // WorkflowID from Temporal
	LedgerTransactionID string          `gorm:"column:ledger_transaction_id"` // TransactionID from Ledger Service
	Status              string          `gorm:"column:status"`
}

type TransactionParam struct {
	OrganizationID string
	UserID         string
	Type           TransactionType
	Description    string
	ReferenceID    string
	TransactionID  string
	Metadata       datatypes.JSON
}

func NewTransaction(p TransactionParam) *Transaction {
	return &Transaction{
		ID:             uuid.NewString(),
		OrganizationID: p.OrganizationID,
		UserID:         p.UserID,
		Type:           p.Type,
		Description:    p.Description,
		ReferenceID:    p.ReferenceID,
		TransactionID:  p.TransactionID,
		Metadata:       p.Metadata,
	}
}

func GenerateTransactionID() (string, error) {
	datePart := time.Now().Format("20060102") // YYMMDD

	r := make([]byte, 3) // 3 bytes = 6 hex chars
	_, err := rand.Read(r)
	if err != nil {
		return "", err
	}
	randomPart := strings.ToUpper(fmt.Sprintf("%x", r))

	return fmt.Sprintf("TRX-%s-%s", datePart, randomPart), nil
}

type PointCredit struct {
	ID              string    `gorm:"column:id"`
	CreatedAt       time.Time `gorm:"column:created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at"`
	OrgID           string    `gorm:"column:org_id"`
	UserID          string    `gorm:"column:user_id"`
	Amount          int64     `gorm:"column:amount"`
	UnitAmount      int64     `gorm:"column:unit_amount"`
	RemainingAmount int64     `gorm:"column:remaining_amount"`
	ExpireAt        time.Time `gorm:"column:expire_at"`
}
