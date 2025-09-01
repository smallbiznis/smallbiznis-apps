package grpc_handler

import (
	"context"

	workflowv1 "github.com/smallbiznis/go-genproto/smallbiznis/workflow/v1"
	"github.com/smallbiznis/smallbiznis-apps/internal/workflow/usecase"
	"go.uber.org/fx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FlowHandler struct {
	workflowv1.UnimplementedWorkflowServiceServer
	workflowUsecase *usecase.FlowUsecase
}

type Params struct {
	fx.In
	Workflow *usecase.FlowUsecase
}

func NewFlowHandler(p Params) *FlowHandler {
	return &FlowHandler{
		workflowUsecase: p.Workflow,
	}
}

func (h *FlowHandler) ListFlowTemplate(ctx context.Context, req *workflowv1.ListFlowTemplatesRequest) (*workflowv1.ListFlowTemplatesResponse, error) {
	return h.workflowUsecase.ListFlowTemplate(ctx, req)
}

func (h *FlowHandler) GetFlowTemplate(ctx context.Context, req *workflowv1.GetFlowTemplateRequest) (*workflowv1.FlowTemplate, error) {
	return h.workflowUsecase.GetFlowTemplate(ctx, req)
}

func (h *FlowHandler) CreateFlow(ctx context.Context, req *workflowv1.CreateFlowRequest) (*workflowv1.Flow, error) {
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organizationId is required")
	}

	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	if len(req.Nodes) == 0 {
		return nil, status.Error(codes.InvalidArgument, "nodes can't empty")
	}

	if req.Nodes[0].Type != workflowv1.NodeType_TRIGGER {
		return nil, status.Error(codes.InvalidArgument, "first node type must TRIGGER")
	}

	if len(req.Edges) == 0 {
		return nil, status.Error(codes.InvalidArgument, "edges can't empty")
	}

	if req.GetOverflow() == nil {
		return nil, status.Error(codes.InvalidArgument, "overflow can't empty")
	}

	return h.workflowUsecase.CreateFlow(ctx, req)
}

func (h *FlowHandler) GetFlow(ctx context.Context, req *workflowv1.GetFlowRequest) (*workflowv1.Flow, error) {
	return h.workflowUsecase.GetFlow(ctx, req)
}

func (h *FlowHandler) ListFlows(ctx context.Context, req *workflowv1.ListFlowsRequest) (*workflowv1.ListFlowsResponse, error) {
	return h.workflowUsecase.ListFlow(ctx, req)
}

func (h *FlowHandler) Updateflow(ctx context.Context, req *workflowv1.UpdateFlowRequest) (*workflowv1.Flow, error) {
	if req.Flow.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organizationId is required")
	}

	if req.Flow.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	if len(req.Flow.Nodes) == 0 {
		return nil, status.Error(codes.InvalidArgument, "nodes can't empty")
	}

	if len(req.Flow.Edges) == 0 {
		return nil, status.Error(codes.InvalidArgument, "edges can't empty")
	}

	if req.Flow.GetOverflow() == nil {
		return nil, status.Error(codes.InvalidArgument, "overflow can't empty")
	}

	return h.workflowUsecase.Updateflow(ctx, req)
}
