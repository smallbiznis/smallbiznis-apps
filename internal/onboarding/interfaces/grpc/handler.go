package grpc_handler

import (
	"context"

	ledgerv1 "github.com/smallbiznis/go-genproto/smallbiznis/ledger/v1"
	"github.com/smallbiznis/smallbiznis-apps/internal/ledger/usecase"
	"go.uber.org/fx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	ledgerv1.UnimplementedLedgerServiceServer
	ledgerUsecase usecase.LedgerUsecase
}

type Params struct {
	fx.In
	LedgerUsecase usecase.LedgerUsecase
}

func NewHandler(p Params) *Handler {
	return &Handler{
		ledgerUsecase: p.LedgerUsecase,
	}
}

func (h *Handler) AddEntry(ctx context.Context, req *ledgerv1.AddEntryRequest) (*ledgerv1.LedgerEntry, error) {

	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organizationId is required")
	}

	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "userId is required")
	}

	if req.Amount <= 0 {
		return nil, status.Error(codes.InvalidArgument, "amount must be greater than 0")
	}

	if req.ReferenceId == "" {
		return nil, status.Error(codes.InvalidArgument, "referenceId is required")
	}

	return h.ledgerUsecase.AddEntry(ctx, req)
}

func (h *Handler) RevertEntry(ctx context.Context, req *ledgerv1.RevertEntryRequest) (*ledgerv1.LedgerEntry, error) {
	return h.ledgerUsecase.ReverEntry(ctx, req)
}

func (h *Handler) ListEntries(ctx context.Context, req *ledgerv1.ListEntriesRequest) (*ledgerv1.ListEntriesResponse, error) {
	return h.ledgerUsecase.ListEntries(ctx, req)
}

func (h *Handler) GetEntry(ctx context.Context, req *ledgerv1.GetEntryRequest) (*ledgerv1.LedgerEntry, error) {
	return h.ledgerUsecase.GetEntry(ctx, req)
}

func (h *Handler) GetBalance(ctx context.Context, req *ledgerv1.GetBalanceRequest) (*ledgerv1.GetBalanceResponse, error) {
	return h.ledgerUsecase.GetBalance(ctx, req)
}

func (h *Handler) VerifyChain(ctx context.Context, req *ledgerv1.VerifyChainRequest) (*ledgerv1.VerifyChainResponse, error) {
	return h.ledgerUsecase.VerifyChain(ctx, req)
}
