package domain

import (
	"time"

	"github.com/bwmarrin/snowflake"
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
	ID          snowflake.ID       `gorm:"column:id"`
	CreatedAt   time.Time          `gorm:"column:created_at"`
	UpdatedAt   time.Time          `gorm:"column:updated_at"`
	Type        string             `gorm:"column:type"`
	Name        string             `gorm:"column:name"`
	Slug        string             `gorm:"column:slug"`
	CountryCode string             `gorm:"column:country_code"`
	Status      OrganizationStatus `gorm:"foreignKey:OrgID"`
}

type OrganizationParams struct {
	ID          snowflake.ID
	Type        string
	Name        string
	Slug        string
	CountryCode string
	Status      OrganizationStatus
}

func NewOrganization(p OrganizationParams) *Organization {
	return &Organization{
		ID:          p.ID,
		Type:        p.Type,
		Name:        p.Name,
		Slug:        p.Slug,
		CountryCode: p.CountryCode,
		Status:      p.Status,
	}
}

type Status string

var (
	Active   Status = "ACTIVE"
	Inactive Status = "INACTIVE"
)

type OrganizationStatus struct {
	ID        string       `gorm:"column:id"`
	OrgID     snowflake.ID `gorm:"column:org_id"`
	Status    Status       `gorm:"column:status"`
	CreatedAt time.Time    `gorm:"column:created_at"`
}

type OrganizationStatusParams struct {
	OrgID  snowflake.ID
	Status Status
}

func NewStatus(p OrganizationStatusParams) *OrganizationStatus {
	return &OrganizationStatus{
		ID:     uuid.NewString(),
		OrgID:  p.OrgID,
		Status: p.Status,
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
	OrgID     snowflake.ID   `gorm:"column:org_id"`
	Email     string         `gorm:"column:email"`
	RoleID    string         `gorm:"column:role_id"`
	Status    string         `gorm:"column:status"`
	Token     string         `gorm:"column:token"`
	AcceptAt  *time.Time     `gorm:"column:accept_at"`
	RevokeAt  *time.Time     `gorm:"column:revoke_at"`
	ExpiryAt  time.Time      `gorm:"column:expiry_at"`
}

func NewInvitation(orgID snowflake.ID, email, role string) *Invitation {
	return &Invitation{
		ID:       uuid.NewString(),
		OrgID:    orgID,
		Email:    email,
		RoleID:   role,
		Status:   orgv1.InvitationStatus_INVITATION_PENDING.String(),
		ExpiryAt: time.Now().Add(24 * time.Hour),
	}
}

type Member struct {
	ID        string       `gorm:"column:id"`
	CreatedAt time.Time    `gorm:"column:created_at"`
	UpdatedAt time.Time    `gorm:"column:updated_at"`
	OrgID     snowflake.ID `gorm:"column:org_id"`
	UserID    string       `gorm:"column:user_id"`
	Email     string       `gorm:"column:email"`
	RoleID    string       `gorm:"column:role_id"`
}

type MemberParams struct {
	OrgID  snowflake.ID
	UserID string
	RoleID string
}

func NewMember(p MemberParams) *Member {
	return &Member{
		ID:     uuid.NewString(),
		OrgID:  p.OrgID,
		UserID: p.UserID,
		RoleID: p.RoleID,
	}
}

type Location struct {
	ID          string       `gorm:"column:id"`
	OrgID       snowflake.ID `gorm:"column:org_id"`
	Name        string       `gorm:"column:name"`
	Address     string       `gorm:"column:address"`
	City        string       `gorm:"column:city"`
	ZipCode     string       `gorm:"column:zip_code"`
	CountryCode string       `gorm:"column:country_code"`
	Timezone    string       `gorm:"column:timezone"`
}

type LocationParams struct {
	OrgID       snowflake.ID
	Name        string
	Address     string
	City        string
	ZipCode     string
	CountryCode string
	Timezone    string
}

func NewLocation(p LocationParams) *Location {
	return &Location{
		ID:          uuid.NewString(),
		OrgID:       p.OrgID,
		Name:        p.Name,
		Address:     p.Address,
		City:        p.City,
		ZipCode:     p.ZipCode,
		CountryCode: p.CountryCode,
		Timezone:    p.Timezone,
	}
}
