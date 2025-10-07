package grpc_handler

import (
	"context"
	"errors"
	"fmt"
	"strings"

	orgv1 "github.com/smallbiznis/go-genproto/smallbiznis/organization/v1"
	"github.com/smallbiznis/smallbiznis-apps/internal/organization/usecase"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// OrganizationHandler implements orgv1.OrganizationServiceServer
// This is a clean, repository-driven implementation scaffold ready for wiring with pgx/sqlc/GORM.
type OrganizationHandler struct {
	orgv1.UnimplementedOrganizationServiceServer
	countryUsecase usecase.ICountry
	orgUsecase     usecase.IOrganizationUsecase
}

func NewOrganization(
	countryUsecase usecase.ICountry,
	orgUsecase usecase.IOrganizationUsecase,
) *OrganizationHandler {
	return &OrganizationHandler{
		orgUsecase: orgUsecase,
	}
}

func (h *OrganizationHandler) ListCountry(ctx context.Context, req *orgv1.ListCountriesRequest) (*orgv1.ListCountriesResponse, error) {
	return h.countryUsecase.ListCountry(ctx, req)
}

func (h *OrganizationHandler) ListTimezone(ctx context.Context, req *orgv1.ListTimezoneRequest) (*orgv1.ListTimezoneResponse, error) {
	return h.orgUsecase.ListTimezone(ctx, req)
}

func (h *OrganizationHandler) ListCurrency(ctx context.Context, req *orgv1.ListCurrencyRequest) (*orgv1.ListCurrencyResponse, error) {
	return h.orgUsecase.ListCurrency(ctx, req)
}

func (h *OrganizationHandler) CreateOrganization(ctx context.Context, req *orgv1.CreateOrganizationRequest) (*orgv1.Organization, error) {

	if req.Type == orgv1.OrganizationType_TYPE_UNSPECIFIED {
		return nil, status.Error(codes.InvalidArgument, "type must be 'PERSONAL' or 'COMPANY'")
	}

	return h.orgUsecase.CreateOrganization(ctx, req)
}

func (h *OrganizationHandler) GetOrganization(ctx context.Context, req *orgv1.GetOrganizationRequest) (*orgv1.Organization, error) {
	return h.orgUsecase.GetOrganization(ctx, req)
}

func (h *OrganizationHandler) ListOrganization(ctx context.Context, req *orgv1.ListOrganizationRequest) (*orgv1.ListOrganizationResponse, error) {
	return h.orgUsecase.ListOrganization(ctx, req)
}

func (h *OrganizationHandler) UpdateOrganization(ctx context.Context, req *orgv1.UpdateOrganizationRequest) (*orgv1.Organization, error) {
	return h.orgUsecase.UpdateOrganization(ctx, req)
}

func (h *OrganizationHandler) CreateInvitation(ctx context.Context, req *orgv1.CreateInvitationRequest) (*orgv1.CreateInvitationResponse, error) {
	return h.orgUsecase.CreateInvitation(ctx, req)
}

func (h *OrganizationHandler) GetInvitation(ctx context.Context, req *orgv1.GetInvitationRequest) (*orgv1.Invitation, error) {
	return h.orgUsecase.GetInvitation(ctx, req)
}

func (h *OrganizationHandler) VerifyInvitation(ctx context.Context, req *orgv1.VerifyInvitationRequest) (*orgv1.VerifyInvitationResponse, error) {
	return h.orgUsecase.VerifyInvitation(ctx, req)
}

func (h *OrganizationHandler) CreateLocation(ctx context.Context, req *orgv1.CreateLocationRequest) (*orgv1.Location, error) {
	return h.orgUsecase.CreateLocation(ctx, req)
}

func (h *OrganizationHandler) GetLocation(ctx context.Context, req *orgv1.GetLocationRequest) (*orgv1.Location, error) {
	return h.orgUsecase.GetLocation(ctx, req)
}

func (h *OrganizationHandler) ListLocation(ctx context.Context, req *orgv1.ListLocationRequest) (*orgv1.ListLocationResponse, error) {
	return h.orgUsecase.ListLocation(ctx, req)
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
