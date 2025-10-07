package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	orgv1 "github.com/smallbiznis/go-genproto/smallbiznis/organization/v1"
	"github.com/smallbiznis/smallbiznis-apps/internal/organization/domain"
	"github.com/smallbiznis/smallbiznis-apps/pkg/db/option"
	"github.com/smallbiznis/smallbiznis-apps/pkg/db/pagination"
	"github.com/smallbiznis/smallbiznis-apps/pkg/gen"
	"github.com/smallbiznis/smallbiznis-apps/pkg/repository"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

//go:generate mockgen -source=organization_usecase.go -destination=./../../usecase/mock_organization_usecase.go -package=usecase
type IOrganizationUsecase interface {
	ListTimezone(context.Context, *orgv1.ListTimezoneRequest) (*orgv1.ListTimezoneResponse, error)
	ListCurrency(context.Context, *orgv1.ListCurrencyRequest) (*orgv1.ListCurrencyResponse, error)

	CreateOrganization(context.Context, *orgv1.CreateOrganizationRequest) (*orgv1.Organization, error)
	GetOrganization(context.Context, *orgv1.GetOrganizationRequest) (*orgv1.Organization, error)
	ListOrganization(context.Context, *orgv1.ListOrganizationRequest) (*orgv1.ListOrganizationResponse, error)
	UpdateOrganization(context.Context, *orgv1.UpdateOrganizationRequest) (*orgv1.Organization, error)

	CreateLocation(context.Context, *orgv1.CreateLocationRequest) (*orgv1.Location, error)
	ListLocation(context.Context, *orgv1.ListLocationRequest) (*orgv1.ListLocationResponse, error)
	GetLocation(context.Context, *orgv1.GetLocationRequest) (*orgv1.Location, error)

	CreateInvitation(context.Context, *orgv1.CreateInvitationRequest) (*orgv1.CreateInvitationResponse, error)
	GetInvitation(context.Context, *orgv1.GetInvitationRequest) (*orgv1.Invitation, error)
	VerifyInvitation(context.Context, *orgv1.VerifyInvitationRequest) (*orgv1.VerifyInvitationResponse, error)
}

type organizationUsecase struct {
	db             *gorm.DB
	snowflake      *gen.SnowflakeNode
	countryRepo    repository.Repository[domain.Country]
	timezoneRepo   repository.Repository[domain.Timezone]
	currencyRepo   repository.Repository[domain.Currency]
	orgRepo        repository.Repository[domain.Organization]
	memberRepo     repository.Repository[domain.Member]
	locationRepo   repository.Repository[domain.Location]
	invitationRepo repository.Repository[domain.Invitation]
}

type OrganizationParams struct {
	fx.In
	DB        *gorm.DB
	Snowflake *gen.SnowflakeNode
}

func NewOrganization(
	p OrganizationParams,
) IOrganizationUsecase {
	return &organizationUsecase{
		db:             p.DB,
		snowflake:      p.Snowflake,
		countryRepo:    repository.ProvideStore[domain.Country](p.DB),
		timezoneRepo:   repository.ProvideStore[domain.Timezone](p.DB),
		currencyRepo:   repository.ProvideStore[domain.Currency](p.DB),
		orgRepo:        repository.ProvideStore[domain.Organization](p.DB),
		memberRepo:     repository.ProvideStore[domain.Member](p.DB),
		locationRepo:   repository.ProvideStore[domain.Location](p.DB),
		invitationRepo: repository.ProvideStore[domain.Invitation](p.DB),
	}
}

func (uc *organizationUsecase) ListCountry(ctx context.Context, req *orgv1.ListCountriesRequest) (*orgv1.ListCountriesResponse, error) {

	countries, err := uc.countryRepo.Find(ctx, &domain.Country{
		Code: req.Code,
	})
	if err != nil {
		return nil, err
	}

	var newCountries []*orgv1.Countries
	for _, c := range countries {
		newCountries = append(newCountries, &orgv1.Countries{
			Code: c.Code,
			Name: c.Name,
		})
	}

	return &orgv1.ListCountriesResponse{
		Data: newCountries,
	}, nil

}

func (uc *organizationUsecase) ListTimezone(ctx context.Context, req *orgv1.ListTimezoneRequest) (*orgv1.ListTimezoneResponse, error) {

	timezones, err := uc.timezoneRepo.Find(ctx, &domain.Timezone{
		CountryCode: req.CountryCode,
	})
	if err != nil {
		return nil, err
	}

	var newTimezones []*orgv1.Timezones
	for _, c := range timezones {
		newTimezones = append(newTimezones, &orgv1.Timezones{
			Id:          c.ID,
			CountryCode: c.CountryCode,
			Tz:          c.Tz,
		})
	}

	return &orgv1.ListTimezoneResponse{
		Data: newTimezones,
	}, nil

}

func (uc *organizationUsecase) ListCurrency(ctx context.Context, req *orgv1.ListCurrencyRequest) (*orgv1.ListCurrencyResponse, error) {

	currencies, err := uc.currencyRepo.Find(ctx, &domain.Currency{
		CountryCode: req.CountryCode,
	})
	if err != nil {
		return nil, err
	}

	var newCurrencies []*orgv1.Currencies
	for _, c := range currencies {
		newCurrencies = append(newCurrencies, &orgv1.Currencies{
			Id:          c.ID,
			CountryCode: c.CountryCode,
			Format:      c.Format,
			Currency:    c.Currency,
		})
	}

	return &orgv1.ListCurrencyResponse{
		Data: newCurrencies,
	}, nil

}

func (uc *organizationUsecase) LookupOrganization(ctx context.Context, req *orgv1.LookupOrganizationRequest) (*orgv1.LookupOrganizationResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	traceID := span.SpanContext().TraceID().String()
	spanID := span.SpanContext().SpanID().String()

	fields := []zap.Field{
		zap.String("trace_id", traceID),
		zap.String("span_id", spanID),
	}

	name := slug.Make(req.Name)
	exist, err := uc.orgRepo.FindOne(ctx, &domain.Organization{
		Slug: name,
	})
	if err != nil {
		zap.L().With(fields...).Error("failed lookup organization", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	if exist != nil {
		return nil, status.Error(codes.InvalidArgument, "organization not found")
	}

	return &orgv1.LookupOrganizationResponse{
		Organization: name,
	}, nil
}

func (uc *organizationUsecase) CreateOrganization(ctx context.Context, req *orgv1.CreateOrganizationRequest) (*orgv1.Organization, error) {

	span := trace.SpanFromContext(ctx)
	defer span.End()

	traceID := span.SpanContext().TraceID().String()
	spanID := span.SpanContext().SpanID().String()

	fields := []zap.Field{
		zap.String("trace_id", traceID),
		zap.String("span_id", spanID),
	}

	country, err := uc.countryRepo.FindOne(ctx, &domain.Country{Code: req.CountryCode})
	if err != nil {
		zap.L().With(fields...).Error("failed get country by code", zap.Error(err))
		return nil, err
	}

	if country == nil {
		return nil, fmt.Errorf("invalid countryCode")
	}

	unqName := slug.Make(
		strings.Trim(req.Name, " "),
	)

	exist, err := uc.orgRepo.FindOne(ctx, &domain.Organization{Slug: unqName})
	if err != nil {
		return nil, err
	}

	if exist != nil {
		return nil, fmt.Errorf("organization already exists")
	}

	orgID := uc.snowflake.GenerateID()
	newOrg := domain.NewOrganization(
		domain.OrganizationParams{
			ID:          orgID,
			Type:        req.Type.String(),
			Name:        req.Name,
			Slug:        unqName,
			CountryCode: country.Code,
			Status: domain.OrganizationStatus{
				ID:     uuid.NewString(),
				OrgID:  orgID,
				Status: domain.Active,
			},
		},
	)

	if err := uc.orgRepo.Create(ctx, newOrg); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return uc.GetOrganization(ctx, &orgv1.GetOrganizationRequest{OrgId: orgID.String()})
}

func (uc *organizationUsecase) GetOrganization(ctx context.Context, req *orgv1.GetOrganizationRequest) (*orgv1.Organization, error) {

	query := domain.Organization{Slug: req.OrgId}
	if orgID, err := snowflake.ParseString(req.OrgId); err == nil {
		query = domain.Organization{
			ID: orgID,
		}
	}

	opts := []option.QueryOption{
		option.WithPreloads("Status", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}),
	}

	exist, err := uc.orgRepo.FindOne(ctx, &query, opts...)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if exist == nil {
		return nil, status.Error(codes.NotFound, "organization not found")
	}

	return &orgv1.Organization{
		OrgId:     exist.ID.String(),
		Type:      orgv1.OrganizationType(orgv1.OrganizationType_value[exist.Type]),
		Slug:      exist.Slug,
		Name:      exist.Name,
		Status:    orgv1.OrganizationState(orgv1.OrganizationState_value[string(exist.Status.Status)]),
		CreatedAt: timestamppb.New(exist.CreatedAt),
		UpdatedAt: timestamppb.New(exist.UpdatedAt),
	}, nil
}

func (uc *organizationUsecase) ListOrganization(ctx context.Context, req *orgv1.ListOrganizationRequest) (*orgv1.ListOrganizationResponse, error) {

	page := req.Page
	opts := []option.QueryOption{
		option.ApplyPagination(
			pagination.Pagination{
				Cursor: page.Cursor,
				Limit:  int(page.Limit),
			},
		),
		option.WithSortBy(option.QuerySortBy{
			Allow: map[string]bool{
				"created_at": true,
			},
		}),
	}

	orgs, err := uc.orgRepo.Find(ctx, &domain.Organization{}, opts...)
	if err != nil {
		return nil, err
	}

	var newOrg []*orgv1.Organization
	for _, v := range orgs {
		newOrg = append(newOrg, &orgv1.Organization{
			OrgId:     v.ID.String(),
			Type:      orgv1.OrganizationType(orgv1.OrganizationType_value[v.Type]),
			Slug:      v.Slug,
			Name:      v.Name,
			CreatedAt: timestamppb.New(v.CreatedAt),
			UpdatedAt: timestamppb.New(v.UpdatedAt),
		})
	}

	return &orgv1.ListOrganizationResponse{
		Data: newOrg,
	}, nil
}

func (uc *organizationUsecase) UpdateOrganization(ctx context.Context, req *orgv1.UpdateOrganizationRequest) (*orgv1.Organization, error) {
	return nil, status.Error(codes.Unimplemented, "Unimplemeted")
}

func (uc *organizationUsecase) CreateLocation(ctx context.Context, req *orgv1.CreateLocationRequest) (*orgv1.Location, error) {

	orgID, err := snowflake.ParseString(req.OrgId)
	if err != nil {
		return nil, err
	}

	newLocation := domain.NewLocation(
		domain.LocationParams{
			OrgID:       orgID,
			Name:        req.Name,
			Address:     req.Address,
			City:        req.City,
			ZipCode:     req.ZipCode,
			CountryCode: req.CountryCode,
			Timezone:    req.Timezone,
		},
	)

	if err := uc.locationRepo.Create(ctx, newLocation); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return uc.GetLocation(ctx, &orgv1.GetLocationRequest{
		LocationId: newLocation.ID,
	})
}

func (uc *organizationUsecase) ListLocation(ctx context.Context, req *orgv1.ListLocationRequest) (*orgv1.ListLocationResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Unimplemeted")
}

func (uc *organizationUsecase) GetLocation(ctx context.Context, req *orgv1.GetLocationRequest) (*orgv1.Location, error) {
	exist, err := uc.locationRepo.FindOne(ctx, &domain.Location{
		ID: req.LocationId,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if exist == nil {
		return nil, fmt.Errorf("invalid locationId")
	}

	return &orgv1.Location{
		LocationId:  exist.ID,
		OrgId:       exist.OrgID.String(),
		Name:        exist.Name,
		Address:     exist.Address,
		City:        exist.City,
		ZipCode:     exist.ZipCode,
		CountryCode: exist.CountryCode,
		Timezone:    exist.Timezone,
	}, nil
}

func (uc *organizationUsecase) CreateInvitation(ctx context.Context, req *orgv1.CreateInvitationRequest) (*orgv1.CreateInvitationResponse, error) {

	orgID, err := snowflake.ParseString(req.OrgId)
	if err != nil {
		return nil, err
	}

	memberExist, err := uc.memberRepo.FindOne(ctx, &domain.Member{
		OrgID: orgID,
		Email: req.Email,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if memberExist != nil {
		return nil, status.Error(codes.InvalidArgument, "member already exist")
	}

	exist, err := uc.invitationRepo.FindOne(ctx, &domain.Invitation{
		Email: req.Email,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if exist != nil {
		return nil, status.Error(codes.AlreadyExists, "invitation already exist")
	}

	newInvitation := domain.NewInvitation(orgID, req.Email, req.Role)
	if err := uc.invitationRepo.Create(ctx, newInvitation); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &orgv1.CreateInvitationResponse{}, nil
}

func (uc *organizationUsecase) GetInvitation(ctx context.Context, req *orgv1.GetInvitationRequest) (*orgv1.Invitation, error) {
	return nil, status.Error(codes.Unimplemented, "Unimplemeted")
}

func (uc *organizationUsecase) VerifyInvitation(ctx context.Context, req *orgv1.VerifyInvitationRequest) (*orgv1.VerifyInvitationResponse, error) {
	exist, err := uc.invitationRepo.FindOne(ctx, &domain.Invitation{
		Token: req.Token,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if exist == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid code invitation")
	}

	now := time.Now()
	if !exist.ExpiryAt.Before(now) {
		newMember := domain.NewMember(
			domain.MemberParams{
				OrgID:  exist.OrgID,
				UserID: "",
				RoleID: exist.RoleID,
			},
		)
		if err := uc.memberRepo.Create(ctx, newMember); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return nil, status.Error(codes.Unimplemented, "Unimplemeted")
}
