package domain

import (
	"time"

	"github.com/google/uuid"
)

type Country struct {
	Code string `gorm:"column:code"`
	Name string `gorm:"column:name"`
}

type Timezone struct {
	ID          string `gorm:"column:id"`
	CountryCode string `gorm:"country_code"`
	Tz          string `gorm:"column:tz"`
}

type Currency struct {
	ID          string `gorm:"column:id"`
	CountryCode string `gorm:"column:country_code"`
	Format      string `gorm:"column:format"`
	Currency    string `gorm:"column:currency"`
}

type Organization struct {
	ID          string    `gorm:"column:id"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
	LogoURL     string    `gorm:"column:logo_url"`
	Slug        string    `gorm:"column:slug"`
	Name        string    `gorm:"column:name"`
	CountryCode string    `gorm:"column:country_code"`
}

func NewOrganization() *Organization {
	return &Organization{
		ID: uuid.NewString(),
	}
}

type Invitation struct {
	ID        string    `gorm:"column:id"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	OrgID     string    `gorm:"column:org_id"`
	Email     string    `gorm:"column:email"`
	Role      string    `gorm:"column:role"`
}

func NewInvitation(orgID string) *Invitation {
	return &Invitation{
		ID:    uuid.NewString(),
		OrgID: orgID,
	}
}

type Member struct {
	ID        string    `gorm:"column:id"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	OrgID     string    `gorm:"column:org_id"`
	UserID    string    `gorm:"column:user_id"`
	Email     string    `gorm:"column:email"`
	Role      string    `gorm:"column:role"`
}

func NewMember(OrgID, UserID, Email, Role string) *Member {
	return &Member{
		ID:     uuid.NewString(),
		OrgID:  OrgID,
		UserID: UserID,
		Email:  Email,
		Role:   Role,
	}
}
