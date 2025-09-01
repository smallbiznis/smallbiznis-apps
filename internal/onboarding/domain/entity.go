package domain

import "time"

type OnboardingType string

var (
	Organization OnboardingType = "org"
)

type Onboarding struct {
	ID         string         `gorm:"column:id"`
	CreatedAt  time.Time      `gorm:"column:created_at"`
	UpdatedAt  time.Time      `gorm:"column:updated_at"`
	Type       OnboardingType `gorm:"column:type"`
	IsTemplate bool           `gorm:"column:is_template"`
	Name       string         `gorm:"column:name"`
}

type OnboardingState struct {
	ID           string    `gorm:"column:id"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
	OnboardingID string    `gorm:"column:onboarding_id"`
}
