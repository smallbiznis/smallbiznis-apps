package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"testing"
	"time"

	"gorm.io/datatypes"
)

func TestNewLedgerEntryAssignsFields(t *testing.T) {
	params := LedgerParams{
		OrgID:         "org-123",
		UserID:        "user-456",
		Type:          "credit",
		Amount:        1000,
		ReferenceID:   "ref-789",
		TransactionID: "txn-101",
		Description:   "sample-description",
		PreviousHash:  "prev-hash",
		Metadata:      datatypes.JSON([]byte(`{"key":"value"}`)),
	}

	entry := NewLedgerEntry(params)

	if entry.ID == "" {
		t.Fatal("expected ID to be generated")
	}

	if entry.OrgID != params.OrgID {
		t.Fatalf("expected OrgID %q, got %q", params.OrgID, entry.OrgID)
	}

	if entry.Metadata == nil {
		t.Fatal("expected Metadata to be assigned")
	}
}

func TestLedgerEntryHashFields(t *testing.T) {
	createdAt := time.Date(2025, 1, 2, 3, 4, 5, 600000000, time.UTC)
	entry := &LedgerEntry{
		ID:            "entry-id",
		OrgID:         "org-id",
		UserID:        "user-id",
		Type:          "debit",
		Amount:        2500,
		TransactionID: "txn-id",
		ReferenceID:   "ref-id",
		Description:   "hash fields",
		PreviousHash:  "prev",
		CreatedAt:     createdAt,
	}

	fields := entry.HashFields()

	expected := map[string]string{
		"id":             entry.ID,
		"org_id":         entry.OrgID,
		"user_id":        entry.UserID,
		"type":           entry.Type,
		"amount":         fmt.Sprintf("%d", entry.Amount),
		"transaction_id": entry.TransactionID,
		"reference_id":   entry.ReferenceID,
		"description":    entry.Description,
		"created_at":     createdAt.UTC().Format(time.RFC3339Nano),
		"previous_hash":  entry.PreviousHash,
	}

	if len(fields) != len(expected) {
		t.Fatalf("expected %d fields, got %d", len(expected), len(fields))
	}

	for key, want := range expected {
		if got := fields[key]; got != want {
			t.Fatalf("expected field %q to be %q, got %q", key, want, got)
		}
	}
}

func TestLedgerEntryGenerateHash(t *testing.T) {
	createdAt := time.Date(2025, 1, 2, 3, 4, 5, 600000000, time.UTC)
	entry := &LedgerEntry{
		ID:            "entry-id",
		OrgID:         "org-id",
		UserID:        "user-id",
		Type:          "credit",
		Amount:        1000,
		TransactionID: "txn-id",
		ReferenceID:   "ref-id",
		Description:   "generate hash",
		PreviousHash:  "prev",
		CreatedAt:     createdAt,
	}

	got := entry.GenerateHash()

	fields := entry.HashFields()
	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, fields[k]))
	}

	joined := strings.Join(parts, "|")
	expectedBytes := sha256.Sum256([]byte(joined))
	expected := hex.EncodeToString(expectedBytes[:])

	if got != expected {
		t.Fatalf("expected hash %q, got %q", expected, got)
	}
}

func TestGenerateTransactionIDFormat(t *testing.T) {
	id, err := GenerateTransactionID()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	pattern := regexp.MustCompile(`^\d{8}-[0-9A-F]{6}$`)
	if !pattern.MatchString(id) {
		t.Fatalf("transaction ID %q does not match expected pattern", id)
	}

	parts := strings.Split(id, "-")
	if len(parts) != 2 {
		t.Fatalf("expected transaction ID to have two parts, got %d", len(parts))
	}

	if _, err := time.Parse("20060102", parts[0]); err != nil {
		t.Fatalf("date part %q is not a valid date: %v", parts[0], err)
	}
}
