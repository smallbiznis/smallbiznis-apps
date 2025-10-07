package domain

import (
	"testing"
	"time"

	orgv1 "github.com/smallbiznis/go-genproto/smallbiznis/organization/v1"
)

func TestNewOrganization(t *testing.T) {
	org := NewOrganization()

	if org == nil {
		t.Fatal("expected organization instance")
	}

	if org.ID == "" {
		t.Fatal("expected organization ID to be set")
	}
}

func TestNewOrgCountry(t *testing.T) {
	country := NewOrgCountry("ID")

	if country.ID == "" {
		t.Fatal("expected country ID to be set")
	}

	if country.CountryCode != "ID" {
		t.Fatalf("expected country code to be ID, got %q", country.CountryCode)
	}
}

func TestNewOrgPlan(t *testing.T) {
	plan := NewOrgPlan("pro")

	if plan.ID == "" {
		t.Fatal("expected plan ID to be set")
	}

	if plan.PlanID != "pro" {
		t.Fatalf("expected plan ID to be pro, got %q", plan.PlanID)
	}
}

func TestNewInvitation(t *testing.T) {
	before := time.Now()
	inv := NewInvitation("org", "user@example.com", "admin")
	after := time.Now()

	if inv.ID == "" {
		t.Fatal("expected invitation ID to be set")
	}

	if inv.OrgID != "org" || inv.Email != "user@example.com" || inv.Role != "admin" {
		t.Fatalf("invitation fields not populated correctly: %+v", inv)
	}

	if inv.Status != orgv1.InvitationStatus_INVITATION_PENDING.String() {
		t.Fatalf("expected status %q, got %q", orgv1.InvitationStatus_INVITATION_PENDING, inv.Status)
	}

	lowerBound := before.Add(24 * time.Hour)
	upperBound := after.Add(24 * time.Hour)

	if inv.ExpiryAt.Before(lowerBound) || inv.ExpiryAt.After(upperBound) {
		t.Fatalf("expected expiry to be within 24h window, got %v (bounds %v - %v)", inv.ExpiryAt, lowerBound, upperBound)
	}

	if inv.AcceptAt != nil {
		t.Fatal("expected AcceptAt to be nil")
	}

	if inv.RevokeAt != nil {
		t.Fatal("expected RevokeAt to be nil")
	}

	if inv.Token != "" {
		t.Fatalf("expected token to be empty, got %q", inv.Token)
	}
}

func TestNewMember(t *testing.T) {
	member := NewMember("org", "user", "user@example.com", "role")

	if member.ID == "" {
		t.Fatal("expected member ID to be set")
	}

	if member.OrgID != "org" || member.UserID != "user" || member.Email != "user@example.com" || member.Role != "role" {
		t.Fatalf("member fields not populated correctly: %+v", member)
	}
}

func TestNewLocation(t *testing.T) {
	location := NewLocation("org")

	if location.ID == "" {
		t.Fatal("expected location ID to be set")
	}

	if location.OrgID != "org" {
		t.Fatalf("expected OrgID to be org, got %q", location.OrgID)
	}
}
