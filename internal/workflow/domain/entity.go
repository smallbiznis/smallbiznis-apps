package domain

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	workflowv1 "github.com/smallbiznis/go-genproto/smallbiznis/workflow/v1"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/structpb"
	"gorm.io/datatypes"
)

type FlowTemplate struct {
	ID          string         `gorm:"column:id;primaryKey"`
	Name        string         `gorm:"column:name"`
	Description string         `gorm:"column:description"`
	Trigger     string         `gorm:"column:trigger"`
	Status      string         `gorm:"column:status"`
	Nodes       datatypes.JSON `gorm:"column:nodes"` // serialized []Node
	Edges       datatypes.JSON `gorm:"column:edges"` // serialized []Edge
	Overflow    datatypes.JSON `gorm:"column:overflow"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (m *FlowTemplate) SetNodes(nodes []*workflowv1.Node) error {

	var internalNodes []Node

	for _, n := range nodes {
		node := Node{
			ID:   n.GetId(),
			Type: n.GetType(),
			Position: Position{
				X: float64(n.Position.GetX()),
				Y: float64(n.Position.GetY()),
			},
		}

		switch n.Type {
		case workflowv1.NodeType_TRIGGER:

			b, err := n.Data.MarshalJSON()
			if err != nil {
				return err
			}

			var trigger NodeTrigger
			if err := json.Unmarshal(b, &trigger); err != nil {
				return err
			}

			m.Trigger = trigger.Key
			node.Trigger = &trigger

		case workflowv1.NodeType_CONDITION:
			b, err := n.Data.MarshalJSON()
			if err != nil {
				return err
			}

			var cond NodeCondition
			if err := json.Unmarshal(b, &cond); err != nil {
				return err
			}

			conditions := make([]*Condition, 0)
			for _, c := range cond.Conditions {
				conditions = append(conditions, &Condition{
					Field:    c.Field,
					Operator: c.Operator,
					Value:    c.Value,
				})
			}

			node.Condition = &NodeCondition{
				Conditions: conditions,
			}

		case workflowv1.NodeType_ACTION:
			b, err := n.Data.MarshalJSON()
			if err != nil {
				return err
			}

			var act NodeAction
			if err := json.Unmarshal(b, &act); err != nil {
				return err
			}

			node.Action = &NodeAction{
				Type:       act.Type,
				Parameters: act.Parameters,
			}
		}

		internalNodes = append(internalNodes, node)
	}

	data, err := json.Marshal(internalNodes)
	if err != nil {
		zap.L().Error("failed to marshal internalNodes", zap.Error(err))
		return err
	}

	m.Nodes = datatypes.JSON(data)
	return nil
}

func (m *FlowTemplate) SetEdges(edges []*workflowv1.Edge) error {

	var newEdges []*Edge
	for _, e := range edges {
		edge := &Edge{
			ID:     e.Id,
			Type:   e.Type,
			Source: e.Source,
			Target: e.Target,
		}

		newEdges = append(newEdges, edge)
	}

	b, err := json.Marshal(newEdges)
	if err != nil {
		zap.L().Error("failed to marshal edges", zap.Error(err))
		return err
	}

	m.Edges = datatypes.JSON(b)
	return nil
}

func (m *FlowTemplate) SetOverview(overflow *workflowv1.Overflow) error {
	b, err := json.Marshal(overflow)
	if err != nil {
		zap.L().Error("failed to marshal overflow", zap.Error(err))
		return err
	}

	m.Overflow = datatypes.JSON(b)
	return nil
}

func (m *FlowTemplate) GetNodes() []*workflowv1.Node {
	var rawNode []*Node
	if err := json.Unmarshal(m.Nodes, &rawNode); err != nil {
		zap.L().Error("failed to unmarshal Nodes", zap.Error(err))
		return nil
	}

	var nodes []*workflowv1.Node
	for _, n := range rawNode {

		node := &workflowv1.Node{
			Id:   n.ID,
			Type: n.Type,
			Position: &workflowv1.Position{
				X: float32(n.Position.X),
				Y: float32(n.Position.Y),
			},
		}

		switch n.Type {
		case workflowv1.NodeType_TRIGGER:

			b, err := json.Marshal(&n.Trigger)
			if err != nil {
				return nil
			}

			var data map[string]any
			if err := json.Unmarshal(b, &data); err != nil {
				return nil
			}

			obj, err := structpb.NewStruct(data)
			if err != nil {
				return nil
			}

			node.Data = obj

		case workflowv1.NodeType_CONDITION:
			b, err := json.Marshal(&n.Condition)
			if err != nil {
				return nil
			}

			var data map[string]any
			if err := json.Unmarshal(b, &data); err != nil {
				return nil
			}

			obj, err := structpb.NewStruct(data)
			if err != nil {
				return nil
			}

			node.Data = obj
		case workflowv1.NodeType_ACTION:
			b, err := json.Marshal(&n.Action)
			if err != nil {
				return nil
			}

			var data map[string]any
			if err := json.Unmarshal(b, &data); err != nil {
				return nil
			}

			obj, err := structpb.NewStruct(data)
			if err != nil {
				return nil
			}

			node.Data = obj
		}

		nodes = append(nodes, node)
	}

	return nodes
}

func (m *FlowTemplate) GetEdges() []*workflowv1.Edge {
	var edges []*workflowv1.Edge
	if err := json.Unmarshal(m.Edges, &edges); err != nil {
		zap.L().Error("failed to unmarshal Nodes", zap.Error(err))
		return nil
	}
	return edges
}

func (m *FlowTemplate) GetOverflow() *workflowv1.Overflow {
	var overflow *workflowv1.Overflow
	if err := json.Unmarshal(m.Overflow, &overflow); err != nil {
		zap.L().Error("failed to unmarshal Nodes", zap.Error(err))
		return nil
	}
	return overflow
}

type Flow struct {
	ID             string         `gorm:"column:id;primaryKey"`
	OrganizationID string         `gorm:"column:organization_id"`
	Name           string         `gorm:"column:name"`
	Description    string         `gorm:"column:description"`
	Trigger        string         `gorm:"column:trigger"`
	Status         string         `gorm:"column:status"` // enum FlowStatus as string
	Nodes          datatypes.JSON `gorm:"column:nodes"`  // serialized []Node
	Edges          datatypes.JSON `gorm:"column:edges"`  // serialized []Edge
	Overflow       datatypes.JSON `gorm:"column:overflow"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

type FlowParams struct {
	OrganizationID string
	Name           string
}

func NewFlow(p FlowParams) *Flow {
	return &Flow{
		ID:             uuid.NewString(),
		OrganizationID: p.OrganizationID,
		Name:           p.Name,
		Status:         workflowv1.FlowStatus_ACTIVE.String(),
	}
}

type FlowGraph struct {
	Nodes map[string]*Node
	Edges map[string][]*Edge
}

func (m *Flow) BuildFlowGraph() (*FlowGraph, error) {
	graph := &FlowGraph{
		Nodes: map[string]*Node{},
		Edges: map[string][]*Edge{},
	}

	for _, n := range m.GetNodes() {
		node := &Node{
			ID:   n.GetId(),
			Type: n.GetType(),
			Position: Position{
				X: float64(n.Position.GetX()),
				Y: float64(n.Position.GetY()),
			},
		}

		switch n.Type {
		case workflowv1.NodeType_TRIGGER:
			b, err := n.Data.MarshalJSON()
			if err != nil {
				return nil, err
			}

			var trigger NodeTrigger
			if err := json.Unmarshal(b, &trigger); err != nil {
				return nil, err
			}

			node.Trigger = &trigger

		case workflowv1.NodeType_CONDITION:
			b, err := n.Data.MarshalJSON()
			if err != nil {
				return nil, err
			}

			var cond NodeCondition
			if err := json.Unmarshal(b, &cond); err != nil {
				return nil, err
			}

			node.Condition = &cond

		case workflowv1.NodeType_ACTION:
			b, err := n.Data.MarshalJSON()
			if err != nil {
				return nil, err
			}

			var act NodeAction
			if err := json.Unmarshal(b, &act); err != nil {
				return nil, err
			}

			node.Action = &act
		}

		graph.Nodes[n.Id] = node
	}

	for _, e := range m.GetEdges() {
		edge := &Edge{
			ID:     e.Id,
			Type:   e.Type,
			Source: e.Source,
			Target: e.Target,
		}

		graph.Edges[e.Source] = append(graph.Edges[e.Source], edge)
	}

	return graph, nil
}

func (m *Flow) SetNodes(nodes []*workflowv1.Node) error {

	var internalNodes []Node

	for _, n := range nodes {
		node := Node{
			ID:   n.GetId(),
			Type: n.GetType(),
			Position: Position{
				X: float64(n.Position.GetX()),
				Y: float64(n.Position.GetY()),
			},
		}

		switch n.Type {
		case workflowv1.NodeType_TRIGGER:

			b, err := n.Data.MarshalJSON()
			if err != nil {
				return err
			}

			var trigger NodeTrigger
			if err := json.Unmarshal(b, &trigger); err != nil {
				return err
			}

			m.Trigger = trigger.Key
			node.Trigger = &trigger

		case workflowv1.NodeType_CONDITION:
			b, err := n.Data.MarshalJSON()
			if err != nil {
				return err
			}

			var cond NodeCondition
			if err := json.Unmarshal(b, &cond); err != nil {
				return err
			}

			exp, err := generateCEL(cond.Conditions)
			if err != nil {
				return err
			}
			cond.Expression = exp

			node.Condition = &cond

		case workflowv1.NodeType_ACTION:
			b, err := n.Data.MarshalJSON()
			if err != nil {
				return err
			}

			var act NodeAction
			if err := json.Unmarshal(b, &act); err != nil {
				return err
			}

			node.Action = &act
		}

		internalNodes = append(internalNodes, node)
	}

	data, err := json.Marshal(internalNodes)
	if err != nil {
		zap.L().Error("failed to marshal internalNodes", zap.Error(err))
		return err
	}

	m.Nodes = datatypes.JSON(data)
	return nil
}

func (m *Flow) SetEdges(edges []*workflowv1.Edge) error {
	b, err := json.Marshal(edges)
	if err != nil {
		zap.L().Error("failed to marshal edges", zap.Error(err))
		return err
	}

	m.Edges = datatypes.JSON(b)
	return nil
}

func (m *Flow) SetOverview(overflow *workflowv1.Overflow) error {
	b, err := json.Marshal(overflow)
	if err != nil {
		zap.L().Error("failed to marshal overflow", zap.Error(err))
		return err
	}

	m.Overflow = datatypes.JSON(b)
	return nil
}

func (m *Flow) GetNodes() []*workflowv1.Node {
	var rawNode []*Node
	if err := json.Unmarshal(m.Nodes, &rawNode); err != nil {
		zap.L().Error("failed to unmarshal Nodes", zap.Error(err))
		return nil
	}

	var nodes []*workflowv1.Node
	for _, n := range rawNode {

		node := &workflowv1.Node{
			Id:   n.ID,
			Type: n.Type,
			Position: &workflowv1.Position{
				X: float32(n.Position.X),
				Y: float32(n.Position.Y),
			},
		}

		switch n.Type {
		case workflowv1.NodeType_TRIGGER:
			b, err := json.Marshal(&n.Trigger)
			if err != nil {
				return nil
			}

			var data map[string]any
			if err := json.Unmarshal(b, &data); err != nil {
				return nil
			}

			obj, err := structpb.NewStruct(data)
			if err != nil {
				return nil
			}

			node.Data = obj
		case workflowv1.NodeType_CONDITION:
			b, err := json.Marshal(&n.Condition)
			if err != nil {
				return nil
			}

			var data map[string]any
			if err := json.Unmarshal(b, &data); err != nil {
				return nil
			}

			obj, err := structpb.NewStruct(data)
			if err != nil {
				return nil
			}

			node.Data = obj
		case workflowv1.NodeType_ACTION:
			b, err := json.Marshal(&n.Action)
			if err != nil {
				return nil
			}

			var data map[string]any
			if err := json.Unmarshal(b, &data); err != nil {
				return nil
			}

			obj, err := structpb.NewStruct(data)
			if err != nil {
				return nil
			}

			node.Data = obj
		}

		nodes = append(nodes, node)
	}

	return nodes
}

func (m *Flow) GetEdges() []*workflowv1.Edge {
	var edges []*workflowv1.Edge
	if err := json.Unmarshal(m.Edges, &edges); err != nil {
		zap.L().Error("failed to unmarshal Nodes", zap.Error(err))
		return nil
	}
	return edges
}

func (m *Flow) GetOverflow() *workflowv1.Overflow {
	var overflow *workflowv1.Overflow
	if err := json.Unmarshal(m.Overflow, &overflow); err != nil {
		zap.L().Error("failed to unmarshal Nodes", zap.Error(err))
		return nil
	}
	return overflow
}

type Node struct {
	ID        string              `json:"id" validate:"required"`
	Type      workflowv1.NodeType `json:"node_type" validate:"required"` // enum as string: TRIGGER, CONDITION, ACTION
	Trigger   *NodeTrigger        `json:"trigger,omitempty" validate:"required_if=NodeType NODE_TYPE_TRIGGER"`
	Condition *NodeCondition      `json:"condition,omitempty" validate:"required_if=NodeType NODE_TYPE_CONDITION"`
	Action    *NodeAction         `json:"action,omitempty" validate:"required_if=NodeType ACTION"`
	Position  Position            `json:"position,omitempty" validate:"required"`
}

type NodeTrigger struct {
	Key         string `json:"key"` // TriggerType as string
	Label       string `json:"label,omitempty"`
	Description string `json:"description,omitempty"`
}

type Condition struct {
	Field    string          `json:"field" validate:"required"`
	Operator string          `json:"operator" validate:"required"`
	Value    *structpb.Value `json:"value" validate:"required"`
}

type NodeCondition struct {
	Conditions []*Condition `json:"conditions" validate:"required"`
	Expression string       `json:"expression"`
}

type NodeAction struct {
	Type       string         `json:"type" validate:"required"`
	Parameters datatypes.JSON `json:"parameters" validate:"required"`
}

type Edge struct {
	ID     string `json:"id" validate:"required"`
	Type   string `json:"type" validate:"required"`
	Source string `json:"source" validate:"required"`
	Target string `json:"target" validate:"required"`
}

type Position struct {
	X float64 `json:"x" validate:"required"`
	Y float64 `json:"y" validate:"required"`
}

type Overflow struct {
	X    float64 `json:"x" validate:"required"`
	Y    float64 `json:"y" validate:"required"`
	Zoom float64 `json:"zoom"`
}

func generateCEL(conditions []*Condition) (string, error) {
	exprs := make([]string, 0)
	for _, c := range conditions {
		var valueStr string
		switch c.Value.Kind.(type) {
		case *structpb.Value_StringValue:
			valueStr = fmt.Sprintf(`"%s"`, c.Value.GetStringValue()) // tambahkan tanda kutip
		case *structpb.Value_NumberValue:
			valueStr = fmt.Sprintf(`%v`, c.Value.GetNumberValue())
		case *structpb.Value_BoolValue:
			valueStr = fmt.Sprintf(`%t`, c.Value.GetBoolValue())
		default:
			return "", fmt.Errorf("unsupported value type: %T", c.Value)
		}
		part := fmt.Sprintf(`%s %s %s`, c.Field, c.Operator, valueStr)
		exprs = append(exprs, part)
	}
	return strings.Join(exprs, " && "), nil
}
