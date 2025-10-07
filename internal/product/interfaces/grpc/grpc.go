package grpc_handler

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/smallbiznis/go-genproto/smallbiznis/common"
	productv1 "github.com/smallbiznis/go-genproto/smallbiznis/product"
	"github.com/smallbiznis/smallbiznis-apps/internal/product/domain"
	"github.com/smallbiznis/smallbiznis-apps/internal/product/usecase"
	"github.com/smallbiznis/smallbiznis-apps/pkg/db/pagination"
	"github.com/smallbiznis/smallbiznis-apps/pkg/errutil"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
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
		return nil, errutil.ToGRPCError(err)
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
		return nil, errutil.ToGRPCError(err)
	}

	return &productv1.GetProductResponse{
		Product: product.ToProto(),
	}, nil
}

func (h *ProductHandler) ListProducts(ctx context.Context, req *productv1.ListProductsRequest) (*productv1.ListProductsResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	traceID := span.SpanContext().TraceID().String()
	spanID := span.SpanContext().SpanID().String()

	traceOpt := []zap.Field{
		zap.String("trace_id", traceID),
		zap.String("span_id", spanID),
	}

	if req.OrgId == "" {
		return nil, errutil.ToGRPCError(fmt.Errorf("orgId is required"))
	}

	orgID, err := snowflake.ParseString(req.OrgId)
	if err != nil {
		zap.L().With(
			traceOpt...,
		).Error("failed parse orgId", zap.Error(err))
		return nil, errutil.ToGRPCError(err)
	}

	page := req.Page
	paginate := pagination.Pagination{
		Limit:  int(page.Limit),
		Cursor: page.Cursor,
	}

	if page.Cursor != "" {
		paginate.Cursor = page.Cursor
	}

	// Service
	products, count, err := h.productUsecase.ListProduct(ctx, domain.Product{
		OrgID: orgID,
	}, paginate)
	if err != nil {
		zap.L().With(
			traceOpt...,
		).Error("failed get list product", zap.Error(err))
		return nil, errutil.ToGRPCError(err)
	}

	// Parse to Proto
	newProducts := make([]*productv1.Product, 0)
	for _, v := range products {
		newProducts = append(newProducts, v.ToProto())
	}

	// Build Cursor Pagination
	pageInfo := pagination.BuildCursorPageInfo(products, page.Limit, func(t *domain.Product) string {
		cursor, err := pagination.EncodeCursor(
			pagination.Cursor{
				ID:        t.ID.String(),
				CreatedAt: t.CreatedAt.Format(time.DateTime),
			},
		)
		if err != nil {
			zap.L().With(
				traceOpt...,
			).Error("failed encode cursor", zap.Error(err))
			return ""
		}
		return cursor
	})

	var cursorResponse common.CursorResponse
	if pageInfo != nil {
		cursorResponse = common.CursorResponse{
			NextCursor: pageInfo.NextCursor,
			PrevCursor: pageInfo.PreviousCursor,
			HasMore:    pageInfo.HasMore,
			TotalCount: int32(count),
		}
	}

	return &productv1.ListProductsResponse{
		Data: newProducts,
		Page: &cursorResponse,
	}, nil
}

func (h *ProductHandler) UpdateProduct(ctx context.Context, req *productv1.UpdateProductRequest) (*productv1.Product, error) {
	return &productv1.Product{}, nil
}

func (h *ProductHandler) DeleteProduct(ctx context.Context, req *productv1.DeleteProductRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
