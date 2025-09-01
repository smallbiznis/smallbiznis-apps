package grpc

import (
	"context"
	"errors"
	"fmt"
	"strings"

	orgv1 "github.com/smallbiznis/go-genproto/smallbiznis/organization/v1"
	"github.com/smallbiznis/smallbiznis-apps/internal/organization/domain"
	"github.com/smallbiznis/smallbiznis-apps/internal/organization/usecase"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// OrganizationHandler implements orgv1.OrganizationServiceServer
// This is a clean, repository-driven implementation scaffold ready for wiring with pgx/sqlc/GORM.
type OrganizationHandler struct {
	orgv1.UnimplementedOrganizationServiceServer
	usecase usecase.IOrganizationUsecase
}

func NewOrganizationHandler(
	usecase usecase.IOrganizationUsecase,
) *OrganizationHandler {
	return &OrganizationHandler{
		usecase: usecase,
	}
}

func (h *OrganizationHandler) ListCountry(ctx context.Context) {}

func (h *OrganizationHandler) ListTimezone(ctx context.Context) {}

func (h *OrganizationHandler) ListCurrency(ctx context.Context) {}

func (h *OrganizationHandler) CreateOrganization(ctx context.Context, req *orgv1.CreateOrganizationRequest) (*orgv1.Organization, error) {
	return h.usecase.CreateOrganization(ctx, req)
}

func (h *OrganizationHandler) GetOrganization(ctx context.Context, req *orgv1.GetOrganizationRequest) (*orgv1.Organization, error) {
	return h.usecase.GetOrganization(ctx, req)
}

func (h *OrganizationHandler) ListOrganization(ctx context.Context, req *orgv1.ListOrganizationRequest) (*orgv1.ListOrganizationResponse, error) {
	return h.usecase.ListOrganization(ctx, req)
}

func (h *OrganizationHandler) UpdateOrganization(ctx context.Context, req *orgv1.UpdateOrganizationRequest) (*orgv1.Organization, error) {
	return h.usecase.UpdateOrganization(ctx, req)
}

func toProtoOrg(o domain.Organization) *orgv1.Organization {
	var trial *timestamppb.Timestamp
	return &orgv1.Organization{
		Id:          o.ID,
		Slug:        o.Slug,
		Name:        o.Name,
		LogoUrl:     o.LogoURL,
		TrialEndsAt: trial,
		CreatedAt:   timestamppb.New(o.CreatedAt),
		UpdatedAt:   timestamppb.New(o.UpdatedAt),
	}
}

func slugify(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	s = strings.ReplaceAll(s, " ", "-")

	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			b.WriteRune(r)
		}
	}
	out := b.String()
	if out == "" {
		out = uuid.NewString()[:8]
	}
	return out
}

var (
	errNotFound      = errors.New("not_found")
	errAlreadyExists = errors.New("already_exists")
)

func toStatus(err error) error {
	switch {
	case errors.Is(err, errNotFound):
		return status.Error(codes.NotFound, "not found")
	case errors.Is(err, errAlreadyExists):
		return status.Error(codes.AlreadyExists, "already exists")
	default:
		return status.Errorf(codes.Internal, "internal error: %v", err)
	}
}

// readIdempotencyKey demonstrates how to access the Idempotency-Key header (if you need it in repo layer)
func readIdempotencyKey(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		vals := md.Get("idempotency-key")
		if len(vals) > 0 {
			return vals[0]
		}
	}
	return ""
}

// Example: enforce simple server-side validation utilities
func require(cond bool, code codes.Code, format string, a ...any) error {
	if cond {
		return nil
	}
	return status.Error(code, fmt.Sprintf(format, a...))
}
