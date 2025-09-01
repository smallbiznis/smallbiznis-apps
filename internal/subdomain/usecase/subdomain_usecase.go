package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	subdomainv1 "github.com/smallbiznis/go-genproto/smallbiznis/subdomain/v1"
	"github.com/smallbiznis/smallbiznis-apps/internal/subdomain/domain"
	"github.com/smallbiznis/smallbiznis-apps/pkg/repository"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type ISubdomaiinUsecase interface {
	CreateSubdomain(context.Context, *subdomainv1.CreateSubdomainRequest) (*subdomainv1.Domain, error)
	GetSubdomain(context.Context, *subdomainv1.GetSubdomainRequest) (*subdomainv1.Domain, error)
	ListSubdomain(context.Context, *subdomainv1.ListSubdomainsRequest) (*subdomainv1.ListSubdomainsResponse, error)
	UpdateSubdomain(context.Context, *subdomainv1.UpdateSubdomainRequest) (*subdomainv1.Domain, error)
}

type subdomainUsecase struct {
	db        *gorm.DB
	subdomain repository.Repository[domain.Subdomain]
}

func NewSubdomainUsecase(
	db *gorm.DB,
	subdomain repository.Repository[domain.Subdomain],
) ISubdomaiinUsecase {
	return &subdomainUsecase{
		db:        db,
		subdomain: subdomain,
	}
}

func (uc *subdomainUsecase) CreateSubdomain(ctx context.Context, req *subdomainv1.CreateSubdomainRequest) (*subdomainv1.Domain, error) {
	newSubdomain := &domain.Subdomain{
		ID:             uuid.NewString(),
		OrganizationID: req.OrganizationId,
		Name:           req.Name,
	}

	if err := uc.subdomain.Create(ctx, newSubdomain); err != nil {
		return nil, err
	}

	return uc.GetSubdomain(ctx, &subdomainv1.GetSubdomainRequest{
		Id: newSubdomain.ID,
	})
}

func (uc *subdomainUsecase) GetSubdomain(ctx context.Context, req *subdomainv1.GetSubdomainRequest) (*subdomainv1.Domain, error) {

	exist, err := uc.subdomain.FindOne(ctx, &domain.Subdomain{
		ID: req.Id,
	})
	if err != nil {
		return nil, err
	}

	if exist == nil {
		return nil, fmt.Errorf("subdomain not found")
	}

	return &subdomainv1.Domain{
		DomainId:       exist.ID,
		Name:           exist.Name,
		OrganizationId: exist.OrganizationID,
		CreatedAt:      timestamppb.New(exist.CreatedAt),
		UpdatedAt:      timestamppb.New(exist.UpdatedAt),
	}, nil
}

func (uc *subdomainUsecase) ListSubdomain(ctx context.Context, req *subdomainv1.ListSubdomainsRequest) (*subdomainv1.ListSubdomainsResponse, error) {

	subdomains, err := uc.subdomain.Find(ctx, &domain.Subdomain{
		OrganizationID: req.OrganizationId,
	})
	if err != nil {
		return nil, err
	}

	data := make([]*subdomainv1.Domain, 0, len(subdomains))
	for _, v := range subdomains {
		data = append(data, &subdomainv1.Domain{
			DomainId:       v.ID,
			Name:           v.Name,
			OrganizationId: v.OrganizationID,
			CreatedAt:      timestamppb.New(v.CreatedAt),
			UpdatedAt:      timestamppb.New(v.UpdatedAt),
		})
	}

	return &subdomainv1.ListSubdomainsResponse{
		Data: data,
	}, nil
}

func (uc *subdomainUsecase) UpdateSubdomain(ctx context.Context, req *subdomainv1.UpdateSubdomainRequest) (*subdomainv1.Domain, error) {
	return nil, nil
}
