package grpc_handler

import (
	"context"

	pointv1 "github.com/smallbiznis/go-genproto/smallbiznis/loyalty/point/v1"
	"github.com/smallbiznis/smallbiznis-apps/internal/loyalty/domain"
	"github.com/smallbiznis/smallbiznis-apps/pkg/workflow"
	"go.temporal.io/sdk/client"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PointHandler struct {
	pointv1.UnimplementedPointServiceServer
	workflow client.Client
}

func NewPointHandler(
	workflow client.Client,
) *PointHandler {
	return &PointHandler{
		workflow: workflow,
	}
}

func (h *PointHandler) Earning(ctx context.Context, req *pointv1.EarningRequest) (*pointv1.EarningResponse, error) {
	fields := []zap.Field{
		zap.String("organization_id", req.OrganizationId),
		zap.String("user_id", req.UserId),
		zap.String("reference_id", req.ReferenceId),
		zap.Any("attributes", req.Attributes),
	}

	zap.L().With(fields...).Info("Earning")

	transactionID, err := domain.GenerateTransactionID()
	if err != nil {
		return nil, err
	}

	transaction := domain.NewTransaction(domain.TransactionParam{
		OrganizationID: req.OrganizationId,
		UserID:         req.UserId,
		Type:           domain.EARNING,
		ReferenceID:    req.ReferenceId,
		TransactionID:  transactionID,
	})

	workflowOption := client.StartWorkflowOptions{
		ID:        transaction.ID,
		TaskQueue: workflow.POINT_TASK_QUEUE.String(),
	}

	req.ReferenceId = transaction.TransactionID
	w, err := h.workflow.ExecuteWorkflow(ctx, workflowOption, workflow.WorkflowEarnPoint, req)
	if err != nil {
		zap.L().With(fields...).Error("failed to ExecuteWorkflow", zap.Error(err))
		return nil, err
	}

	transaction.WorkflowID = w.GetID()

	return &pointv1.EarningResponse{
		TransactionId: transaction.TransactionID,
		Status:        pointv1.Status_PENDING,
		CreatedAt:     timestamppb.New(transaction.CreatedAt),
	}, nil
}

func (h *PointHandler) GetEarningStatus(ctx context.Context, req *pointv1.StatusEarningRequest) (*pointv1.StatusEarningResponse, error) {
	return &pointv1.StatusEarningResponse{}, nil
}

func (h *PointHandler) ValidateRedeem(ctx context.Context, req *pointv1.RedeemValidationRequest) (*pointv1.RedeemValidationResponse, error) {
	return &pointv1.RedeemValidationResponse{}, nil
}

func (h *PointHandler) Redemption(ctx context.Context, req *pointv1.RedeemRequest) (*pointv1.RedeemResponse, error) {
	return &pointv1.RedeemResponse{}, nil
}

func (h *PointHandler) GetRedemptionStatus(ctx context.Context, req *pointv1.StatusRedeemRequest) (*pointv1.StatusRedeemResponse, error) {
	return &pointv1.StatusRedeemResponse{}, nil
}
