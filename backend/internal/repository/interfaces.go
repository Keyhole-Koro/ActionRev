package repository

import (
	"context"

	"github.com/keyhole-koro/actionrev/internal/domain"
)

// ---------------------------------------------------------------------------
// ListOptions
// ---------------------------------------------------------------------------

type ListOptions struct {
	Limit  int
	Offset int
}

// ---------------------------------------------------------------------------
// Document
// ---------------------------------------------------------------------------

type DocumentRepository interface {
	Create(ctx context.Context, doc domain.Document) error
	GetByID(ctx context.Context, id string) (domain.Document, error)
	UpdateStatus(ctx context.Context, id string, status domain.DocumentStatus) error
	ListByWorkspace(ctx context.Context, workspaceID string, opts ListOptions) ([]domain.Document, error)
	Delete(ctx context.Context, id string) error
}

type ChunkRepository interface {
	BatchCreate(ctx context.Context, chunks []domain.DocumentChunk) error
	ListByDocument(ctx context.Context, documentID string) ([]domain.DocumentChunk, error)
}

// ---------------------------------------------------------------------------
// Graph
// ---------------------------------------------------------------------------

type NodeRepository interface {
	BatchUpsert(ctx context.Context, nodes []domain.Node) error
	ListByDocument(ctx context.Context, documentID string) ([]domain.Node, error)
	ListBySourceFile(ctx context.Context, documentID, sourceFilename string) ([]domain.Node, error)
}

type EdgeRepository interface {
	BatchUpsert(ctx context.Context, edges []domain.Edge) error
	ListByDocument(ctx context.Context, documentID string) ([]domain.Edge, error)
}

// ---------------------------------------------------------------------------
// Workspace
// ---------------------------------------------------------------------------

type WorkspaceRepository interface {
	Create(ctx context.Context, ws domain.Workspace) error
	GetByID(ctx context.Context, id string) (domain.Workspace, error)
	Update(ctx context.Context, ws domain.Workspace) error
	ListByOwner(ctx context.Context, ownerID string) ([]domain.Workspace, error)
	AddMember(ctx context.Context, member domain.WorkspaceMember) error
	UpdateMemberRole(ctx context.Context, workspaceID, userID string, role domain.MemberRole) error
	RemoveMember(ctx context.Context, workspaceID, userID string) error
	GetMember(ctx context.Context, workspaceID, userID string) (domain.WorkspaceMember, error)
	ListMembers(ctx context.Context, workspaceID string) ([]domain.WorkspaceMember, error)
	GetPlanLimits(ctx context.Context, plan domain.WorkspacePlan) (domain.PlanLimits, error)
}

// ---------------------------------------------------------------------------
// User
// ---------------------------------------------------------------------------

type UserRepository interface {
	Upsert(ctx context.Context, user domain.User) (isNew bool, err error)
	GetByID(ctx context.Context, userID string) (domain.User, error)
}

// ---------------------------------------------------------------------------
// Upload (GCS)
// ---------------------------------------------------------------------------

type UploadRepository interface {
	// 署名付き PUT URL を発行する (フロントエンドから直接 GCS にアップロード)
	CreateSignedUploadURL(ctx context.Context, objectPath, mimeType string, fileSizeBytes int64) (url string, err error)
	Delete(ctx context.Context, objectPath string) error
}
