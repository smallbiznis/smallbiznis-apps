package grpc

import (
	"context"

	subdomainv1 "github.com/smallbiznis/go-genproto/smallbiznis/subdomain/v1"
	"github.com/smallbiznis/smallbiznis-apps/internal/subdomain/usecase"
)

type SubdomainHandler struct {
	subdomainv1.UnimplementedSubdomainServiceServer
	subdomainUsecase usecase.ISubdomaiinUsecase
}

func NewSubdomainHandler(
	subdomainUsecase usecase.ISubdomaiinUsecase,
) *SubdomainHandler {
	return &SubdomainHandler{
		subdomainUsecase: subdomainUsecase,
	}
}

func (h *SubdomainHandler) CreateSubdomain(ctx context.Context, req *subdomainv1.CreateSubdomainRequest) (*subdomainv1.Domain, error) {
	return h.subdomainUsecase.CreateSubdomain(ctx, req)
}

func (h *SubdomainHandler) GetSubdomain(ctx context.Context, req *subdomainv1.GetSubdomainRequest) (*subdomainv1.Domain, error) {
	return h.subdomainUsecase.GetSubdomain(ctx, req)
}

func (h *SubdomainHandler) ListSubdomains(ctx context.Context, req *subdomainv1.ListSubdomainsRequest) (*subdomainv1.ListSubdomainsResponse, error) {
	return h.subdomainUsecase.ListSubdomain(ctx, req)
}

func (h *SubdomainHandler) UpdateSubdomain(ctx context.Context, req *subdomainv1.UpdateSubdomainRequest) (*subdomainv1.Domain, error) {
	return h.subdomainUsecase.UpdateSubdomain(ctx, req)
}
