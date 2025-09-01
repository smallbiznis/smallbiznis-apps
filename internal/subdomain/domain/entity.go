package domain

import "time"

type Subdomain struct {
	ID             string    `gorm:"column:id"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
	OrganizationID string    `gorm:"column:organization_id"`
	Name           string    `gorm:"column:name"`
}
