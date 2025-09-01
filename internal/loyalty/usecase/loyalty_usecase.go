package usecase

import (
	"context"

	pointv1 "github.com/smallbiznis/go-genproto/smallbiznis/loyalty/point/v1"
	"go.temporal.io/sdk/client"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

//go:generate mockgen -source=point_usecase.go -destination=mock_point_usecase.go -package=usecase
type ILoyaltyUsecase interface{}

type loyaltyUsecase struct {
	db       *gorm.DB
	temporal client.Client
}

func NewPointUsecase(
	db *gorm.DB,
	temporal client.Client,
) ILoyaltyUsecase {
	return &loyaltyUsecase{
		db,
		temporal,
	}
}

func (s *loyaltyUsecase) Earning(ctx context.Context) {}

func (s *loyaltyUsecase) StatusEarning(ctx context.Context) (*pointv1.StatusEarningResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Unimplemented")
}

func (s *loyaltyUsecase) GetRedemption(ctx context.Context) {}

func (s *loyaltyUsecase) Redemption(ctx context.Context, req *pointv1.RedeemRequest) (*pointv1.RedeemResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Unimplemented")
}
