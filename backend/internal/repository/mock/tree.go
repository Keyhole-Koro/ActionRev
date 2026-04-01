package mock

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"slices"
	"sync"
	"time"

	graphv1 "github.com/synthify/backend/gen/synthify/graph/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const seedWorkspaceID = "00000000-0000-4000-8000-000000000001"

type TreeRepository struct {
	mu               sync.RWMutex
	rootByWorkspace  map[string]string
	nodesByWorkspace map[string]map[string]*graphv1.PaperNode
	notesByNode      map[string][]*graphv1.PaperNote
	actionsByNode    map[string][]*graphv1.ActionRequest
}

func NewTreeRepository() *TreeRepository {
	now := time.Now()
	rootID := "paper_root_strategy"
	claimID := "paper_claim_growth"
	evidenceID := "paper_evidence_retention"
	actionID := "paper_action_validate"

	root := &graphv1.PaperNode{
		PaperNodeId:  rootID,
		WorkspaceId:  seedWorkspaceID,
		ChildIds:     []string{claimID, actionID},
		Title:        "Growth Strategy",
		Description:  "Workspace 全体の戦略仮説を束ねるルート",
		Content:      "主要な仮説、根拠、次のアクションをここから整理する。",
		Category:     graphv1.PaperNodeCategory_PAPER_NODE_CATEGORY_CONCEPT,
		Scope:        graphv1.PaperNodeScope_PAPER_NODE_SCOPE_WORKSPACE,
		DisplayOrder: 0,
		Status:       graphv1.PaperNodeStatus_PAPER_NODE_STATUS_READY,
		Meta: &graphv1.PaperMeta{
			Badges: []string{"root"},
		},
		CreatedAt: timestamppb.New(now.Add(-2 * time.Hour)),
		UpdatedAt: timestamppb.New(now.Add(-30 * time.Minute)),
	}
	claim := &graphv1.PaperNode{
		PaperNodeId:       claimID,
		WorkspaceId:       seedWorkspaceID,
		ParentId:          rootID,
		ChildIds:          []string{evidenceID},
		Title:             "Retention improvement drives growth",
		Description:       "成長戦略の主張ノード",
		Content:           "解約率を下げる施策が最も短期で LTV を押し上げる。",
		Category:          graphv1.PaperNodeCategory_PAPER_NODE_CATEGORY_CLAIM,
		Scope:             graphv1.PaperNodeScope_PAPER_NODE_SCOPE_CANONICAL,
		DisplayOrder:      0,
		Status:            graphv1.PaperNodeStatus_PAPER_NODE_STATUS_READY,
		SourceDocumentIds: []string{"doc_demo"},
		Meta: &graphv1.PaperMeta{
			Badges: []string{"claim"},
			Relations: []*graphv1.NodeRelationSummary{
				{
					TargetNodeId: actionID,
					RelationType: graphv1.PaperRelationType_PAPER_RELATION_TYPE_REFERENCES,
					Label:        "needs validation",
				},
			},
		},
		Blocks: []*graphv1.PaperBlock{
			{
				Kind: &graphv1.PaperBlock_Documents{
					Documents: &graphv1.DocumentsBlock{DocumentIds: []string{"doc_demo"}},
				},
			},
		},
		CreatedAt: timestamppb.New(now.Add(-90 * time.Minute)),
		UpdatedAt: timestamppb.New(now.Add(-20 * time.Minute)),
	}
	evidence := &graphv1.PaperNode{
		PaperNodeId:       evidenceID,
		WorkspaceId:       seedWorkspaceID,
		ParentId:          claimID,
		Title:             "Users cite onboarding clarity",
		Description:       "インタビュー由来の根拠",
		Content:           "最近のユーザーインタビューでは、初期オンボーディングの理解しやすさが継続利用の要因として挙がっている。",
		Category:          graphv1.PaperNodeCategory_PAPER_NODE_CATEGORY_EVIDENCE,
		Scope:             graphv1.PaperNodeScope_PAPER_NODE_SCOPE_DOCUMENT,
		DisplayOrder:      0,
		Status:            graphv1.PaperNodeStatus_PAPER_NODE_STATUS_READY,
		SourceDocumentIds: []string{"doc_demo"},
		CreatedAt:         timestamppb.New(now.Add(-70 * time.Minute)),
		UpdatedAt:         timestamppb.New(now.Add(-10 * time.Minute)),
	}
	action := &graphv1.PaperNode{
		PaperNodeId:  actionID,
		WorkspaceId:  seedWorkspaceID,
		ParentId:     rootID,
		Title:        "Validate with pricing cohort",
		Description:  "次の検証アクション",
		Content:      "価格改定前後のコホートで retention の差を確認する。",
		Category:     graphv1.PaperNodeCategory_PAPER_NODE_CATEGORY_ACTION,
		Scope:        graphv1.PaperNodeScope_PAPER_NODE_SCOPE_WORKSPACE,
		DisplayOrder: 1,
		Status:       graphv1.PaperNodeStatus_PAPER_NODE_STATUS_PENDING,
		CreatedAt:    timestamppb.New(now.Add(-80 * time.Minute)),
		UpdatedAt:    timestamppb.New(now.Add(-5 * time.Minute)),
	}

	note := &graphv1.PaperNote{
		NoteId:      "note_review_pricing",
		WorkspaceId: seedWorkspaceID,
		PaperNodeId: actionID,
		Kind:        graphv1.PaperNoteKind_PAPER_NOTE_KIND_REVIEW,
		Title:       "AI asks for review",
		Body:        "価格コホートの切り方が妥当か確認してください。",
		Priority:    graphv1.NotePriority_NOTE_PRIORITY_MEDIUM,
		CreatedAt:   timestamppb.New(now.Add(-15 * time.Minute)),
		UpdatedAt:   timestamppb.New(now.Add(-15 * time.Minute)),
	}
	actionRequest := &graphv1.ActionRequest{
		ActionRequestId: "ar_review_pricing",
		WorkspaceId:     seedWorkspaceID,
		PaperNodeId:     actionID,
		NoteId:          note.NoteId,
		Type:            graphv1.ActionRequestType_ACTION_REQUEST_TYPE_REVIEW,
		Title:           "Review pricing cohort split",
		Body:            "有意差を見る区切り方としてこの cohort 分割で良いか判断が必要。",
		Priority:        graphv1.ActionRequestPriority_ACTION_REQUEST_PRIORITY_MEDIUM,
		Status:          graphv1.ActionRequestStatus_ACTION_REQUEST_STATUS_OPEN,
		RequestedBy:     "ai",
		AssignedTo:      "human",
		CreatedAt:       timestamppb.New(now.Add(-15 * time.Minute)),
	}

	return &TreeRepository{
		rootByWorkspace: map[string]string{
			seedWorkspaceID: rootID,
		},
		nodesByWorkspace: map[string]map[string]*graphv1.PaperNode{
			seedWorkspaceID: {
				rootID:     root,
				claimID:    claim,
				evidenceID: evidence,
				actionID:   action,
			},
		},
		notesByNode: map[string][]*graphv1.PaperNote{
			actionID: {note},
		},
		actionsByNode: map[string][]*graphv1.ActionRequest{
			actionID: {actionRequest},
		},
	}
}

func (r *TreeRepository) GetWorkspaceTree(_ context.Context, workspaceID string) (*graphv1.WorkspaceTree, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	rootID, ok := r.rootByWorkspace[workspaceID]
	if !ok {
		return nil, nil
	}
	nodes := r.nodesByWorkspace[workspaceID]
	items := make([]*graphv1.PaperNode, 0, len(nodes))
	for _, node := range nodes {
		items = append(items, clonePaperNode(node))
	}
	slices.SortFunc(items, func(a, b *graphv1.PaperNode) int {
		if a.GetDisplayOrder() == b.GetDisplayOrder() {
			return compareStrings(a.GetPaperNodeId(), b.GetPaperNodeId())
		}
		if a.GetDisplayOrder() < b.GetDisplayOrder() {
			return -1
		}
		return 1
	})

	return &graphv1.WorkspaceTree{
		WorkspaceId: workspaceID,
		RootNodeId:  rootID,
		Nodes:       items,
	}, nil
}

func (r *TreeRepository) ListPaperNodeChildren(_ context.Context, workspaceID string, parentID string) ([]*graphv1.PaperNode, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	parent := r.nodesByWorkspace[workspaceID][parentID]
	if parent == nil {
		return nil, nil
	}

	children := make([]*graphv1.PaperNode, 0, len(parent.GetChildIds()))
	for _, childID := range parent.GetChildIds() {
		if child := r.nodesByWorkspace[workspaceID][childID]; child != nil {
			children = append(children, clonePaperNode(child))
		}
	}

	return children, nil
}

func (r *TreeRepository) GetPaperNode(_ context.Context, workspaceID string, paperNodeID string) (*graphv1.PaperNode, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return clonePaperNode(r.nodesByWorkspace[workspaceID][paperNodeID]), nil
}

func (r *TreeRepository) CreatePaperNode(_ context.Context, workspaceID string, parentID string, title string, description string, content string, category graphv1.PaperNodeCategory, scope graphv1.PaperNodeScope, sourceDocumentIDs []string) (*graphv1.PaperNode, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	parent := r.nodesByWorkspace[workspaceID][parentID]
	if parent == nil {
		return nil, errors.New("parent paper node not found")
	}

	now := timestamppb.Now()
	nodeID, err := newTreeID("paper")
	if err != nil {
		return nil, err
	}

	node := &graphv1.PaperNode{
		PaperNodeId:       nodeID,
		WorkspaceId:       workspaceID,
		ParentId:          parentID,
		Title:             title,
		Description:       description,
		Content:           content,
		Category:          category,
		Scope:             scope,
		DisplayOrder:      uint32(len(parent.GetChildIds())),
		Status:            graphv1.PaperNodeStatus_PAPER_NODE_STATUS_DRAFT,
		SourceDocumentIds: append([]string(nil), sourceDocumentIDs...),
		Meta: &graphv1.PaperMeta{
			IsNew: true,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	r.nodesByWorkspace[workspaceID][nodeID] = node
	parent.ChildIds = append(parent.ChildIds, nodeID)
	parent.UpdatedAt = timestamppb.Now()

	return clonePaperNode(node), nil
}

func (r *TreeRepository) UpdatePaperNode(_ context.Context, workspaceID string, paperNodeID string, title string, description string, content string, status graphv1.PaperNodeStatus) (*graphv1.PaperNode, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	node := r.nodesByWorkspace[workspaceID][paperNodeID]
	if node == nil {
		return nil, nil
	}

	if title != "" {
		node.Title = title
	}
	if description != "" {
		node.Description = description
	}
	if content != "" {
		node.Content = content
	}
	if status != graphv1.PaperNodeStatus_PAPER_NODE_STATUS_UNSPECIFIED {
		node.Status = status
	}
	if node.Meta != nil {
		node.Meta.IsNew = false
	}
	node.UpdatedAt = timestamppb.Now()

	return clonePaperNode(node), nil
}

func (r *TreeRepository) ReorderPaperNode(_ context.Context, workspaceID string, paperNodeID string, newParentID string, insertBeforeID string) (*graphv1.PaperNode, []string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	workspaceNodes := r.nodesByWorkspace[workspaceID]
	node := workspaceNodes[paperNodeID]
	newParent := workspaceNodes[newParentID]
	if node == nil || newParent == nil {
		return nil, nil, nil
	}

	if oldParent := workspaceNodes[node.GetParentId()]; oldParent != nil {
		oldParent.ChildIds = removeString(oldParent.ChildIds, paperNodeID)
		reindexChildren(workspaceNodes, oldParent.ChildIds)
		oldParent.UpdatedAt = timestamppb.Now()
	}

	childIDs := append([]string(nil), newParent.GetChildIds()...)
	insertAt := len(childIDs)
	if insertBeforeID != "" {
		for idx, childID := range childIDs {
			if childID == insertBeforeID {
				insertAt = idx
				break
			}
		}
	}
	childIDs = slices.Insert(childIDs, insertAt, paperNodeID)
	newParent.ChildIds = childIDs
	newParent.UpdatedAt = timestamppb.Now()

	node.ParentId = newParentID
	node.UpdatedAt = timestamppb.Now()
	reindexChildren(workspaceNodes, newParent.ChildIds)

	return clonePaperNode(node), append([]string(nil), newParent.ChildIds...), nil
}

func (r *TreeRepository) ListNodeNotes(_ context.Context, workspaceID string, paperNodeID string) ([]*graphv1.PaperNote, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.nodesByWorkspace[workspaceID][paperNodeID] == nil {
		return nil, nil
	}
	notes := r.notesByNode[paperNodeID]
	out := make([]*graphv1.PaperNote, 0, len(notes))
	for _, note := range notes {
		out = append(out, clonePaperNote(note))
	}
	return out, nil
}

func (r *TreeRepository) CreateNodeNote(_ context.Context, workspaceID string, paperNodeID string, kind graphv1.PaperNoteKind, title string, body string, priority graphv1.NotePriority) (*graphv1.PaperNote, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.nodesByWorkspace[workspaceID][paperNodeID] == nil {
		return nil, errors.New("paper node not found")
	}

	now := timestamppb.Now()
	noteID, err := newTreeID("note")
	if err != nil {
		return nil, err
	}
	note := &graphv1.PaperNote{
		NoteId:      noteID,
		WorkspaceId: workspaceID,
		PaperNodeId: paperNodeID,
		Kind:        kind,
		Title:       title,
		Body:        body,
		Priority:    priority,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	r.notesByNode[paperNodeID] = append(r.notesByNode[paperNodeID], note)
	return clonePaperNote(note), nil
}

func (r *TreeRepository) ListNodeActionRequests(_ context.Context, workspaceID string, paperNodeID string) ([]*graphv1.ActionRequest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.nodesByWorkspace[workspaceID][paperNodeID] == nil {
		return nil, nil
	}
	items := r.actionsByNode[paperNodeID]
	out := make([]*graphv1.ActionRequest, 0, len(items))
	for _, item := range items {
		out = append(out, cloneActionRequest(item))
	}
	return out, nil
}

func (r *TreeRepository) ResolveActionRequest(_ context.Context, workspaceID string, actionRequestID string) (*graphv1.ActionRequest, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for nodeID, items := range r.actionsByNode {
		for _, item := range items {
			if item.GetWorkspaceId() == workspaceID && item.GetActionRequestId() == actionRequestID {
				item.Status = graphv1.ActionRequestStatus_ACTION_REQUEST_STATUS_RESOLVED
				item.ResolvedAt = timestamppb.Now()
				r.actionsByNode[nodeID] = items
				return cloneActionRequest(item), nil
			}
		}
	}
	return nil, nil
}

func (r *TreeRepository) DismissActionRequest(_ context.Context, workspaceID string, actionRequestID string) (*graphv1.ActionRequest, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for nodeID, items := range r.actionsByNode {
		for _, item := range items {
			if item.GetWorkspaceId() == workspaceID && item.GetActionRequestId() == actionRequestID {
				item.Status = graphv1.ActionRequestStatus_ACTION_REQUEST_STATUS_DISMISSED
				item.ResolvedAt = timestamppb.Now()
				r.actionsByNode[nodeID] = items
				return cloneActionRequest(item), nil
			}
		}
	}
	return nil, nil
}

func clonePaperNode(node *graphv1.PaperNode) *graphv1.PaperNode {
	if node == nil {
		return nil
	}
	return proto.Clone(node).(*graphv1.PaperNode)
}

func clonePaperNote(note *graphv1.PaperNote) *graphv1.PaperNote {
	if note == nil {
		return nil
	}
	return proto.Clone(note).(*graphv1.PaperNote)
}

func cloneActionRequest(item *graphv1.ActionRequest) *graphv1.ActionRequest {
	if item == nil {
		return nil
	}
	return proto.Clone(item).(*graphv1.ActionRequest)
}

func reindexChildren(nodes map[string]*graphv1.PaperNode, childIDs []string) {
	for idx, childID := range childIDs {
		if child := nodes[childID]; child != nil {
			child.DisplayOrder = uint32(idx)
			child.UpdatedAt = timestamppb.Now()
		}
	}
}

func removeString(items []string, target string) []string {
	out := make([]string, 0, len(items))
	for _, item := range items {
		if item != target {
			out = append(out, item)
		}
	}
	return out
}

func compareStrings(a string, b string) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}

func newTreeID(prefix string) (string, error) {
	value := make([]byte, 8)
	if _, err := rand.Read(value); err != nil {
		return "", err
	}
	return prefix + "_" + hex.EncodeToString(value), nil
}
