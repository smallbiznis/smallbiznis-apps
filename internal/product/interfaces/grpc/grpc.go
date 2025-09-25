package grpc_handler

import (
	"context"

	"github.com/bwmarrin/snowflake"
	productv1 "github.com/smallbiznis/go-genproto/smallbiznis/product"
	"github.com/smallbiznis/smallbiznis-apps/internal/product/domain"
	"github.com/smallbiznis/smallbiznis-apps/internal/product/usecase"
	"github.com/smallbiznis/smallbiznis-apps/pkg/db/pagination"
	"go.uber.org/fx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ProductHandler struct {
	productv1.UnimplementedProductServiceServer
	productUsecase usecase.ProductUsecase
}

type Params struct {
	fx.In
	ProductUsecase usecase.ProductUsecase
}

func NewProductHandler(p Params) *ProductHandler {
	return &ProductHandler{
		productUsecase: p.ProductUsecase,
	}
}

func (h *ProductHandler) CreateProduct(ctx context.Context, req *productv1.CreateProductRequest) (*productv1.Product, error) {
	product, err := h.productUsecase.CreateProduct(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return product.ToProto(), nil
}

func (h *ProductHandler) GetProduct(ctx context.Context, req *productv1.GetProductRequest) (*productv1.GetProductResponse, error) {
	productID, err := snowflake.ParseString(req.ProductId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	product, err := h.productUsecase.GetProduct(ctx, productID.Int64())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &productv1.GetProductResponse{
		Product: product.ToProto(),
	}, nil
}

func (h *ProductHandler) ListProducts(ctx context.Context, req *productv1.ListProductsRequest) (*productv1.ListProductsResponse, error) {

	orgID, err := snowflake.ParseString(req.OrgId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	products, count, err := h.productUsecase.ListProduct(ctx, domain.Product{
		OrgID: orgID,
	}, pagination.Pagination{
		Limit: int(req.Limit),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &productv1.ListProductsResponse{
		Products:   products.ToProto(),
		TotalCount: int32(count),
	}, nil
}

func (h *ProductHandler) UpdateProduct(ctx context.Context, req *productv1.UpdateProductRequest) (*productv1.Product, error) {
	return &productv1.Product{}, nil
}

func (h *ProductHandler) DeleteProduct(ctx context.Context, req *productv1.DeleteProductRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
