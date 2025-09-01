package usecase

import (
	"context"
	"fmt"

	workflowv1 "github.com/smallbiznis/go-genproto/smallbiznis/workflow/v1"
	"github.com/smallbiznis/smallbiznis-apps/internal/workflow/domain"
	"github.com/smallbiznis/smallbiznis-apps/pkg/repository"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type FlowUsecase struct {
	db       *gorm.DB
	template repository.Repository[domain.FlowTemplate]
	flow     repository.Repository[domain.Flow]
}

type Params struct {
	fx.In
	DB       *gorm.DB
	Template repository.Repository[domain.FlowTemplate]
	Flow     repository.Repository[domain.Flow]
}

func NewFlowUsecase(p Params) *FlowUsecase {
	return &FlowUsecase{
		db:       p.DB,
		template: p.Template,
		flow:     p.Flow,
	}
}

func (u *FlowUsecase) ListFlowTemplate(ctx context.Context, req *workflowv1.ListFlowTemplatesRequest) (*workflowv1.ListFlowTemplatesResponse, error) {
	return &workflowv1.ListFlowTemplatesResponse{}, nil
}

func (u *FlowUsecase) GetFlowTemplate(ctx context.Context, req *workflowv1.GetFlowTemplateRequest) (*workflowv1.FlowTemplate, error) {
	template, err := u.template.FindOne(ctx, &domain.FlowTemplate{
		ID: req.GetId(),
	})
	if err != nil {
		zap.L().Error("failed to query get template", zap.Error(err))
		return nil, err
	}

	return &workflowv1.FlowTemplate{
		Id:          template.ID,
		Name:        template.Name,
		Description: template.Description,
		Status:      workflowv1.FlowStatus(workflowv1.FlowStatus_value[template.Status]),
		Nodes:       template.GetNodes(),
		Edges:       template.GetEdges(),
		Overflow:    template.GetOverflow(),
	}, nil
}

func (u *FlowUsecase) CreateFlow(ctx context.Context, req *workflowv1.CreateFlowRequest) (*workflowv1.Flow, error) {

	fields := []zap.Field{
		zap.String("organization_id", req.OrganizationId),
		zap.String("name", req.Name),
	}

	flow := domain.NewFlow(domain.FlowParams{
		OrganizationID: req.OrganizationId,
		Name:           req.Name,
	})

	if err := flow.SetNodes(req.Nodes); err != nil {
		zap.L().With(fields...).Error("failed setNodes", zap.Error(err))
		return nil, err
	}

	if err := flow.SetEdges(req.Edges); err != nil {
		zap.L().With(fields...).Error("failed setEdges", zap.Error(err))
		return nil, err
	}

	if err := flow.SetOverview(req.Overflow); err != nil {
		zap.L().With(fields...).Error("failed setOverview", zap.Error(err))
		return nil, err
	}

	if err := u.flow.Create(ctx, flow); err != nil {
		zap.L().With(fields...).Error("failed create flow", zap.Error(err))
		return nil, err
	}

	return u.GetFlow(ctx, &workflowv1.GetFlowRequest{Id: flow.ID})
}

func (u *FlowUsecase) GetFlow(ctx context.Context, req *workflowv1.GetFlowRequest) (*workflowv1.Flow, error) {
	flow, err := u.flow.FindOne(ctx, &domain.Flow{
		ID: req.GetId(),
	})
	if err != nil {
		zap.L().Error("failed to query get flow", zap.Error(err))
		return nil, err
	}

	return &workflowv1.Flow{
		Id:             flow.ID,
		OrganizationId: flow.OrganizationID,
		Name:           flow.Name,
		Description:    flow.Description,
		Status:         workflowv1.FlowStatus(workflowv1.FlowStatus_value[flow.Status]),
		Nodes:          flow.GetNodes(),
		Edges:          flow.GetEdges(),
		Overflow:       flow.GetOverflow(),
		CreatedAt:      timestamppb.New(flow.CreatedAt),
		UpdatedAt:      timestamppb.New(flow.UpdatedAt),
	}, nil
}

func (u *FlowUsecase) ListFlow(ctx context.Context, req *workflowv1.ListFlowsRequest) (*workflowv1.ListFlowsResponse, error) {
	flows, err := u.flow.Find(ctx, &domain.Flow{
		OrganizationID: req.OrganizationId,
		Status:         req.Status.String(),
		Trigger:        req.Trigger,
	})
	if err != nil {
		return nil, err
	}

	var result []*workflowv1.Flow
	for _, flow := range flows {
		result = append(result, &workflowv1.Flow{
			Id:             flow.ID,
			OrganizationId: flow.OrganizationID,
			Name:           flow.Name,
			Description:    flow.Description,
			Status:         workflowv1.FlowStatus(workflowv1.ActionType_value[flow.Status]),
			Nodes:          flow.GetNodes(),
			Edges:          flow.GetEdges(),
			Overflow:       flow.GetOverflow(),
			CreatedAt:      timestamppb.New(flow.CreatedAt),
			UpdatedAt:      timestamppb.New(flow.UpdatedAt),
		})
	}

	return &workflowv1.ListFlowsResponse{
		Data: result,
	}, nil
}

func (u *FlowUsecase) Updateflow(ctx context.Context, req *workflowv1.UpdateFlowRequest) (*workflowv1.Flow, error) {
	exist, err := u.flow.FindOne(ctx, &domain.Flow{
		ID: req.Id,
	})
	if err != nil {
		zap.L().Error("failed to query get flow", zap.Error(err))
		return nil, err
	}

	if exist == nil {
		zap.L().Error("flow not found")
		return nil, fmt.Errorf("flow not found")
	}

	update := domain.Flow{
		Name:        req.Flow.Name,
		Description: req.Flow.Description,
		Status:      req.Flow.Status.String(),
	}

	update.SetNodes(req.Flow.Nodes)
	update.SetEdges(req.Flow.Edges)
	update.SetOverview(req.Flow.Overflow)

	if err := u.flow.Update(ctx, exist.ID, &update); err != nil {
		zap.L().Error("failed to update flow", zap.Error(err), zap.String("flow_id", exist.ID))
		return nil, err
	}

	return &workflowv1.Flow{
		Id:             exist.ID,
		OrganizationId: exist.OrganizationID,
		Name:           exist.Name,
		Description:    exist.Description,
		Status:         workflowv1.FlowStatus(workflowv1.ActionType_value[exist.Status]),
		Nodes:          exist.GetNodes(),
		Edges:          exist.GetEdges(),
		Overflow:       exist.GetOverflow(),
		CreatedAt:      timestamppb.New(exist.CreatedAt),
		UpdatedAt:      timestamppb.New(exist.UpdatedAt),
	}, nil
}
