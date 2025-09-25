package domain

import (
	"time"

	"github.com/bwmarrin/snowflake"
	productv1 "github.com/smallbiznis/go-genproto/smallbiznis/product"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Option struct {
	ID        snowflake.ID `gorm:"column:id;primaryKey"`
	OrgID     snowflake.ID `gorm:"column:org_id;not null"`
	Name      string       `gorm:"column:name;not null"`
	CreatedAt time.Time    `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time    `gorm:"column:updated_at;autoUpdateTime"`
}

func (m *Option) ToProto() *productv1.Option {
	return &productv1.Option{
		OptionId:   m.ID.String(),
		OrgId:      m.OrgID.String(),
		OptionName: m.Name,
	}
}

type OptionValue struct {
	ID        snowflake.ID `gorm:"column:id;primaryKey"`
	OptionID  snowflake.ID `gorm:"column:option_id;not null"`
	Value     string       `gorm:"column:value;not null"`
	CreatedAt time.Time    `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time    `gorm:"column:updated_at;autoUpdateTime"`
}

type OptionValueParams struct {
	ID       snowflake.ID
	OptionID snowflake.ID
	Value    string
}

func NewOptionValue(p OptionValueParams) (*OptionValue, error) {
	return &OptionValue{
		ID:       p.ID,
		OptionID: p.OptionID,
		Value:    p.Value,
	}, nil
}

type Product struct {
	ID        snowflake.ID   `gorm:"column:id;primaryKey"`
	OrgID     snowflake.ID   `gorm:"column:org_id;not null"`
	Type      string         `gorm:"column:type;not null"`
	Title     string         `gorm:"column:title;not null"`
	Slug      string         `gorm:"column:slug;not null"`
	BodyHTML  string         `gorm:"column:body_html"`
	Options   ProductOptions `gorm:"foreignKey:ProductID"`
	Variants  Variants       `gorm:"foreignKey:ProductID"`
	Taxable   bool           `gorm:"column:taxable"`
	Status    string         `gorm:"column:status;not null"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
}

func (m *Product) ToProto() *productv1.Product {
	return &productv1.Product{
		ProductId: m.ID.String(),
		OrgId:     m.OrgID.String(),
		Type:      productv1.ProductType(productv1.ProductType_value[m.Type]),
		Title:     m.Title,
		Slug:      m.Slug,
		BodyHtml:  m.BodyHTML,
		Options:   m.Options.ToProto(),
		Variants:  m.Variants.ToProto(),
		Status:    productv1.ProductStatus(productv1.ProductStatus_value[m.Status]),
		CreatedAt: timestamppb.New(m.CreatedAt),
		UpdatedAt: timestamppb.New(m.UpdatedAt),
	}
}

type Products []*Product

func (m Products) ToProto() []*productv1.Product {
	var res []*productv1.Product
	for _, v := range m {
		res = append(res, v.ToProto())
	}
	return res
}

type ProductOption struct {
	ID         snowflake.ID        `gorm:"column:id;primaryKey"`
	ProductID  snowflake.ID        `gorm:"column:product_id;not null"`
	OptionName string              `gorm:"column:option_name;not null"`
	Position   int                 `gorm:"column:position;not null"`
	Values     ProductOptionValues `gorm:"foreignKey:ProductOptionID"`
}

func (m *ProductOption) ToProto() *productv1.ProductOption {
	return &productv1.ProductOption{
		OptionName: m.OptionName,
		Position:   int32(m.Position),
		Values:     m.Values.ToProto(),
	}
}

func NewProductOption(p ProductOption) (*ProductOption, error) {
	return &ProductOption{
		ID:         p.ID,
		ProductID:  p.ProductID,
		OptionName: p.OptionName,
		Position:   p.Position,
		Values:     p.Values,
	}, nil
}

type ProductOptions []*ProductOption

func (m ProductOptions) ToProto() []*productv1.ProductOption {
	var res []*productv1.ProductOption
	for _, v := range m {
		res = append(res, v.ToProto())
	}
	return res
}

type ProductOptionValue struct {
	ID              snowflake.ID
	ProductOptionID snowflake.ID
	Value           string
}

func (m *ProductOptionValue) ToProto() *productv1.ProductOptionValue {
	return &productv1.ProductOptionValue{
		Value: m.Value,
	}
}

func NewProductOptionValue(p ProductOptionValue) (*ProductOptionValue, error) {
	return &ProductOptionValue{
		ID:              p.ID,
		ProductOptionID: p.ProductOptionID,
		Value:           p.Value,
	}, nil
}

type ProductOptionValues []*ProductOptionValue

func (m ProductOptionValues) ToProto() []*productv1.ProductOptionValue {
	var res []*productv1.ProductOptionValue
	for _, v := range m {
		res = append(res, v.ToProto())
	}
	return res
}

type Variant struct {
	ID        snowflake.ID `gorm:"column:id;primaryKey"`
	OrgID     snowflake.ID `gorm:"column:org_id;not null"`
	ProductID snowflake.ID `gorm:"column:product_id"`
	SKU       string       `gorm:"column:sku;not null"`
	Title     string       `gorm:"column:title;not null"`
	Taxable   bool         `gorm:"column:taxable;not null"`
	Prices    Prices       `gorm:"foreignKey:VariantID"`
	Dimension Dimension    `gorm:"foreignKey:VariantID"`
	Attribute []*Attribute `gorm:"foreignKey:VariantID"`
	CreatedAt time.Time    `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time    `gorm:"column:updated_at;autoUpdateTime"`
}

type VariantParams struct {
	ID        snowflake.ID
	OrgID     snowflake.ID
	ProductID snowflake.ID
	SKU       string
	Title     string
	Taxable   bool
	Prices    []*Price
	Dimension Dimension
}

func NewVariant(p VariantParams) (*Variant, error) {
	return &Variant{
		ID:        p.ID,
		OrgID:     p.OrgID,
		ProductID: p.ProductID,
		SKU:       p.SKU,
		Title:     p.Title,
		Taxable:   p.Taxable,
		Prices:    p.Prices,
		Dimension: p.Dimension,
	}, nil
}

func (m *Variant) ToProto() *productv1.Variant {
	return &productv1.Variant{
		VariantId: m.ID.String(),
		ProductId: m.ProductID.String(),
		Sku:       m.SKU,
		Title:     m.Title,
		Taxable:   m.Taxable,
		Prices:    m.Prices.ToProto(),
	}
}

type Variants []*Variant

func (m Variants) ToProto() []*productv1.Variant {
	var res []*productv1.Variant
	for _, v := range m {
		res = append(res, v.ToProto())
	}
	return res
}

type Attribute struct {
	ID        snowflake.ID `gorm:"column:id;primaryKey"`
	VariantID snowflake.ID `gorm:"column:variant_id;not null"`
	Key       string       `gorm:"column:key"`
	Value     string       `gorm:"column:value"`
}

type Price struct {
	ID        snowflake.ID `gorm:"column:id;primaryKey"`
	VariantID snowflake.ID `gorm:"column:variant_id;not null"`
	Currency  string       `gorm:"column:currency"`
	Price     float64      `gorm:"column:price;not null"`
	CompareAt float64      `gorm:"column:compare_at_price;not null"`
	Cost      float64      `gorm:"column:cost;not null"`
	CreatedAt time.Time    `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time    `gorm:"column:updated_at;autoUpdateTime"`
}

type PriceParams struct {
	ID        snowflake.ID
	VariantID snowflake.ID
	Currency  string
	Price     float64
	CompareAt float64
	Cost      float64
}

func NewPrice(p PriceParams) *Price {
	return &Price{
		ID:        p.ID,
		VariantID: p.VariantID,
		Currency:  p.Currency,
		Price:     p.Price,
		Cost:      p.Cost,
		CompareAt: p.CompareAt,
	}
}

func (m *Price) ToProto() *productv1.Price {
	return &productv1.Price{
		PriceId:        m.ID.String(),
		VariantId:      m.VariantID.String(),
		Currency:       m.Currency,
		Price:          m.Price,
		Cost:           m.Cost,
		CompareAtPrice: m.CompareAt,
	}
}

type Prices []*Price

func (m Prices) ToProto() []*productv1.Price {
	var res []*productv1.Price
	for _, v := range m {
		res = append(res, v.ToProto())
	}
	return res
}

type Dimension struct {
	ID         snowflake.ID `gorm:"column:id;primaryKey"`
	VariantID  snowflake.ID `gorm:"column:variant_id;not null"`
	Weight     float64      `gorm:"column:weight;not null"`
	WeightUnit string       `gorm:"column:weight_unit"`
	Height     float64      `gorm:"column:height;not null"`
	HeightUnit string       `gorm:"column:height_unit"`
	Depth      float64      `gorm:"column:depth;not null"`
	CreatedAt  time.Time    `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time    `gorm:"column:updated_at;autoUpdateTime"`
}

type DimensionParams struct {
	ID         snowflake.ID
	VariantID  snowflake.ID
	Weight     float64
	WeightUnit string
	Height     float64
	HeightUnit string
	Depth      float64
}

func NewDimension(p DimensionParams) (*Dimension, error) {
	return &Dimension{
		ID:         p.ID,
		VariantID:  p.VariantID,
		Weight:     p.Weight,
		WeightUnit: p.WeightUnit,
		Height:     p.Height,
		HeightUnit: p.HeightUnit,
		Depth:      p.Depth,
	}, nil
}
