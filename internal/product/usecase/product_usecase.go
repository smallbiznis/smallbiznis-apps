package usecase

import (
	"context"
	"errors"

	"github.com/bwmarrin/snowflake"
	"github.com/gosimple/slug"
	productv1 "github.com/smallbiznis/go-genproto/smallbiznis/product"
	"github.com/smallbiznis/smallbiznis-apps/internal/product/domain"
	"github.com/smallbiznis/smallbiznis-apps/pkg/db/option"
	"github.com/smallbiznis/smallbiznis-apps/pkg/db/pagination"
	"github.com/smallbiznis/smallbiznis-apps/pkg/gen"
	"github.com/smallbiznis/smallbiznis-apps/pkg/repository"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

//go:generate mockgen -source=product_usecase.go -destination=./../../usecase/mock_product_usecase.go -package=usecase
type ProductUsecase interface {
	ListProduct(context.Context, domain.Product, pagination.Pagination) (domain.Products, int64, error)
	GetProduct(context.Context, int64) (*domain.Product, error)
	CreateProduct(context.Context, *productv1.CreateProductRequest) (*domain.Product, error)
	UpdateProduct(context.Context, *productv1.UpdateProductRequest) (*domain.Product, error)
	DeleteProduct(context.Context, int64) error
}

type ProductParams struct {
	fx.In
	DB        *gorm.DB
	Snowflake *gen.SnowflakeNode
}

type productUsecase struct {
	db                *gorm.DB
	snowflake         *gen.SnowflakeNode
	optionRepo        repository.Repository[domain.Option]
	productRepo       repository.Repository[domain.Product]
	productOptionRepo repository.Repository[domain.ProductOption]
	variantRepo       repository.Repository[domain.Variant]
}

func NewProduct(p ProductParams) ProductUsecase {
	return &productUsecase{
		db:                p.DB,
		snowflake:         p.Snowflake,
		optionRepo:        repository.ProvideStore[domain.Option](p.DB),
		productRepo:       repository.ProvideStore[domain.Product](p.DB),
		productOptionRepo: repository.ProvideStore[domain.ProductOption](p.DB),
		variantRepo:       repository.ProvideStore[domain.Variant](p.DB),
	}
}

func (uc *productUsecase) ListProduct(ctx context.Context, filter domain.Product, pagination pagination.Pagination) (res domain.Products, count int64, err error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	traceID := span.SpanContext().TraceID().String()
	spanID := span.SpanContext().SpanID().String()

	traceOpt := []zap.Field{
		zap.String("trace_id", traceID),
		zap.String("span_id", spanID),
	}

	opts := []option.QueryOption{
		option.ApplyPagination(pagination),
		option.WithPreloads("Options", "Options.Values", "Variants", "Variants.Prices"),
	}

	products, err := uc.productRepo.Find(ctx, &filter, opts...)
	if err != nil {
		zap.L().With(traceOpt...).Error("Failed to find item", zap.Error(err))
		return nil, 0, status.Error(codes.Internal, err.Error())
	}

	totalData, err := uc.productRepo.Count(ctx, &filter)
	if err != nil {
		return nil, 0, status.Error(codes.Internal, err.Error())
	}

	return products, totalData, nil
}

func (uc *productUsecase) GetProduct(ctx context.Context, productID int64) (res *domain.Product, err error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	traceID := span.SpanContext().TraceID().String()
	spanID := span.SpanContext().SpanID().String()

	fields := []zap.Field{
		zap.String("trace_id", traceID),
		zap.String("span_id", spanID),
	}

	opts := []option.QueryOption{
		option.WithPreloads("Options", "Options.Values", "Variants", "Variants.Prices"),
	}

	product, err := uc.productRepo.FindOne(ctx, &domain.Product{ID: snowflake.ParseInt64(productID)}, opts...)
	if err != nil {
		zap.L().With(fields...).Error("Failed to find product", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	if product == nil {
		zap.L().With(fields...).Error("Product not found", zap.Int64("product_id", productID))
		return nil, status.Error(codes.NotFound, "Product not found")
	}

	return product, nil
}

func (uc *productUsecase) CreateProduct(ctx context.Context, req *productv1.CreateProductRequest) (res *domain.Product, err error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	traceID := span.SpanContext().TraceID().String()
	spanID := span.SpanContext().SpanID().String()

	fields := []zap.Field{
		zap.String("trace_id", traceID),
		zap.String("span_id", spanID),
	}

	orgID, err := snowflake.ParseString(req.OrgId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid org id")
	}

	newSlug := slug.Make(req.Title)
	exist, err := uc.productRepo.FindOne(ctx, &domain.Product{
		OrgID: orgID,
		Slug:  newSlug,
	})
	if err != nil {
		zap.L().With(fields...).Error("failed query get item", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	if exist != nil {
		zap.L().With(fields...).Error("item already exist", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, "item already exist!")
	}

	productID := uc.snowflake.GenerateID()
	newProduct := domain.Product{
		ID:       productID,
		OrgID:    orgID,
		Type:     req.Type.String(),
		Slug:     newSlug,
		Title:    req.Title,
		BodyHTML: req.BodyHtml,
		Status:   req.Status.String(),
	}

	if err := uc.db.WithContext(ctx).Transaction(func(tx *gorm.DB) (err error) {

		if len(req.Options) > 0 {
			// productOptions := make([]*domain.ProductOption, 0)
			for i, o := range req.Options {

				exist, err := uc.createOrUpdateOption(ctx, tx, domain.Option{
					OrgID: orgID,
					Name:  o.OptionName,
				})
				if err != nil {
					zap.L().With(fields...).Error("Failed to create or update option", zap.Error(err))
					return err
				}

				productOptionID := uc.snowflake.GenerateID()
				productOption, err := domain.NewProductOption(
					domain.ProductOption{
						ID:         productOptionID,
						ProductID:  productID,
						OptionName: exist.Name,
						Position:   (i + 1),
					},
				)
				if err != nil {
					return err
				}

				for _, v := range o.Values {
					productOptionValueID := uc.snowflake.GenerateID()
					productOptionValue, err := domain.NewProductOptionValue(
						domain.ProductOptionValue{
							ID:              productOptionValueID,
							ProductOptionID: productOptionID,
							Value:           v.Value,
						},
					)
					if err != nil {
						return err
					}

					optionValue := new(domain.OptionValue)
					if err := tx.Where(domain.OptionValue{
						OptionID: exist.ID,
						Value:    v.Value,
					}).First(optionValue).Error; err != nil {
						if !errors.Is(err, gorm.ErrRecordNotFound) {
							zap.L().With(fields...).Error("Failed to get option value", zap.Error(err))
							return err
						}

						optionValueID := uc.snowflake.GenerateID()
						optionValue = &domain.OptionValue{
							ID:       optionValueID,
							OptionID: exist.ID,
							Value:    v.Value,
						}

						if err = tx.Create(optionValue).Error; err != nil {
							zap.L().With(fields...).Error("Failed to create option value", zap.Error(err))
							return err
						}
					}

					productOption.Values = append(productOption.Values, productOptionValue)
				}

				newProduct.Options = append(newProduct.Options, productOption)
			}
		}

		if len(req.Variants) > 0 {

			variants := make([]*domain.Variant, 0)
			for _, variant := range req.Variants {
				variantID := uc.snowflake.GenerateID()
				newVariant, err := domain.NewVariant(domain.VariantParams{
					ID:        variantID,
					OrgID:     orgID,
					ProductID: productID,
					SKU:       variant.Sku,
					Title:     variant.Title,
					Taxable:   variant.Taxable,
				})
				if err != nil {
					return err
				}

				if len(variant.Prices) > 0 {
					for _, v := range variant.Prices {
						priceID := uc.snowflake.GenerateID()
						price := domain.NewPrice(
							domain.PriceParams{
								ID:           priceID,
								VariantID:    variantID,
								CurrencyCode: v.CurrencyCode,
								Cost:         v.Cost,
								Price:        v.Price,
								CompareAt:    v.CompareAtPrice,
							},
						)

						newVariant.Prices = append(newVariant.Prices, price)
					}
				}
				variants = append(variants, newVariant)

			}

			newProduct.Variants = append(newProduct.Variants, variants...)
		}

		if err = uc.productRepo.WithTrx(tx).Create(ctx, &newProduct); err != nil {
			zap.L().With(fields...).Error("Failed to create item", zap.Error(err))
			return err
		}

		return
	}); err != nil {
		zap.L().With(fields...).Error("Failed to create item", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return uc.GetProduct(ctx, newProduct.ID.Int64())
}

func (uc *productUsecase) UpdateProduct(ctx context.Context, data *productv1.UpdateProductRequest) (res *domain.Product, err error) {
	return nil, status.Error(codes.Unimplemented, "method UpdateProduct not implemented")
}

func (uc *productUsecase) DeleteProduct(ctx context.Context, id int64) (err error) {
	return status.Error(codes.Unimplemented, "method DeleteProduct not implemented")
}

func (uc *productUsecase) buildVariant(product domain.Product) (variants []domain.Variant) {

	// Helper function to recursively generate variants
	var createVariants func(int, []string, string)
	createVariants = func(optionIndex int, currentOptions []string, currentTitle string) {
		// if optionIndex == len(product.Options) {
		// 	variant := domain.Variant{
		// 		SKU:   "",
		// 		Title: currentTitle,
		// 		// Attributes: currentOptions,
		// 	}

		// 	for i, v := range currentOptions {
		// 		if int(i+1) == len(currentOptions) {
		// 			variant.SKU += fmt.Sprintf("-%s", v)
		// 		} else {
		// 			variant.SKU += strings.ToUpper(fmt.Sprintf("%s-%s", product.ItemCode(), v))
		// 		}
		// 	}

		// 	variants = append(variants, variant)
		// 	return
		// }

		// option := product.Options[optionIndex]
		// for _, value := range option.Values {
		// 	// Create a new map to avoid mutating the original one
		// 	newOptions := make([]string, 0)
		// 	newOptions = append(newOptions, currentOptions...)
		// 	newOptions = append(newOptions, value)

		// 	var newTitle string
		// 	if int(optionIndex+1) == len(product.Options) {
		// 		newTitle = fmt.Sprintf("%s, %s", currentTitle, value)
		// 	} else {
		// 		newTitle = fmt.Sprintf("%s - %s", currentTitle, value)
		// 	}
		// 	createVariants(optionIndex+1, newOptions, newTitle)
		// }
	}

	createVariants(0, []string{}, product.Title)

	return variants
}

func (uc *productUsecase) createOrUpdateOption(ctx context.Context, tx *gorm.DB, req domain.Option) (option *domain.Option, err error) {

	f := domain.Option{
		OrgID: req.OrgID,
	}

	if req.ID != 0 {
		f.ID = req.ID
	}

	if req.Name != "" {
		f.Name = req.Name
	}

	exist, err := uc.optionRepo.WithTrx(tx).FindOne(tx.Statement.Context, &f)
	if err != nil {
		return
	}

	if exist == nil {
		req.ID = uc.snowflake.GenerateID()
		if err := uc.optionRepo.WithTrx(tx).Create(ctx, &req); err != nil {
			return nil, err
		}

		return &req, nil
	}

	return exist, nil
}
