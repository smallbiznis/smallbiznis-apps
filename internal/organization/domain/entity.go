package domain

import (
	"time"

	"github.com/google/uuid"
	orgv1 "github.com/smallbiznis/go-genproto/smallbiznis/organization/v1"
	"gorm.io/gorm"
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
	ID        string              `gorm:"column:id"`
	CreatedAt time.Time           `gorm:"column:created_at"`
	UpdatedAt time.Time           `gorm:"column:updated_at"`
	Type      string              `gorm:"column:type"`
	Name      string              `gorm:"column:name"`
	Slug      string              `gorm:"column:slug"`
	LogoURL   string              `gorm:"column:logo_url"`
	Country   OrganizationCountry `gorm:"foreignKey:OrgID"`
	Plan      OrganizationPlan    `gorm:"foreignKey:OrgID"`
}

func NewOrganization() *Organization {
	return &Organization{
		ID: uuid.NewString(),
	}
}

type OrganizationCountry struct {
	ID          string    `gorm:"column:id"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
	OrgID       string    `gorm:"column:org_id"`
	CountryCode string    `gorm:"column:country_code"`
}

func NewOrgCountry(code string) *OrganizationCountry {
	return &OrganizationCountry{
		ID:          uuid.NewString(),
		CountryCode: code,
	}
}

type OrganizationPlan struct {
	ID        string    `gorm:"column:id"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	OrgID     string    `gorm:"column:org_id"`
	PlanID    string    `gorm:"column:plan_id"` // free, pro, growth
}

func NewOrgPlan(plan string) *OrganizationPlan {
	return &OrganizationPlan{
		ID:     uuid.NewString(),
		PlanID: plan,
	}
}

type Invitation struct {
	ID        string         `gorm:"column:id"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
	OrgID     string         `gorm:"column:org_id"`
	Email     string         `gorm:"column:email"`
	Role      string         `gorm:"column:role"`
	Status    string         `gorm:"column:status"`
	Token     string         `gorm:"column:token"`
	AcceptAt  *time.Time     `gorm:"column:accept_at"`
	RevokeAt  *time.Time     `gorm:"column:revoke_at"`
	ExpiryAt  time.Time      `gorm:"column:expiry_at"`
}

func NewInvitation(orgID, email, role string) *Invitation {
	return &Invitation{
		ID:       uuid.NewString(),
		OrgID:    orgID,
		Email:    email,
		Role:     role,
		Status:   orgv1.InvitationStatus_INVITATION_PENDING.String(),
		ExpiryAt: time.Now().Add(24 * time.Hour),
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

type Location struct {
	ID          string `gorm:"column:id"`
	OrgID       string `gorm:"column:org_id"`
	Name        string `gorm:"column:name"`
	Address     string `gorm:"column:address"`
	City        string `gorm:"column:city"`
	ZipCode     string `gorm:"column:zip_code"`
	CountryCode string `gorm:"column:country_code"`
	Timezone    string `gorm:"column:timezone"`
}

func NewLocation(orgId string) *Location {
	return &Location{
		ID:    uuid.NewString(),
		OrgID: orgId,
	}
}
