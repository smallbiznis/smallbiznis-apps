package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	ledgerv1 "github.com/smallbiznis/go-genproto/smallbiznis/ledger/v1"
	"github.com/smallbiznis/smallbiznis-apps/internal/ledger/domain"
	"github.com/smallbiznis/smallbiznis-apps/pkg/db/option"
	"github.com/smallbiznis/smallbiznis-apps/pkg/errutil"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

//go:generate mockgen -source=usecase.go -destination=./../../usecase/mock_ledger_usecase.go -package=usecase
type LedgerUsecase interface {
	AddEntry(ctx context.Context, req *ledgerv1.AddEntryRequest) (*ledgerv1.LedgerEntry, error)
	ReverEntry(ctx context.Context, req *ledgerv1.RevertEntryRequest) (*ledgerv1.LedgerEntry, error)
	ListEntries(ctx context.Context, req *ledgerv1.ListEntriesRequest) (*ledgerv1.ListEntriesResponse, error)
	GetEntry(ctx context.Context, req *ledgerv1.GetEntryRequest) (*ledgerv1.LedgerEntry, error)
	VerifyChain(ctx context.Context, req *ledgerv1.VerifyChainRequest) (*ledgerv1.VerifyChainResponse, error)
	GetBalance(ctx context.Context, req *ledgerv1.GetBalanceRequest) (*ledgerv1.GetBalanceResponse, error)
}

type ledgerUsecase struct {
	fx.In
	DB                   *gorm.DB
	LedgerRepository     domain.LedgerRepository
	CreditPoolRepository domain.CreditPoolRepository
	BalanceRepository    domain.BalanceRepository
}

func NewLedger(p ledgerUsecase) LedgerUsecase {
	return &p
}

func (s *ledgerUsecase) GetBalance(ctx context.Context, req *ledgerv1.GetBalanceRequest) (*ledgerv1.GetBalanceResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	traceID := span.SpanContext().TraceID().String()
	spanID := span.SpanContext().SpanID().String()

	opts := []zap.Field{
		zap.String("trace_id", traceID),
		zap.String("span_id", spanID),
	}

	lastEntry, err := s.BalanceRepository.FindOne(ctx, &domain.Balance{OrganizationID: req.OrganizationId, UserID: req.UserId}, option.WithSortBy(option.QuerySortBy{OrderBy: "DESC"}))
	if err != nil {
		zap.L().With(opts...).Error("failed to query FindOne entry", zap.Error(err))
		return nil, err
	}

	var lastBalance int64 = 0
	if lastEntry != nil {
		lastBalance = lastEntry.Balance
	}

	return &ledgerv1.GetBalanceResponse{
		Balance:       lastBalance,
		LastUpdatedAt: timestamppb.New(lastEntry.CreatedAt),
	}, nil
}

func (s *ledgerUsecase) AddEntry(ctx context.Context, req *ledgerv1.AddEntryRequest) (*ledgerv1.LedgerEntry, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	traceID := span.SpanContext().TraceID().String()
	spanID := span.SpanContext().SpanID().String()

	opts := []zap.Field{
		zap.String("trace_id", traceID),
		zap.String("span_id", spanID),
	}

	exist, err := s.LedgerRepository.FindOne(ctx, &domain.LedgerEntry{
		OrganizationID: req.OrganizationId,
		ReferenceID:    req.ReferenceId,
	})
	if err != nil {
		zap.L().With(opts...).Error("failed to query FindOne entry", zap.Error(err))
		return nil, err
	}

	if exist != nil {
		zap.L().With(opts...).Error("failed to create new entry", zap.Error(fmt.Errorf("reference_id %s already exists", req.ReferenceId)))
		return nil, errutil.BadRequest("failed to create new entry; reference_id already exists", nil)
	}

	if err := s.processAddEntry(ctx, req); err != nil {
		zap.L().Error("failed process add entry", zap.Error(err))
		return nil, err
	}

	entry, err := s.LedgerRepository.FindOne(ctx, &domain.LedgerEntry{
		ReferenceID: req.ReferenceId,
	})
	if err != nil {
		return nil, err
	}

	return &ledgerv1.LedgerEntry{
		Id:             entry.ID,
		OrganizationId: entry.OrganizationID,
		UserId:         entry.UserID,
		Type:           ledgerv1.EntryType(ledgerv1.EntryType_value[entry.Type]),
		Amount:         entry.Amount,
		TransactionId:  entry.TransactionID,
		ReferenceId:    entry.ReferenceID,
		Description:    entry.Description,
	}, nil
}

func (s *ledgerUsecase) ReverEntry(ctx context.Context, req *ledgerv1.RevertEntryRequest) (*ledgerv1.LedgerEntry, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	traceID := span.SpanContext().TraceID().String()
	spanID := span.SpanContext().SpanID().String()

	opts := []zap.Field{
		zap.String("trace_id", traceID),
		zap.String("span_id", spanID),
	}

	if err := s.DB.Transaction(func(tx *gorm.DB) error {

		originalEntry, err := s.LedgerRepository.FindOne(ctx, &domain.LedgerEntry{
			ID: req.EntryId,
		}, option.WithLockingUpdate())
		if err != nil {
			zap.L().With(opts...).Error("failed to query FindOne entry", zap.Error(err))
			return err
		}

		return s.processRevertCredit(ctx, tx, originalEntry)

	}); err != nil {
		return nil, err
	}

	current, err := s.LedgerRepository.FindOne(ctx, &domain.LedgerEntry{ID: req.EntryId})
	if err != nil {
		return nil, err
	}

	return &ledgerv1.LedgerEntry{
		Id:             current.ID,
		OrganizationId: current.OrganizationID,
		UserId:         current.UserID,
		Type:           ledgerv1.EntryType(ledgerv1.EntryType_value[current.Type]),
		Amount:         current.Amount,
		TransactionId:  current.TransactionID,
		ReferenceId:    current.ReferenceID,
		Description:    current.Description,
	}, nil
}

func (s *ledgerUsecase) getLastEntry(tx *gorm.DB, ctx context.Context, req *domain.LedgerEntry) (*domain.LedgerEntry, error) {
	lastEntry, err := s.LedgerRepository.WithTrx(tx).FindOne(ctx, &domain.LedgerEntry{
		OrganizationID: req.OrganizationID,
		UserID:         req.UserID,
	}, option.WithSortBy(
		option.QuerySortBy{
			SortBy:  "created_at",
			OrderBy: "desc",
			Allow: map[string]bool{
				"created_at": true,
			},
		},
	), option.WithLockingUpdate())
	if err != nil {
		return nil, err
	}

	return lastEntry, nil
}

func (s *ledgerUsecase) processAddEntry(ctx context.Context, req *ledgerv1.AddEntryRequest) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {

		lastEntry, err := s.getLastEntry(tx, ctx, &domain.LedgerEntry{
			OrganizationID: req.OrganizationId,
			UserID:         req.UserId,
		})
		if err != nil {
			return err
		}

		// Handle DEBIT
		if req.Type == ledgerv1.EntryType_DEBIT {
			return s.processDebit(ctx, tx, lastEntry, req)
		}

		// Handle CREDIT
		return s.processCredit(ctx, tx, lastEntry, req)
	})
}

func (s *ledgerUsecase) processDebit(ctx context.Context, tx *gorm.DB, lastEntry *domain.LedgerEntry, req *ledgerv1.AddEntryRequest) error {

	entries, err := s.CreditPoolRepository.WithTrx(tx).Find(ctx, &domain.CreditPool{
		OrganizationID: req.OrganizationId,
		UserID:         req.UserId,
	},
		option.ApplyOperator(option.Condition{
			Field:    "remaining",
			Operator: option.GT,
			Value:    0,
		}),
		option.WithSortBy(
			option.QuerySortBy{
				SortBy:  "created_at",
				OrderBy: "asc",
				Allow: map[string]bool{
					"created_at": true,
				},
			},
		),
		option.WithLockingUpdate(),
	)
	if err != nil {
		return err
	}

	if len(entries) == 0 {
		return fmt.Errorf("insufficient points")
	}

	transactionID, err := domain.GenerateTransactionID()
	if err != nil {
		zap.L().Error("failed to generate transactionId", zap.Error(err))
		return err
	}

	balance, err := s.BalanceRepository.FindOne(ctx, &domain.Balance{
		OrganizationID: req.OrganizationId,
		UserID:         req.UserId,
	},
		option.WithLockingUpdate(),
	)
	if err != nil {
		return err
	}

	if balance == nil {
		return fmt.Errorf("balance not found")
	}

	var totalAvailable int64
	for _, e := range entries {
		totalAvailable += e.Remaining
	}
	if totalAvailable < req.Amount {
		return fmt.Errorf("insufficient points: need=%d available=%d", req.Amount, totalAvailable)
	}

	remaining := req.Amount
	allocations := make([]domain.RedeemAllocation, 0, len(entries))
	for _, entry := range entries {
		if remaining == 0 {
			break
		}

		allocatable := min(entry.Remaining, remaining)
		allocations = append(allocations, domain.RedeemAllocation{
			CreditPoolID:    entry.ID,
			SourceID:        entry.LedgerEntryID,
			Amount:          allocatable,
			RemainingAmount: entry.Remaining - allocatable,
		})

		remaining -= allocatable
	}
	if remaining > 0 {
		return fmt.Errorf("insufficient points")
	}

	metadebit := make([]domain.MetaDebit, 0, len(allocations))
	for _, a := range allocations {
		metadebit = append(metadebit, domain.MetaDebit{
			LedgerEntryID: a.SourceID,
			Amount:        a.Amount,
		})
	}

	meta := make(map[string]any, len(req.Metadata)+1)
	for k, v := range req.Metadata {
		meta[k] = v
	}
	meta["sources"] = metadebit

	b, _ := json.Marshal(meta)
	entry := domain.NewLedgerEntry(domain.LedgerParams{
		Type:           ledgerv1.EntryType_DEBIT.String(),
		OrganizationID: req.OrganizationId,
		UserID:         req.UserId,
		Amount:         req.Amount,
		TransactionID:  transactionID,
		ReferenceID:    req.ReferenceId,
		Description:    req.Description,
		PreviousHash:   lastEntry.Hash,
		Metadata:       datatypes.JSON(b),
	})
	entry.Hash = entry.GenerateHash()

	if err := s.LedgerRepository.WithTrx(tx).Create(ctx, entry); err != nil {
		return err
	}

	for _, alloc := range allocations {
		updates := map[string]any{
			"remaining":   gorm.Expr("remaining - ?", alloc.Amount),
			"consumed_at": time.Now(),
		}
		if err := s.CreditPoolRepository.WithTrx(tx).Update(ctx, alloc.CreditPoolID, &updates); err != nil {
			zap.L().Error("failed to update credit pools", zap.Error(err))
			return err
		}
	}

	updates := map[string]any{
		"balance":    gorm.Expr("balance - ?", req.Amount),
		"updated_at": time.Now(),
	}
	if err := s.BalanceRepository.WithTrx(tx).Update(ctx, balance.ID, &updates); err != nil {
		return err
	}

	return nil
}

func (s *ledgerUsecase) processCredit(ctx context.Context, tx *gorm.DB, lastEntry *domain.LedgerEntry, req *ledgerv1.AddEntryRequest) error {
	var (
		previousHash    string = "GENESIS"
		previousBalance int64  = 0
	)

	balance, err := s.BalanceRepository.WithTrx(tx).FindOne(ctx, &domain.Balance{
		OrganizationID: req.OrganizationId,
		UserID:         req.UserId,
	}, option.WithLockingUpdate())
	if err != nil {
		zap.L().Error("failed to query balance", zap.Error(err))
		return err
	}

	transactionID, err := domain.GenerateTransactionID()
	if err != nil {
		zap.L().Error("failed to generate transactionId", zap.Error(err))
		return err
	}

	entry := domain.NewLedgerEntry(domain.LedgerParams{
		OrganizationID: req.OrganizationId,
		UserID:         req.UserId,
		Type:           req.Type.String(),
		Amount:         req.Amount,
		TransactionID:  transactionID,
		ReferenceID:    req.ReferenceId,
		Description:    req.Description,
	})

	if lastEntry != nil {
		previousHash = lastEntry.Hash
		previousBalance = balance.Balance
	}

	entry.PreviousHash = previousHash
	entry.Hash = entry.GenerateHash()

	if err := s.LedgerRepository.WithTrx(tx).Create(ctx, entry); err != nil {
		zap.L().Error("failed to create entry", zap.Error(err))
		return err
	}

	if err := s.CreditPoolRepository.WithTrx(tx).Create(ctx, &domain.CreditPool{
		ID:             uuid.NewString(),
		OrganizationID: req.OrganizationId,
		UserID:         req.UserId,
		LedgerEntryID:  entry.ID,
		Remaining:      req.Amount,
		CreatedAt:      time.Now(),
	}); err != nil {
		zap.L().Error("failed to create credit pools", zap.Error(err))
		return err
	}

	if balance == nil {
		if err := s.BalanceRepository.WithTrx(tx).Create(ctx, &domain.Balance{
			ID:             uuid.NewString(),
			OrganizationID: req.OrganizationId,
			UserID:         req.UserId,
			Balance:        entry.Amount,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}); err != nil {
			zap.L().Error("failed to create balance", zap.Error(err))
			return err
		}
	} else {
		if err := s.BalanceRepository.WithTrx(tx).Update(ctx, balance.ID, &domain.Balance{
			Balance:   entry.Amount + previousBalance,
			UpdatedAt: time.Now(),
		}); err != nil {
			zap.L().Error("failed to update balance", zap.Error(err))
			return err
		}
	}

	return nil
}

func (s *ledgerUsecase) processRevertCredit(ctx context.Context, tx *gorm.DB, lastEntry *domain.LedgerEntry) error {

	balance, err := s.BalanceRepository.FindOne(ctx, &domain.Balance{
		OrganizationID: lastEntry.OrganizationID,
		UserID:         lastEntry.UserID,
	},
		option.WithLockingUpdate(),
	)
	if err != nil {
		return err
	}

	if balance == nil {
		return fmt.Errorf("balance not found")
	}

	transactionID, err := domain.GenerateTransactionID()
	if err != nil {
		zap.L().Error("failed to generate transactionId", zap.Error(err))
		return err
	}

	entry := domain.NewLedgerEntry(domain.LedgerParams{
		OrganizationID: lastEntry.OrganizationID,
		UserID:         lastEntry.UserID,
		Type:           lastEntry.Type,
		Amount:         lastEntry.Amount,
		TransactionID:  transactionID,
		ReferenceID:    lastEntry.TransactionID,
		Description:    fmt.Sprintf("Revert of %s", lastEntry.ID),
	})

	entry.PreviousHash = lastEntry.Hash
	entry.Hash = entry.GenerateHash()

	if err := s.LedgerRepository.WithTrx(tx).Create(ctx, entry); err != nil {
		return err
	}

	return s.BalanceRepository.WithTrx(tx).Update(ctx, balance.ID, &domain.Balance{
		Balance:   balance.Balance - lastEntry.Amount,
		UpdatedAt: time.Now(),
	})
}

func (s *ledgerUsecase) ListEntries(ctx context.Context, req *ledgerv1.ListEntriesRequest) (*ledgerv1.ListEntriesResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	traceID := span.SpanContext().TraceID().String()
	spanID := span.SpanContext().SpanID().String()

	opts := []zap.Field{
		zap.String("trace_id", traceID),
		zap.String("span_id", spanID),
	}

	entries, err := s.LedgerRepository.Find(ctx, &domain.LedgerEntry{
		OrganizationID: req.OrganizationId,
		UserID:         req.UserId,
	})
	if err != nil {
		zap.L().With(opts...).Error("failed to query list entries", zap.Error(err))
		return nil, err
	}

	newEntries := make([]*ledgerv1.LedgerEntry, 0)
	for _, entry := range entries {
		newEntries = append(newEntries, &ledgerv1.LedgerEntry{
			Id:             entry.ID,
			OrganizationId: entry.OrganizationID,
			UserId:         entry.UserID,
			Type:           ledgerv1.EntryType(ledgerv1.EntryType_value[entry.Type]),
			Amount:         entry.Amount,
			TransactionId:  entry.TransactionID,
			ReferenceId:    entry.ReferenceID,
			Description:    entry.Description,
		})
	}

	return &ledgerv1.ListEntriesResponse{
		Data: newEntries,
	}, nil
}

func (s *ledgerUsecase) GetEntry(ctx context.Context, req *ledgerv1.GetEntryRequest) (*ledgerv1.LedgerEntry, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	traceID := span.SpanContext().TraceID().String()
	spanID := span.SpanContext().SpanID().String()

	opts := []zap.Field{
		zap.String("trace_id", traceID),
		zap.String("span_id", spanID),
	}

	entry, err := s.LedgerRepository.FindOne(ctx, &domain.LedgerEntry{
		ID: req.Id,
	})
	if err != nil {
		zap.L().With(opts...).Error("failed to FindOne entry", zap.Error(err))
		return nil, err
	}

	return &ledgerv1.LedgerEntry{
		Id:             entry.ID,
		OrganizationId: entry.OrganizationID,
		UserId:         entry.UserID,
		Type:           ledgerv1.EntryType(ledgerv1.EntryType_value[entry.Type]),
		Amount:         entry.Amount,
		TransactionId:  entry.TransactionID,
		ReferenceId:    entry.ReferenceID,
		Description:    entry.Description,
	}, nil
}

func (s *ledgerUsecase) VerifyChain(ctx context.Context, req *ledgerv1.VerifyChainRequest) (*ledgerv1.VerifyChainResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	traceID := span.SpanContext().TraceID().String()
	spanID := span.SpanContext().SpanID().String()

	opts := []zap.Field{
		zap.String("trace_id", traceID),
		zap.String("span_id", spanID),
	}

	entries, err := s.LedgerRepository.Find(ctx, &domain.LedgerEntry{
		OrganizationID: req.OrganizationId,
		UserID:         req.UserId,
	})
	if err != nil {
		zap.L().With(opts...).Error("failed to query Find entries", zap.Error(err))
		return nil, err
	}

	var lastHash string
	for _, entry := range entries {
		expectedHash := entry.GenerateHash()
		if entry.Hash != expectedHash || entry.PreviousHash != lastHash {
			return &ledgerv1.VerifyChainResponse{
				Valid: false,
			}, nil
		}
		lastHash = entry.Hash
	}

	return &ledgerv1.VerifyChainResponse{
		Valid: true,
	}, nil
}
