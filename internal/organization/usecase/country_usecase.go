package usecase

import (
	"context"
	"fmt"

	orgv1 "github.com/smallbiznis/go-genproto/smallbiznis/organization/v1"
	"github.com/smallbiznis/smallbiznis-apps/internal/organization/domain"
	"github.com/smallbiznis/smallbiznis-apps/pkg/repository"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type ICountry interface {
	ListCountry(context.Context, *orgv1.ListCountriesRequest) (*orgv1.ListCountriesResponse, error)
	GetCountry(ctx context.Context, code string) (*domain.Country, error)
}

type CountryParams struct {
	fx.In
	DB *gorm.DB
}

type country struct {
	db          *gorm.DB
	countryRepo repository.Repository[domain.Country]
}

func NewCountry(p CountryParams) ICountry {
	return &country{
		db:          p.DB,
		countryRepo: repository.ProvideStore[domain.Country](p.DB),
	}
}

func (uc *country) ListCountry(ctx context.Context, req *orgv1.ListCountriesRequest) (*orgv1.ListCountriesResponse, error) {

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

func (uc *country) GetCountry(ctx context.Context, code string) (*domain.Country, error) {

	country, err := uc.countryRepo.FindOne(ctx, &domain.Country{
		Code: code,
	})
	if err != nil {
		return nil, err
	}

	if country == nil {
		return nil, fmt.Errorf("invalid countryCode")
	}

	return country, nil
}
